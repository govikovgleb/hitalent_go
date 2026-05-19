package handlers

import (
	"encoding/json"
	"errors"
	"hitalent_go/internal/service"
	"hitalent_go/internal/validator"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	srv *service.Service
}

func New(srv *service.Service) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	dept, err := h.srv.CreateDepartment(req.Name, req.ParentID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			respondError(w, http.StatusNotFound, "parent department not found")
		case errors.Is(err, service.ErrConflict):
			respondError(w, http.StatusConflict, "department with this name already exists")
		case errors.Is(err, validator.ErrValidation):
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondJSON(w, http.StatusCreated, dept)
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	deptIDStr := strings.TrimPrefix(r.URL.Path, "/departments/")
	deptIDStr = strings.TrimSuffix(deptIDStr, "/employees/")

	deptID, err := strconv.ParseUint(deptIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid department ID")
		return
	}

	var req struct {
		FullName string     `json:"full_name"`
		Position string     `json:"position"`
		HiredAt  *time.Time `json:"hired_at,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	emp, err := h.srv.CreateEmployee(uint(deptID), req.FullName, req.Position, req.HiredAt)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			respondError(w, http.StatusNotFound, "department not found")
		case errors.Is(err, validator.ErrValidation):
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondJSON(w, http.StatusCreated, emp)
}

func (h *Handler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	deptIDStr := strings.TrimPrefix(r.URL.Path, "/departments/")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadGateway, "invalid department ID")
		return
	}

	depth := 1
	if d := r.URL.Query().Get("depth"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val >= 1 && val <= 5 {
			depth = val
		}
	}

	includeEmployees := true
	if ie := r.URL.Query().Get("include_employees"); ie == "false" {
		includeEmployees = false
	}

	dept, err := h.srv.GetDepartment(uint(deptID), depth, includeEmployees)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			respondError(w, http.StatusNotFound, "department not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"department": dept,
		"employees":  dept.Employees,
		"children":   dept.Children,
	})
}

func (h *Handler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	deptIDStr := strings.TrimPrefix(r.URL.Path, "/departments/")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid department ID")
		return
	}

	var req struct {
		Name     *string `json:"name"`
		ParentID *uint   `json:"parent_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	dept, err := h.srv.UpdateDepartment(uint(deptID), req.Name, req.ParentID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			respondError(w, http.StatusNotFound, "department not found")
		case errors.Is(err, service.ErrConflict):
			respondError(w, http.StatusConflict, "department with this name already exists")
		case errors.Is(err, service.ErrCycleDetected):
			respondError(w, http.StatusBadRequest, "cycle detected")
		case errors.Is(err, validator.ErrValidation):
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondJSON(w, http.StatusOK, dept)
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	deptIDStr := strings.TrimPrefix(r.URL.Path, "/departments/")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid department ID")
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode != "cascade" && mode != "reassign" {
		respondError(w, http.StatusBadRequest, "invalid mode")
		return
	}

	var reassignDeptID *uint
	if rdi := r.URL.Query().Get("reassign_dept_id"); rdi != "" {
		val, err := strconv.ParseUint(rdi, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid reassign department ID")
			return
		}
		uval := uint(val)
		reassignDeptID = &uval
	}

	if err := h.srv.DeleteDepartment(uint(deptID), mode, reassignDeptID); err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			respondError(w, http.StatusNotFound, "department not found")
		case errors.Is(err, service.ErrDeleteModeInvalid):
			respondError(w, http.StatusBadRequest, "invalid delete mode")
		case errors.Is(err, validator.ErrValidation):
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
