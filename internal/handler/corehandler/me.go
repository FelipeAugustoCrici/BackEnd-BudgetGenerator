package corehandler

import (
	"net/http"

	"budgetgen/internal/auth"
	"budgetgen/internal/db"
	"budgetgen/internal/model"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	userID := auth.UserID(c)
	var user model.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "name": user.Name, "email": user.Email})
}
