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

func clientService() *service.ClientService {
	return service.NewClientService(repository.NewClientRepository(db.DB))
}

func ListClients(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	var q dto.ClientFilterQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clients, total, err := clientService().List(userID, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedResponse{Data: clients, Total: total, Page: q.Page, Limit: q.Limit})
}

func GetClient(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	client, err := clientService().GetByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

func CreateClient(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	var req dto.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client, err := clientService().Create(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, client)
}

func UpdateClient(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req dto.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client, err := clientService().Update(id, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, client)
}

func DeleteClient(c *gin.Context) {
	userID := uuid.MustParse(auth.UserID(c))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := clientService().Delete(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
