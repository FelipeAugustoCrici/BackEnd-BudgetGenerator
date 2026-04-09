package corehandler

import (
	"net/http"

	"budgetgen/internal/auth"
	"budgetgen/internal/db"
	"budgetgen/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListTemplates(c *gin.Context) {
	userID := auth.UserID(c)
	var templates []model.Template
	db.DB.Where("user_id = ?", userID).Order("created_at asc").Find(&templates)
	c.JSON(http.StatusOK, templates)
}

func GetTemplate(c *gin.Context) {
	userID := auth.UserID(c)
	var t model.Template
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&t).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template não encontrado"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func CreateTemplate(c *gin.Context) {
	userID := auth.UserID(c)
	var t model.Template
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.ID = uuid.NewString()
	t.UserID = userID
	if err := db.DB.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar template"})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func UpdateTemplate(c *gin.Context) {
	userID := auth.UserID(c)
	var t model.Template
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&t).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template não encontrado"})
		return
	}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.ID = c.Param("id")
	t.UserID = userID
	db.DB.Save(&t)
	c.JSON(http.StatusOK, t)
}

func DeleteTemplate(c *gin.Context) {
	userID := auth.UserID(c)
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).Delete(&model.Template{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template não encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
