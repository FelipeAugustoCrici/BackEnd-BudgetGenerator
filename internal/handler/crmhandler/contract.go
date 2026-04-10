package crmhandler

import (
	"net/http"

	"budgetgen/internal/auth"
	"budgetgen/internal/db"
	"budgetgen/internal/dto"
	"budgetgen/internal/repository"
	"budgetgen/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func contractService() *service.ContractService {
	return service.NewContractService(repository.NewContractRepository(db.DB))
}

func ListContracts(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	var q dto.ContractFilterQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contracts, total, err := contractService().List(userID, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedResponse{Data: contracts, Total: total, Page: q.Page, Limit: q.Limit})
}

func GetContract(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	contract, err := contractService().GetByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
		return
	}
	c.JSON(http.StatusOK, contract)
}

func CreateContract(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	var req dto.CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contract, err := contractService().Create(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, contract)
}

func UpdateContract(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req dto.UpdateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contract, err := contractService().Update(id, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contract)
}

func SendContract(c *gin.Context) {
	action(c, "send")
}

func ViewContract(c *gin.Context) {
	action(c, "view")
}

func SignContract(c *gin.Context) {
	action(c, "sign")
}

func RefuseContract(c *gin.Context) {
	action(c, "refuse")
}

func GetContractByBudget(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	budgetID, err := uuid.Parse(c.Param("budgetId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget id"})
		return
	}
	contract, err := contractService().GetByBudgetID(budgetID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
		return
	}
	c.JSON(http.StatusOK, contract)
}

func ListContractEvents(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	events, err := contractService().ListEvents(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
		return
	}
	c.JSON(http.StatusOK, events)
}

func action(c *gin.Context, act string) {
	userID := uuid.MustParse(auth.UserID(c))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := contractService()
	var contract interface{}
	switch act {
	case "send":
		contract, err = svc.Send(id, userID)
	case "view":
		contract, err = svc.View(id, userID)
	case "sign":
		contract, err = svc.Sign(id, userID)
	case "refuse":
		contract, err = svc.Refuse(id, userID)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contract)
}
