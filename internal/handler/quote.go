package handler

import (
	"net/http"

	"budgetgen/internal/auth"
	"budgetgen/internal/db"
	"budgetgen/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListQuotes(c *gin.Context) {
	userID := auth.UserID(c)
	var quotes []model.Quote
	db.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&quotes)
	c.JSON(http.StatusOK, quotes)
}

func GetQuote(c *gin.Context) {
	userID := auth.UserID(c)
	var quote model.Quote
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&quote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "orçamento não encontrado"})
		return
	}
	c.JSON(http.StatusOK, quote)
}

func CreateQuote(c *gin.Context) {
	userID := auth.UserID(c)
	var quote model.Quote
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quote.ID = uuid.NewString()
	quote.UserID = userID
	if err := db.DB.Create(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar orçamento"})
		return
	}
	c.JSON(http.StatusCreated, quote)
}

func UpdateQuote(c *gin.Context) {
	userID := auth.UserID(c)
	var quote model.Quote
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&quote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "orçamento não encontrado"})
		return
	}
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quote.ID = c.Param("id")
	quote.UserID = userID
	db.DB.Save(&quote)
	c.JSON(http.StatusOK, quote)
}

func DeleteQuote(c *gin.Context) {
	userID := auth.UserID(c)
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).Delete(&model.Quote{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "orçamento não encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
