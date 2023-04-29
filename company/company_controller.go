package company

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) CreateCompany(ctx *gin.Context) {
	var company Company

	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCompany, err := c.service.CreateCompany(ctx.Request.Context(), &company)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdCompany)
}

func (c *Controller) GetCompany(ctx *gin.Context) {
	companyUuid, err := processUuid(ctx)
	if err != nil {
		return
	}
	company, err := c.service.GetCompanyByID(ctx.Request.Context(), *companyUuid)

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

func (c *Controller) UpdateCompany(ctx *gin.Context) {
	companyUuid, err := processUuid(ctx)
	if err != nil {
		return
	}
	var company Company
	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCompany, err := c.service.UpdateCompany(ctx.Request.Context(), *companyUuid, &company)

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

func (c *Controller) DeleteCompany(ctx *gin.Context) {
	companyUuid, err := processUuid(ctx)
	if err != nil {
		return
	}
	err = c.service.DeleteCompany(ctx.Request.Context(), *companyUuid)

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
