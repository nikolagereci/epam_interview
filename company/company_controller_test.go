package company

import (
	"encoding/json"
	"errors"
	mock_company_service "github.com/ngereci/xm_interview/mocks/mock_company/service"
	"github.com/ngereci/xm_interview/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestController_CreateCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	mockController := NewController(mockService)

	// Test case: Successful creation
	newCompany := &model.Company{Name: "Test Company", Employees: 100, Type: model.Corporation}
	expectedCompany := &model.Company{Name: "Test Company", ID: uuid.New(), Type: model.Corporation}
	mockService.EXPECT().CreateCompany(newCompany).Return(expectedCompany, nil)
	// Create a test user
	requestBody, _ := json.Marshal(newCompany)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r

	mockController.CreateCompany(ctx)
	expectedJsonString, _ := json.Marshal(expectedCompany)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, string(expectedJsonString), w.Body.String())
}

func TestController_CreateCompany_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	mockController := NewController(mockService)

	newCompany := &model.Company{Name: "Test Company", Employees: 100, Type: model.Corporation}
	// Test case: Failed creation due to service error
	mockService.EXPECT().CreateCompany(newCompany).Return(nil, errors.New("Test error"))
	requestBody, _ := json.Marshal(newCompany)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r

	mockController.CreateCompany(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"Test error"}`, w.Body.String())
}

func TestController_CreateCompany_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	mockController := NewController(mockService)

	newCompany := &model.Company{Name: "Test Company", Employees: 100}
	// Test case 1: validation fail on type
	requestBody, _ := json.Marshal(newCompany)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r

	mockController.CreateCompany(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Key: 'Company.Type' Error:Field validation for 'Type' failed on the 'required' tag"}`, w.Body.String())

	// Test case 2: validation fail on employees
	newCompany = &model.Company{Name: "Test Company", Type: model.Corporation}
	requestBody, _ = json.Marshal(newCompany)
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = r

	mockController.CreateCompany(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Key: 'Company.Employees' Error:Field validation for 'Employees' failed on the 'required' tag"}`, w.Body.String())
	// Test case 3: validation fail on name
	newCompany = &model.Company{Employees: 100, Type: model.Corporation}
	requestBody, _ = json.Marshal(newCompany)
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = r

	mockController.CreateCompany(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Key: 'Company.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`, w.Body.String())
}

func TestController_GetCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	controller := NewController(mockService)

	companyID := uuid.New()
	dummyCompany := &model.Company{
		ID:   companyID,
		Name: "Test Company",
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = gin.Params{{Key: "id", Value: companyID.String()}}
	mockService.EXPECT().GetCompanyByID(companyID).Return(dummyCompany, nil)
	controller.GetCompany(ctx)
	expectedJsonString, _ := json.Marshal(dummyCompany)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expectedJsonString), w.Body.String())
}

func TestController_GetCompany_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	controller := NewController(mockService)

	companyID := uuid.New()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = gin.Params{{Key: "id", Value: companyID.String()}}
	mockService.EXPECT().GetCompanyByID(companyID).Return(nil, nil)
	controller.GetCompany(ctx)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "company not found")
}

func TestController_GetCompany_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	controller := NewController(mockService)

	companyID := uuid.New()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = gin.Params{{Key: "id", Value: companyID.String()}}
	mockService.EXPECT().GetCompanyByID(companyID).Return(nil, errors.New("something went wrong"))
	controller.GetCompany(ctx)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "something went wrong")
}

func TestUpdateCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	mockController := NewController(mockService)

	// Test case: Successful update
	companyID := uuid.New()
	newCompany := &model.Company{Name: "Test Company", Employees: 100, Type: model.Corporation}
	expectedCompany := &model.Company{Name: "Test Company", ID: companyID, Type: model.Corporation}
	mockService.EXPECT().UpdateCompany(companyID, gomock.Any()).Return(&model.Company{ID: companyID, Name: "Test Company", Type: model.Corporation}, nil).Times(1)
	// Create a test user
	requestBody, _ := json.Marshal(newCompany)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(requestBody)))
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = gin.Params{{Key: "id", Value: companyID.String()}}

	mockController.UpdateCompany(ctx)
	expectedJsonString, _ := json.Marshal(expectedCompany)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(expectedJsonString), w.Body.String())

}

func TestUpdateCompany_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	mockController := NewController(mockService)

	// Test case: Successful update
	companyID := ""
	newCompany := &model.Company{Name: "Test Company", Employees: 100, Type: model.Corporation}
	//mockService.EXPECT().UpdateCompany(gomock.Any(), companyID, gomock.Any()).Return(&Company{ID: companyID, Name: "Test Company", Type: Corporation}, nil).Times(1)
	// Create a test user
	requestBody, _ := json.Marshal(newCompany)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(requestBody)))
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = gin.Params{{Key: "id", Value: companyID}}

	mockController.UpdateCompany(ctx)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(t, `{"error":"invalid UUID length: 0"}`, w.Body.String())

}

func TestUpdateCompany_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	mockController := NewController(mockService)

	// Test case: Successful update
	companyID := uuid.New()
	newCompany := &model.Company{Name: "Test Company", Employees: 100, Type: model.Corporation}
	mockService.EXPECT().UpdateCompany(companyID, gomock.Any()).Return(nil, errors.New("something went wrong")).Times(1)
	// Create a test user
	requestBody, _ := json.Marshal(newCompany)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(requestBody)))
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = gin.Params{{Key: "id", Value: companyID.String()}}

	mockController.UpdateCompany(ctx)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"something went wrong"}`, w.Body.String())

}
func TestUpdateCompany_CompanyNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	mockController := NewController(mockService)

	// Test case: Successful update
	companyID := uuid.New()
	newCompany := &model.Company{Name: "Test Company", Employees: 100, Type: model.Corporation}
	mockService.EXPECT().UpdateCompany(companyID, gomock.Any()).Return(nil, nil).Times(1)
	// Create a test user
	requestBody, _ := json.Marshal(newCompany)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(requestBody)))
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = gin.Params{{Key: "id", Value: companyID.String()}}

	mockController.UpdateCompany(ctx)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"company not found"}`, w.Body.String())

}

func TestDeleteCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	mockController := NewController(mockService)

	// Test case: Successful update
	companyID := uuid.New()
	mockService.EXPECT().DeleteCompany(companyID).Return(nil).Times(1)
	// Create a test user
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = gin.Params{{Key: "id", Value: companyID.String()}}

	mockController.DeleteCompany(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, "{}", w.Body.String())

}
func TestDeleteCompany_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_company_service.NewMockService(ctrl)
	mockController := NewController(mockService)

	// Test case: Successful update
	companyID := uuid.New()
	mockService.EXPECT().DeleteCompany(companyID).Return(errors.New("something went wrong")).Times(1)
	// Create a test user
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = gin.Params{{Key: "id", Value: companyID.String()}}

	mockController.DeleteCompany(ctx)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"something went wrong"}`, w.Body.String())

}

func TestProcessUuid_Success(t *testing.T) {
	// Prepare test case
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uuid.New().String()})

	// Execute function
	result, err := processUuid(ctx)

	// Assert result
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestProcessUuid_Error(t *testing.T) {
	// Prepare test case
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "invalid"})

	// Execute function
	result, err := processUuid(ctx)

	// Assert result
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
