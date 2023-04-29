package company

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	CreateCompany(ctx context.Context, company *Company) (*Company, error)
	GetCompanyByID(ctx context.Context, id uuid.UUID) (*Company, error)
	UpdateCompany(ctx context.Context, id uuid.UUID, company *Company) (*Company, error)
	DeleteCompany(ctx context.Context, id uuid.UUID) error
}

type companyService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &companyService{repo: repo}
}

func (s *companyService) CreateCompany(ctx context.Context, company *Company) (*Company, error) {
	// Generate a new UUID for the company
	company.ID = uuid.New()
	// determine new company name is unique
	count, err := s.repo.CountByName(ctx, company.Name)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrCompanyExists{company.Name}
	}
	err = s.repo.Create(ctx, company)

	if err != nil {
		return nil, err
	}

	return company, nil
}

func (s *companyService) GetCompanyByID(ctx context.Context, id uuid.UUID) (*Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *companyService) UpdateCompany(ctx context.Context, id uuid.UUID, company *Company) (*Company, error) {
	existingCompany, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if existingCompany == nil {
		return nil, ErrCompanyNotFound{id}
	}

	// Copy over the fields that can be updated
	existingCompany.Name = company.Name
	existingCompany.Description = company.Description
	existingCompany.Employees = company.Employees
	existingCompany.Registered = company.Registered
	existingCompany.Type = company.Type

	updatedCompany, err := s.repo.Update(ctx, existingCompany)

	if err != nil {
		return nil, err
	}

	return updatedCompany, nil
}

func (s *companyService) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	existingCompany, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return err
	}

	if existingCompany == nil {
		return ErrCompanyNotFound{id}
	}

	return s.repo.Delete(ctx, id)
}
