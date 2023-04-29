package company

import (
	"fmt"
	"github.com/google/uuid"
)

type CompanyType string

const (
	Corporation        CompanyType = "Corporation"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "SoleProprietorship"
)

type Company struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Employees   int         `json:"employees"`
	Registered  bool        `json:"registered"`
	Type        CompanyType `json:"type"`
}

type ErrCompanyNotFound struct {
	id uuid.UUID
}

func (e ErrCompanyNotFound) Error() string {
	return fmt.Sprintf("company %v not found", e.id)
}

type ErrCompanyExists struct {
	name string
}

func (e ErrCompanyExists) Error() string {
	return fmt.Sprintf("company %v exists", e.name)
}
