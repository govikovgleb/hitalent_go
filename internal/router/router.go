package router

import (
	"hitalent_go/internal/handlers"
	"net/http"
	"strings"
)

func New(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/departments/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			if strings.HasSuffix(r.URL.Path, "/employees/") {
				h.CreateEmployee(w, r)
			} else {
				h.CreateDepartment(w, r)
			}
		case http.MethodGet:
			h.GetDepartment(w, r)
		case http.MethodPatch:
			h.UpdateDepartment(w, r)
		case http.MethodDelete:
			h.DeleteDepartment(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
