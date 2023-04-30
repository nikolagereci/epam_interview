package model

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
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description,omitempty"`
	Employees   int         `json:"employees" binding:"required"`
	Registered  bool        `json:"registered"`
	Type        CompanyType `json:"type" binding:"required"`
}

type ErrCompanyNotFound struct {
	Id uuid.UUID
}

func (e ErrCompanyNotFound) Error() string {
	return fmt.Sprintf("company %v not found", e.Id)
}

type ErrCompanyExists struct {
	Name string
}

func (e ErrCompanyExists) Error() string {
	return fmt.Sprintf("company %v exists", e.Name)
}
