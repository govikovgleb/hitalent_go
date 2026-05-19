package models

import "time"

type Department struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(200);not null" json:"name"`
	ParentID  *uint     `gorm:"index" json:"parent_id,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Parent    *Department  `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"parent,omitempty"`
	Children  []Department `gorm:"foreignKey:ParentID" json:"-"`
	Employees []Employee   `gorm:"foreignKey:DepartmentID" json:"-"`
}

type Employee struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName     string     `gorm:"type:varchar(200);not null" json:"full_name"`
	Position     string     `gorm:"type:varchar(200);not null" json:"position"`
	DepartmentID uint       `gorm:"index" json:"department_id"`
	HiredAt      *time.Time `gorm:"type:datetime" json:"hired_at,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Department Department `gorm:"foreignKey:DepartmentID;constraint:OnDelete:CASCADE" json:"-"`
}
