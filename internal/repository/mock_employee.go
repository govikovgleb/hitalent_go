package repository

import (
	"hitalent_go/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) CreateEmployee(emp *models.Employee) error {
	args := m.Called(emp)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetEmployeeByID(id uint) (*models.Employee, error) {
	args := m.Called(id)
	emp, _ := args.Get(0).(*models.Employee)
	return emp, args.Error(1)
}

func (m *MockEmployeeRepository) GetEmployeesByDepartment(deptID uint) ([]models.Employee, error) {
	args := m.Called(deptID)
	emps, _ := args.Get(0).([]models.Employee)
	return emps, args.Error(1)
}

func (m *MockEmployeeRepository) ReassignEmployees(fromDeptID, toDeptID uint) error {
	args := m.Called(fromDeptID, toDeptID)
	return args.Error(0)
}
