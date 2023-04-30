package company

import (
	"context"
	"github.com/google/uuid"
	"github.com/ngereci/xm_interview/event"
	"github.com/ngereci/xm_interview/model"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	CreateCompany(ctx context.Context, newCompany *model.Company) (*model.Company, error)
	GetCompanyByID(ctx context.Context, id uuid.UUID) (*model.Company, error)
	UpdateCompany(ctx context.Context, id uuid.UUID, forUpdateCompany *model.Company) (*model.Company, error)
	DeleteCompany(ctx context.Context, id uuid.UUID) error
}

type companyService struct {
	repo          Repository
	kafkaProducer event.KafkaAdapter
}

func NewService(repo Repository, kafkaProducer event.KafkaAdapter) Service {
	return &companyService{repo: repo, kafkaProducer: kafkaProducer}
}

func (s *companyService) CreateCompany(ctx context.Context, newCompany *model.Company) (*model.Company, error) {
	// Generate a new UUID for the company
	newCompany.ID = uuid.New()
	// determine new company name is unique
	count, err := s.repo.CountByName(ctx, newCompany.Name)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, model.ErrCompanyExists{newCompany.Name}
	}
	err = s.repo.Create(ctx, newCompany)
	if err != nil {
		return nil, err
	}
	if kafkaErr := s.kafkaProducer.SendEventWithPayload(event.EVENT_CREATE, newCompany); kafkaErr != nil {
		log.Errorf("company:%v created but send event failed, rolling back. error:%v", newCompany.ID, err)
		//handle rollback
		if err = s.repo.Delete(ctx, newCompany.ID); err != nil {
			log.Errorf("company:%v rollback failed. error:%v", newCompany.ID, err)
			return nil, err
		}
		log.Infof("company:%v rollback success", newCompany.ID)
		return nil, kafkaErr
	}
	return newCompany, nil
}

func (s *companyService) GetCompanyByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *companyService) UpdateCompany(ctx context.Context, id uuid.UUID, forUpdateCompany *model.Company) (*model.Company, error) {
	existingCompany, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if existingCompany == nil {
		return nil, model.ErrCompanyNotFound{Id: id}
	}

	// Copy over the fields that can't be updated
	forUpdateCompany.ID = existingCompany.ID

	updatedCompany, err := s.repo.Update(ctx, forUpdateCompany)

	if err != nil {
		return nil, err
	}
	if kafkaErr := s.kafkaProducer.SendEventWithPayload(event.EVENT_UPDATE, forUpdateCompany); kafkaErr != nil {
		log.Errorf("forUpdateCompany:%v updated but send event failed, rolling back. error:%v", forUpdateCompany.ID, err)
		//handle rollback
		if _, err = s.repo.Update(ctx, existingCompany); err != nil {
			log.Errorf("forUpdateCompany:%v rollback failed. error:%v", forUpdateCompany.ID, err)
			return nil, err
		}
		log.Infof("forUpdateCompany:%v rollback success", forUpdateCompany.ID)
		return nil, kafkaErr
	}

	return updatedCompany, nil
}

func (s *companyService) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	existingCompany, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return err
	}

	if existingCompany == nil {
		return model.ErrCompanyNotFound{id}
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if kafkaErr := s.kafkaProducer.SendEventWithPayload(event.EVENT_DELETE, existingCompany); kafkaErr != nil {
		log.Errorf("company:%v deleted but send event failed, rolling back. error:%v", existingCompany.ID, err)
		//handle rollback
		if err = s.repo.Create(ctx, existingCompany); err != nil {
			log.Errorf("company:%v rollback failed. error:%v", existingCompany.ID, err)
			return err
		}
		log.Infof("company:%v rollback success", existingCompany.ID)
		return kafkaErr
	}
	return nil
}
