package service

import (
	"hitalent_go/internal/models"
	"hitalent_go/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateDepartment(t *testing.T) {
	t.Run("success create root department", func(t *testing.T) {
		mockDeptRepo := new(repository.MockDepartmentRepository)
		mockEmpRepo := new(repository.MockEmployeeRepository)

		svc := New(mockDeptRepo, mockEmpRepo)

		mockDeptRepo.On("CheckDuplicateName", "BA Department", (*uint)(nil), uint(0)).
			Return(false, nil).Once()

		mockDeptRepo.On("CreateDepartment", mock.AnythingOfType("*models.Department")).
			Return(nil).Once().
			Run(func(args mock.Arguments) {
				dept := args.Get(0).(*models.Department)
				assert.Equal(t, "BA Department", dept.Name)
				assert.Nil(t, dept.ParentID)
			})

		dept, err := svc.CreateDepartment("BA Department", nil)

		assert.NoError(t, err)
		assert.NotNil(t, dept)
		assert.Equal(t, "BA Department", dept.Name)
		assert.Nil(t, dept.ParentID)

		mockEmpRepo.AssertNotCalled(t, "CreateEmployee")

		mockDeptRepo.AssertExpectations(t)
	})

	t.Run("success create child department", func(t *testing.T) {
		mockDeptRepo := new(repository.MockDepartmentRepository)
		mockEmpRepo := new(repository.MockEmployeeRepository)
		svc := New(mockDeptRepo, mockEmpRepo)

		parentID := uint(1)

		mockDeptRepo.On("CheckDepartmentExists", uint(1)).
			Return(true, nil).Once()

		mockDeptRepo.On("CheckDuplicateName", "Backend", &parentID, uint(0)).
			Return(false, nil).Once()

		mockDeptRepo.On("CreateDepartment", mock.AnythingOfType("*models.Department")).
			Return(nil).Once()

		dept, err := svc.CreateDepartment("Backend", &parentID)

		assert.NoError(t, err)
		assert.NotNil(t, dept)
		mockDeptRepo.AssertExpectations(t)
	})

	t.Run("parent not found", func(t *testing.T) {
		mockDeptRepo := new(repository.MockDepartmentRepository)
		mockEmpRepo := new(repository.MockEmployeeRepository)
		svc := New(mockDeptRepo, mockEmpRepo)

		parentID := uint(999)

		mockDeptRepo.On("CheckDepartmentExists", uint(999)).
			Return(false, nil).Once()

		dept, err := svc.CreateDepartment("Backend", &parentID)

		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, dept)
		mockDeptRepo.AssertExpectations(t)
	})

	t.Run("duplicate name", func(t *testing.T) {
		mockDeptRepo := new(repository.MockDepartmentRepository)
		mockEmpRepo := new(repository.MockEmployeeRepository)
		svc := New(mockDeptRepo, mockEmpRepo)

		parentID := uint(1)

		mockDeptRepo.On("CheckDepartmentExists", uint(1)).
			Return(true, nil).Once()

		mockDeptRepo.On("CheckDuplicateName", "Backend", &parentID, uint(0)).
			Return(true, nil).Once()

		dept, err := svc.CreateDepartment("Backend", &parentID)

		assert.ErrorIs(t, err, ErrConflict)
		assert.Nil(t, dept)
	})
}
