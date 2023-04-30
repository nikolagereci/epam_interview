package company

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ngereci/xm_interview/event"
	mock_kafka "github.com/ngereci/xm_interview/mocks/mock_company/event"
	mock_company_repository "github.com/ngereci/xm_interview/mocks/mock_company/repository"
	"github.com/ngereci/xm_interview/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testCompany = &model.Company{
		ID:          uuid.MustParse("56f86115-a58f-43db-8a1b-9aa2908f7a18"),
		Name:        "Test Company",
		Description: "Test Description",
		Employees:   100,
		Registered:  false,
		Type:        model.Corporation,
	}
	testCompanyUpdate = &model.Company{
		ID:          uuid.MustParse("56f86115-a58f-43db-8a1b-9aa2908f7a18"),
		Name:        "Test Company Update",
		Description: "Test Description Update",
		Employees:   200,
		Registered:  true,
		Type:        model.Corporation,
	}
	testErr = errors.New("test error")
)

func TestCompanyService_CreateCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	newCompany := &model.Company{
		Name: "Test Company",
	}

	mockRepo.EXPECT().CountByName(newCompany.Name).Return(0, nil)
	mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(company *model.Company) error {
		assert.Equal(t, newCompany.Name, company.Name)
		assert.NotEqual(t, uuid.Nil, company.ID)
		*testCompany = *company
		return nil
	})
	mockKafka.EXPECT().SendEventWithPayload(event.EVENT_CREATE, testCompany).Return(nil)

	svc := NewService(mockRepo, mockKafka)
	company, err := svc.CreateCompany(newCompany)

	assert.NoError(t, err)
	assert.Equal(t, testCompany, company)
}

func TestCompanyService_CreateCompany_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	newCompany := &model.Company{
		Name: "Test Company",
	}

	mockRepo.EXPECT().CountByName(newCompany.Name).Return(1, nil)
	svc := NewService(mockRepo, mockKafka)
	_, err := svc.CreateCompany(newCompany)
	assert.Error(t, err)
	assert.IsType(t, model.ErrCompanyExists{}, err)

	mockRepo.EXPECT().CountByName(newCompany.Name).Return(0, testErr)

	_, err = svc.CreateCompany(newCompany)
	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

func TestCompanyService_CreateCompany_CreateFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	newCompany := &model.Company{
		Name: "Test Company",
	}

	mockRepo.EXPECT().CountByName(newCompany.Name).Return(0, nil)
	mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(company *model.Company) error {
		assert.Equal(t, newCompany.Name, company.Name)
		assert.NotEqual(t, uuid.Nil, company.ID)
		*testCompany = *company
		return testErr
	})
	//mockKafka.EXPECT().SendEventWithPayload(event.EVENT_CREATE, testCompany).Return(testErr)
	svc := NewService(mockRepo, mockKafka)
	_, err := svc.CreateCompany(newCompany)
	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

func TestCompanyService_CreateCompany_KafkaFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	newCompany := &model.Company{
		Name: "Test Company",
	}

	mockRepo.EXPECT().CountByName(newCompany.Name).Return(0, nil)
	mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(company *model.Company) error {
		assert.Equal(t, newCompany.Name, company.Name)
		assert.NotEqual(t, uuid.Nil, company.ID)
		*testCompany = *company
		return nil
	})
	mockRepo.EXPECT().Delete(gomock.Any()).Return(nil)
	mockKafka.EXPECT().SendEventWithPayload(event.EVENT_CREATE, testCompany).Return(testErr)
	svc := NewService(mockRepo, mockKafka)
	_, err := svc.CreateCompany(newCompany)
	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

func TestGetCompanyByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafkaProducer := mock_kafka.NewMockKafkaAdapter(ctrl)

	mockRepo.EXPECT().GetByID(testCompany.ID).Return(testCompany, nil)

	companyService := NewService(mockRepo, mockKafkaProducer)
	company, err := companyService.GetCompanyByID(testCompany.ID)

	assert.NoError(t, err)
	assert.Equal(t, testCompany, company)
}

func TestGetCompanyByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafkaProducer := mock_kafka.NewMockKafkaAdapter(ctrl)

	mockRepo.EXPECT().GetByID(testCompany.ID).Return(nil, errors.New("something went wrong"))

	companyService := NewService(mockRepo, mockKafkaProducer)
	company, err := companyService.GetCompanyByID(testCompany.ID)

	assert.Error(t, err)
	assert.Nil(t, company)
}

func TestCompanyService_UpdateCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	mockRepo.EXPECT().GetByID(testCompany.ID).Return(testCompany, nil)
	mockRepo.EXPECT().Update(gomock.Any()).Return(testCompanyUpdate, nil)
	mockKafka.EXPECT().SendEventWithPayload(event.EVENT_UPDATE, testCompanyUpdate).Return(nil)

	svc := NewService(mockRepo, mockKafka)
	company, err := svc.UpdateCompany(testCompany.ID, testCompanyUpdate)

	assert.NoError(t, err)
	assert.Equal(t, testCompanyUpdate, company)
}

func TestCompanyService_UpdateCompany_UpdateFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	mockRepo.EXPECT().GetByID(testCompany.ID).Return(testCompany, nil)
	mockRepo.EXPECT().Update(gomock.Any()).Return(nil, errors.New("something went wrong"))

	svc := NewService(mockRepo, mockKafka)
	company, err := svc.UpdateCompany(testCompany.ID, testCompanyUpdate)

	assert.Error(t, err)
	assert.Nil(t, company)
}

func TestCompanyService_UpdateCompany_KafkaFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	mockRepo.EXPECT().GetByID(testCompany.ID).Return(testCompany, nil)
	//mockRepo.EXPECT().CountByName(testCompany.Name).Return(0, nil)
	mockRepo.EXPECT().Update(gomock.Any()).Return(testCompanyUpdate, nil).Times(2)
	mockKafka.EXPECT().SendEventWithPayload(event.EVENT_UPDATE, testCompanyUpdate).Return(errors.New("something went wrong"))

	svc := NewService(mockRepo, mockKafka)
	company, err := svc.UpdateCompany(testCompany.ID, testCompanyUpdate)

	assert.Error(t, err)
	assert.Nil(t, company)
}
func TestCompanyService_DeleteCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	mockRepo.EXPECT().GetByID(testCompany.ID).Return(testCompany, nil)
	//mockRepo.EXPECT().CountByName(testCompany.Name).Return(0, nil)
	mockRepo.EXPECT().Delete(gomock.Any()).Return(nil)
	mockKafka.EXPECT().SendEventWithPayload(event.EVENT_DELETE, testCompany).Return(nil)

	svc := NewService(mockRepo, mockKafka)
	err := svc.DeleteCompany(testCompany.ID)

	assert.NoError(t, err)
}

func TestCompanyService_DeleteCompany_DeleteFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	mockRepo.EXPECT().GetByID(testCompany.ID).Return(testCompany, nil)
	mockRepo.EXPECT().Delete(gomock.Any()).Return(errors.New("something went wrong"))

	svc := NewService(mockRepo, mockKafka)
	err := svc.DeleteCompany(testCompany.ID)

	assert.Error(t, err)
}

func TestCompanyService_DeleteCompany_KafkaFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_company_repository.NewMockRepository(ctrl)
	mockKafka := mock_kafka.NewMockKafkaAdapter(ctrl)

	mockRepo.EXPECT().GetByID(testCompany.ID).Return(testCompany, nil)
	mockRepo.EXPECT().Delete(gomock.Any()).Return(nil)
	mockRepo.EXPECT().Create(testCompany).Return(nil)
	mockKafka.EXPECT().SendEventWithPayload(event.EVENT_DELETE, testCompany).Return(errors.New("something went wrong"))

	svc := NewService(mockRepo, mockKafka)
	err := svc.DeleteCompany(testCompany.ID)

	assert.Error(t, err)
}
