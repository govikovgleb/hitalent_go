package repository

import (
	"hitalent_go/internal/models"

	"gorm.io/gorm"
)

type DepartmentRepo struct {
	DB *gorm.DB
}

type DepartmentRepository interface {
	CreateDepartment(dept *models.Department) error
	GetDepartmentByID(id uint) (*models.Department, error)
	GetAllDepartments() ([]models.Department, error)
	GetDepartmentWithEmployees(id uint) (*models.Department, error)
	GetChildren(parentID uint) ([]models.Department, error)
	UpdateDepartment(dept *models.Department) error
	DeleteDepartmentCascade(id uint) error
	DeleteDepartment(id uint) error
	CheckDepartmentExists(id uint) (bool, error)
	CheckDuplicateName(name string, parentID *uint, excludeID uint) (bool, error)
}

func NewDepartmentRepository(db *gorm.DB) (DepartmentRepository, error) {
	return &DepartmentRepo{DB: db}, nil
}

func (r *DepartmentRepo) CreateDepartment(dept *models.Department) error {
	return r.DB.Create(dept).Error
}

func (r *DepartmentRepo) GetDepartmentByID(id uint) (*models.Department, error) {
	var dept models.Department
	err := r.DB.First(&dept, id).Error
	if err != nil {
		return nil, err
	}

	return &dept, nil
}

func (r *DepartmentRepo) GetAllDepartments() ([]models.Department, error) {
	var depts []models.Department
	err := r.DB.Select("id, parent_id").Find(&depts).Error
	return depts, err
}

func (r *DepartmentRepo) GetDepartmentWithEmployees(id uint) (*models.Department, error) {
	var dept models.Department
	err := r.DB.Preload("Employees").First(&dept, id).Error
	if err != nil {
		return nil, err
	}

	return &dept, nil
}

func (r *DepartmentRepo) GetChildren(parentID uint) ([]models.Department, error) {
	var children []models.Department
	err := r.DB.Where("parent_id = ?", parentID).Find(&children).Error
	return children, err
}

func (r *DepartmentRepo) UpdateDepartment(dept *models.Department) error {
	return r.DB.Save(dept).Error
}

func (r *DepartmentRepo) DeleteDepartmentCascade(id uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Get all descendant department IDs recursively
		descendantIDs, err := r.getDescendantIDs(tx, id)
		if err != nil {
			return err
		}

		// Add the root department to the list
		allIDs := append(descendantIDs, id)

		// Delete all employees in all departments (root + descendants)
		if err := tx.Where("department_id IN ?", allIDs).Delete(&models.Employee{}).Error; err != nil {
			return err
		}

		// Delete all descendant departments first (children before parent)
		// Order by depth descending to avoid FK constraint issues
		for i := len(descendantIDs) - 1; i >= 0; i-- {
			if err := tx.Delete(&models.Department{}, descendantIDs[i]).Error; err != nil {
				return err
			}
		}

		// Finally delete the root department
		return tx.Delete(&models.Department{}, id).Error
	})
}

// getDescendantIDs recursively collects all descendant department IDs
func (r *DepartmentRepo) getDescendantIDs(tx *gorm.DB, parentID uint) ([]uint, error) {
	var children []models.Department
	if err := tx.Where("parent_id = ?", parentID).Select("id").Find(&children).Error; err != nil {
		return nil, err
	}

	if len(children) == 0 {
		return nil, nil
	}

	var allDescendants []uint
	for _, child := range children {
		allDescendants = append(allDescendants, child.ID)
		// Recursively get grandchildren
		grandchildren, err := r.getDescendantIDs(tx, child.ID)
		if err != nil {
			return nil, err
		}
		allDescendants = append(allDescendants, grandchildren...)
	}

	return allDescendants, nil
}

func (r *DepartmentRepo) DeleteDepartment(id uint) error {
	return r.DB.Delete(&models.Department{}, id).Error
}

func (r *DepartmentRepo) CheckDepartmentExists(id uint) (bool, error) {
	var exists bool
	err := r.DB.Model(&models.Department{}).Select("count(*) > 0").Where("id = ?", id).Find(&exists).Error
	return exists, err
}

func (r *DepartmentRepo) CheckDuplicateName(name string, parentID *uint, excludeID uint) (bool, error) {
	var count int64
	query := r.DB.Model(&models.Department{}).Where("name = ? AND id != ?", name, excludeID)
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	err := query.Count(&count).Error
	return count > 0, err
}
