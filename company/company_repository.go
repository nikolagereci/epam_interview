package company

import (
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/ngereci/xm_interview/model"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	Create(company *model.Company) error
	GetByID(id uuid.UUID) (*model.Company, error)
	Update(company *model.Company) (*model.Company, error)
	Delete(id uuid.UUID) error
	CountByName(name string) (int, error)
}

type companyRepository struct {
	session *gocql.Session
}

func NewRepository(session *gocql.Session) Repository {
	return &companyRepository{session: session}
}

func (r *companyRepository) Create(company *model.Company) error {
	query := r.session.Query(`
		INSERT INTO company (id, name, description, employees, registered, type)
		VALUES (?, ?, ?, ?, ?, ?)
	`, company.ID.String(), company.Name, company.Description, company.Employees, company.Registered, company.Type)

	return query.Exec()
}

func (r *companyRepository) GetByID(id uuid.UUID) (*model.Company, error) {

	query := r.session.Query(`
		SELECT id, name, description, employees, registered, type
		FROM company
		WHERE id = ?
	`, id.String())
	resultMap := make(map[string]any)
	err := query.MapScan(resultMap)
	if err != nil {
		if err == gocql.ErrNotFound {
			log.Warnf("id:%v GetByID not found", id)
			return nil, nil
		}
		log.Errorf("id:%v GetByID error:%v", id, err)
		return nil, err
	}
	company := model.Company{
		ID:          id,
		Name:        resultMap["name"].(string),
		Description: resultMap["description"].(string),
		Employees:   resultMap["employees"].(int),
		Registered:  resultMap["registered"].(bool),
		Type:        model.CompanyType(resultMap["type"].(string)),
	}

	return &company, nil
}
func (r *companyRepository) CountByName(name string) (count int, err error) {

	query := r.session.Query(`
		SELECT COUNT(*)
		FROM company
		WHERE name = ?
	`, name)
	err = query.Scan(&count)
	if err != nil {
		if err == gocql.ErrNotFound {
			log.Warnf("name:%v CountByName not found", name)
			return count, nil
		}
		log.Errorf("name:%v CountByName error:%v", name, err)
		return count, err
	}
	return
}

func (r *companyRepository) Update(company *model.Company) (*model.Company, error) {
	query := r.session.Query(`
		UPDATE company
		SET name = ?, description = ?, employees = ?, registered = ?, type = ?
		WHERE id = ?
	`, company.Name, company.Description, company.Employees, company.Registered, company.Type, company.ID.String())

	err := query.Exec()

	if err != nil {
		log.Errorf("id:%v Update error:%v", company.ID, err)
		return nil, err
	}

	return r.GetByID(company.ID)
}

func (r *companyRepository) Delete(id uuid.UUID) error {
	query := r.session.Query(`
		DELETE FROM company
		WHERE id = ?
	`, id.String())

	err := query.Exec()
	if err != nil {
		log.Errorf("id:%v Delete error:%v", id, err)
		return err
	}
	return nil
}
