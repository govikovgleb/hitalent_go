package seed

import (
	"fmt"
	"hitalent_go/internal/service"
	"time"

	"log/slog"
)

type Seeder struct {
	service *service.Service
}

func New(svc *service.Service) *Seeder {
	return &Seeder{service: svc}
}

func (s *Seeder) Run() error {
	slog.Info("Starting database seeding...")

	depts, err := s.service.GetDepartment(1, 1, false)
	if err == nil && depts != nil {
		slog.Info("Database already has data, skipping seed")
		return nil
	}

	company, err := s.service.CreateDepartment("Компания", nil)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	slog.Info("Created", "department", company.Name, "id", company.ID)

	it, err := s.service.CreateDepartment("IT", &company.ID)
	if err != nil {
		return fmt.Errorf("failed to create IT: %w", err)
	}
	slog.Info("Created", "department", it.Name, "id", it.ID)

	backend, err := s.service.CreateDepartment("Backend", &it.ID)
	if err != nil {
		return fmt.Errorf("failed to create Backend: %w", err)
	}

	goTeam, err := s.service.CreateDepartment("Go Team", &backend.ID)
	if err != nil {
		return fmt.Errorf("failed to create Go Team: %w", err)
	}

	frontend, err := s.service.CreateDepartment("Frontend", &it.ID)
	if err != nil {
		return fmt.Errorf("failed to create Frontend: %w", err)
	}

	devops, err := s.service.CreateDepartment("DevOps", &it.ID)
	if err != nil {
		return fmt.Errorf("failed to create DevOps: %w", err)
	}

	hr, err := s.service.CreateDepartment("HR", &company.ID)
	if err != nil {
		return fmt.Errorf("failed to create HR: %w", err)
	}

	sales, err := s.service.CreateDepartment("Sales", &company.ID)
	if err != nil {
		return fmt.Errorf("failed to create Sales: %w", err)
	}

	b2b, err := s.service.CreateDepartment("B2B Sales", &sales.ID)
	if err != nil {
		return fmt.Errorf("failed to create B2B Sales: %w", err)
	}

	b2c, err := s.service.CreateDepartment("B2C Sales", &sales.ID)
	if err != nil {
		return fmt.Errorf("failed to create B2C Sales: %w", err)
	}

	slog.Info("Created department hierarchy", "total_departments", 10)

	hireDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	employees := []struct {
		deptID   uint
		name     string
		position string
	}{
		{it.ID, "Иванов Иван", "CTO"},
		{it.ID, "Александр Глеб", "Team Lead"},
		{it.ID, "Имя Имянов", "Architect"},

		{backend.ID, "Backend Developer 1", "Senior Go Developer"},
		{backend.ID, "Backend Developer 2", "Middle Go Developer"},
		{backend.ID, "Backend Developer 3", "Junior Go Developer"},
		{backend.ID, "Database Engineer", "PostgreSQL DBA"},

		{goTeam.ID, "Go Expert 1", "Principal Engineer"},
		{goTeam.ID, "Go Expert 2", "Staff Engineer"},

		{frontend.ID, "React Developer 1", "Senior Frontend"},
		{frontend.ID, "Vue Developer 1", "Middle Frontend"},
		{frontend.ID, "Angular Developer 1", "Junior Frontend"},

		{devops.ID, "DevOps Engineer 1", "Senior DevOps"},
		{devops.ID, "DevOps Engineer 2", "Middle DevOps"},

		{hr.ID, "HR Manager 1", "Senior HR"},
		{hr.ID, "HR Manager 2", "Recruiter"},

		{sales.ID, "Sales Director", "VP of Sales"},
		{sales.ID, "Sales Manager", "Senior Sales"},

		{b2b.ID, "B2B Manager 1", "Account Executive"},
		{b2b.ID, "B2B Manager 2", "Sales Development"},
		{b2b.ID, "B2B Manager 3", "Customer Success"},

		{b2c.ID, "B2C Manager 1", "Retail Sales"},
		{b2c.ID, "B2C Manager 2", "Online Sales"},
		{b2c.ID, "B2C Manager 3", "Call Center"},
	}

	for _, emp := range employees {
		_, err := s.service.CreateEmployee(emp.deptID, emp.name, emp.position, &hireDate)
		if err != nil {
			return fmt.Errorf("failed to create employee %s: %w", emp.name, err)
		}
	}

	slog.Info("Seeding completed successfully",
		"departments", 10,
		"employees", len(employees),
		"message", "Ready for testing cascade deletion!")

	return nil
}
