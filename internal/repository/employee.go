package repository

import (
	"hitalent_go/internal/models"

	"gorm.io/gorm"
)

type EmployeeRepo struct {
	DB *gorm.DB
}

type EmployeeRepository interface {
	CreateEmployee(emp *models.Employee) error
	GetEmployeeByID(id uint) (*models.Employee, error)
	GetEmployeesByDepartment(deptID uint) ([]models.Employee, error)
	ReassignEmployees(fromDeptID, toDeptID uint) error
}

func NewEmployeeRepository(db *gorm.DB) (EmployeeRepository, error) {
	return &EmployeeRepo{DB: db}, nil
}

func (r *EmployeeRepo) ReassignEmployees(fromDeptID, toDeptID uint) error {
	return r.DB.Model(&models.Employee{}).
		Where("department_id = ?", fromDeptID).
		Update("department_id", toDeptID).Error
}

func (r *EmployeeRepo) CreateEmployee(emp *models.Employee) error {
	return r.DB.Create(emp).Error
}

func (r *EmployeeRepo) GetEmployeeByID(id uint) (*models.Employee, error) {
	var emp models.Employee
	err := r.DB.First(&emp, id).Error
	if err != nil {
		return nil, err
	}

	return &emp, nil
}

func (r *EmployeeRepo) GetEmployeesByDepartment(deptID uint) ([]models.Employee, error) {
	var emps []models.Employee
	err := r.DB.Where("department_id = ?", deptID).Order("created_at").Find(&emps).Error
	return emps, err
}
