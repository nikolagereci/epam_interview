package company

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ngereci/xm_interview/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Controller interface {
	CreateCompany(ctx *gin.Context)
	GetCompany(ctx *gin.Context)
	UpdateCompany(ctx *gin.Context)
	DeleteCompany(ctx *gin.Context)
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service: service}
}

func (c *controller) CreateCompany(ctx *gin.Context) {
	var company model.Company

	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdCompany, err := c.service.CreateCompany(&company)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdCompany)
}

func (c *controller) GetCompany(ctx *gin.Context) {
	companyUuid, err := processUuid(ctx)
	if err != nil {
		return
	}
	company, err := c.service.GetCompanyByID(*companyUuid)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if company == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	ctx.JSON(http.StatusOK, company)
}

func (c *controller) UpdateCompany(ctx *gin.Context) {
	companyUuid, err := processUuid(ctx)
	if err != nil {
		return
	}
	var company model.Company
	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedCompany, err := c.service.UpdateCompany(*companyUuid, &company)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updatedCompany == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	ctx.JSON(http.StatusOK, updatedCompany)
}

func (c *controller) DeleteCompany(ctx *gin.Context) {
	companyUuid, err := processUuid(ctx)
	if err != nil {
		return
	}
	err = c.service.DeleteCompany(*companyUuid)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func processUuid(ctx *gin.Context) (*uuid.UUID, error) {
	id := ctx.Param("id")
	companyUuid, err := uuid.Parse(id)
	if err != nil {
		log.Warnf("id:%v UUID parse error:%v", id, err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return nil, err
	}
	return &companyUuid, nil
}
