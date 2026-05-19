package repository

import (
	"hitalent_go/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockDepartmentRepository struct {
	mock.Mock
}

func (m *MockDepartmentRepository) CreateDepartment(dept *models.Department) error {
	args := m.Called(dept)
	return args.Error(0)
}

func (m *MockDepartmentRepository) CheckDepartmentExists(id uint) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDepartmentRepository) CheckDuplicateName(name string, parentID *uint, deptID uint) (bool, error) {
	args := m.Called(name, parentID, deptID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDepartmentRepository) GetDepartmentByID(id uint) (*models.Department, error) {
	args := m.Called(id)
	dept, _ := args.Get(0).(*models.Department)
	return dept, args.Error(1)
}

func (m *MockDepartmentRepository) GetAllDepartments() ([]models.Department, error) {
	args := m.Called()
	depts, _ := args.Get(0).([]models.Department)
	return depts, args.Error(1)
}

func (m *MockDepartmentRepository) GetDepartmentWithEmployees(id uint) (*models.Department, error) {
	args := m.Called(id)
	dept, _ := args.Get(0).(*models.Department)
	return dept, args.Error(1)
}

func (m *MockDepartmentRepository) GetChildren(parentID uint) ([]models.Department, error) {
	args := m.Called(parentID)
	depts, _ := args.Get(0).([]models.Department)
	return depts, args.Error(1)
}

func (m *MockDepartmentRepository) UpdateDepartment(dept *models.Department) error {
	args := m.Called(dept)
	return args.Error(0)
}

func (m *MockDepartmentRepository) DeleteDepartmentCascade(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDepartmentRepository) DeleteDepartment(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
