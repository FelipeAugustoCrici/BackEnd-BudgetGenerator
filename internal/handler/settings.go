package handler

import (
	"net/http"

	"budgetgen/internal/auth"
	"budgetgen/internal/db"
	"budgetgen/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetSettings(c *gin.Context) {
	userID := auth.UserID(c)
	var s model.Settings
	err := db.DB.Where("user_id = ?", userID).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		// return empty settings
		c.JSON(http.StatusOK, model.Settings{UserID: userID})
		return
	}
	c.JSON(http.StatusOK, s)
}

func UpsertSettings(c *gin.Context) {
	userID := auth.UserID(c)
	var input model.Settings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var s model.Settings
	err := db.DB.Where("user_id = ?", userID).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		input.ID = uuid.NewString()
		input.UserID = userID
		db.DB.Create(&input)
		c.JSON(http.StatusCreated, input)
		return
	}

	input.ID = s.ID
	input.UserID = userID
	db.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}
