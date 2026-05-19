package service

import (
	"errors"
	"hitalent_go/internal/models"
	"hitalent_go/internal/repository"
	"hitalent_go/internal/validator"
	"strings"
	"time"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("conflict")
	ErrDeleteModeInvalid = errors.New("invalid delete mode")
	ErrCycleDetected     = errors.New("cycle detected in department tree")
)

type Service struct {
	deptRepo repository.DepartmentRepository
	empRepo  repository.EmployeeRepository
}

type DepartmentService interface {
	CreateDepartment(name string, parentID *uint) (*models.Department, error)
	GetDepartment(id uint, depth int, includeEmployees bool) (*models.Department, error)
	UpdateDepartment(id uint, name *string, parentID *uint) (*models.Department, error)
	DeleteDepartment(id uint, mode string, reassignDeptID *uint) error
}

type EmployeeService interface {
	CreateEmployee(deptID uint, fullName, position string, hiredAt *time.Time) (*models.Employee, error)
}

func New(deptRepo repository.DepartmentRepository, empRepo repository.EmployeeRepository) *Service {
	return &Service{deptRepo: deptRepo, empRepo: empRepo}
}

func (s *Service) CreateDepartment(name string, parentID *uint) (*models.Department, error) {
	name = strings.TrimSpace(name)
	if err := validator.ValidateDepartmentName(name); err != nil {
		return nil, err
	}

	if parentID != nil {
		exists, err := s.deptRepo.CheckDepartmentExists(*parentID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrNotFound
		}
	}

	dup, err := s.deptRepo.CheckDuplicateName(name, parentID, 0)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, ErrConflict
	}

	dept := &models.Department{
		Name:     name,
		ParentID: parentID,
	}

	if err := s.deptRepo.CreateDepartment(dept); err != nil {
		return nil, err
	}

	return dept, nil
}

func (s *Service) CreateEmployee(deptID uint, fullName, position string, hiredAt *time.Time) (*models.Employee, error) {
	fullName = normalizeFullName(strings.TrimSpace(fullName))
	position = strings.TrimSpace(position)

	if err := validator.ValidateEmployeeName(fullName); err != nil {
		return nil, err
	}
	if err := validator.ValidateEmployeePosition(position); err != nil {
		return nil, err
	}

	exists, err := s.deptRepo.CheckDepartmentExists(deptID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	emp := &models.Employee{
		DepartmentID: deptID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}

	if err := s.empRepo.CreateEmployee(emp); err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *Service) GetDepartment(id uint, depth int, includeEmployees bool) (*models.Department, error) {
	if depth < 1 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}

	var dept *models.Department
	var err error

	if depth == 1 && includeEmployees {
		dept, err = s.deptRepo.GetDepartmentWithEmployees(id)
	} else {
		dept, err = s.deptRepo.GetDepartmentByID(id)
	}

	if err != nil {
		return nil, ErrNotFound
	}

	if depth > 1 && includeEmployees {
		emps, err := s.empRepo.GetEmployeesByDepartment(id)
		if err != nil {
			return nil, err
		}
		dept.Employees = emps
	}

	if depth > 1 {
		children, err := s.buildTree(id, depth-1, includeEmployees)
		if err != nil {
			return nil, err
		}
		dept.Children = children
	}

	return dept, nil
}

func (s *Service) UpdateDepartment(id uint, name *string, parentID *uint) (*models.Department, error) {
	dept, err := s.deptRepo.GetDepartmentByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if name != nil {
		trimmedName := strings.TrimSpace(*name)
		if err := validator.ValidateDepartmentName(trimmedName); err != nil {
			return nil, err
		}

		var checkParentID *uint
		if parentID != nil {
			checkParentID = parentID
		} else {
			checkParentID = dept.ParentID
		}

		dup, err := s.deptRepo.CheckDuplicateName(trimmedName, checkParentID, id)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, ErrConflict
		}

		dept.Name = trimmedName
	}

	if parentID != nil {
		if *parentID == id {
			return nil, ErrCycleDetected
		}

		if *parentID != 0 {
			exists, err := s.deptRepo.CheckDepartmentExists(*parentID)
			if err != nil {
				return nil, err
			}
			if !exists {
				return nil, ErrNotFound
			}

			if err := s.checkCycle(id, *parentID); err != nil {
				return nil, err
			}
		} else {
			var zero uint = 0
			parentID = &zero
		}

		dept.ParentID = parentID
	}

	if err := s.deptRepo.UpdateDepartment(dept); err != nil {
		return nil, err
	}

	return dept, nil
}

func (s *Service) DeleteDepartment(id uint, mode string, reassignDeptID *uint) error {
	exists, err := s.deptRepo.CheckDepartmentExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}

	switch mode {
	case "cascade":
		return s.deptRepo.DeleteDepartmentCascade(id)
	case "reassign":
		if reassignDeptID == nil {
			return ErrDeleteModeInvalid
		}
		toExists, err := s.deptRepo.CheckDepartmentExists(*reassignDeptID)
		if err != nil {
			return err
		}
		if !toExists {
			return ErrNotFound
		}
		if err := s.empRepo.ReassignEmployees(id, *reassignDeptID); err != nil {
			return err
		}
		return s.deptRepo.DeleteDepartment(id)
	default:
		return ErrDeleteModeInvalid
	}
}

func (s *Service) buildTree(parentID uint, remainingDepth int, includeEmployees bool) ([]models.Department, error) {
	children, err := s.deptRepo.GetChildren(parentID)
	if err != nil {
		return nil, err
	}

	for i := range children {
		if includeEmployees {
			emps, err := s.empRepo.GetEmployeesByDepartment(children[i].ID)
			if err != nil {
				return nil, err
			}
			children[i].Employees = emps
		}
		if remainingDepth > 1 {
			grandChildren, err := s.buildTree(children[i].ID, remainingDepth-1, includeEmployees)
			if err != nil {
				return nil, err
			}
			children[i].Children = grandChildren
		}
	}

	return children, nil
}

func (s *Service) checkCycle(deptID uint, newParentID uint) error {
	allDepts, err := s.deptRepo.GetAllDepartments()
	if err != nil {
		return err
	}

	parents := make(map[uint]*uint)
	for _, d := range allDepts {
		parents[d.ID] = d.ParentID
	}

	currentID := newParentID
	visited := make(map[uint]bool)

	for currentID != 0 {
		if visited[currentID] {
			return ErrCycleDetected
		}
		visited[currentID] = true

		if currentID == deptID {
			return ErrCycleDetected
		}

		parentID, exists := parents[currentID]
		if !exists || parentID == nil {
			break
		}

		currentID = *parentID
	}
	return nil
}

func normalizeFullName(name string) string {
	fields := strings.Fields(name)
	return strings.Join(fields, " ")
}
