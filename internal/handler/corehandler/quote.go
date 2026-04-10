package corehandler

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
	createQuoteVersion(quote, userID, "Versão inicial")
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
	createQuoteVersion(quote, userID, "")
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

func ListQuoteVersions(c *gin.Context) {
	userID := auth.UserID(c)
	quoteID := c.Param("id")
	var versions []model.QuoteVersion
	db.DB.Where("quote_id = ? AND user_id = ?", quoteID, userID).
		Order("version_number desc").Find(&versions)
	c.JSON(http.StatusOK, versions)
}

func ActivateQuoteVersion(c *gin.Context) {
	userID := auth.UserID(c)
	quoteID := c.Param("id")
	versionID := c.Param("versionId")

	var version model.QuoteVersion
	if err := db.DB.Where("id = ? AND quote_id = ? AND user_id = ?", versionID, quoteID, userID).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "versão não encontrada"})
		return
	}

	// deactivate all, then activate selected
	db.DB.Model(&model.QuoteVersion{}).
		Where("quote_id = ? AND user_id = ?", quoteID, userID).
		Update("is_active", false)
	db.DB.Model(&version).Update("is_active", true)

	// restore quote from snapshot
	snap := version.Snapshot
	db.DB.Model(&model.Quote{}).Where("id = ? AND user_id = ?", quoteID, userID).Updates(map[string]interface{}{
		"client_name":   snap.ClientName,
		"date":          snap.Date,
		"notes":         snap.Notes,
		"status":        snap.Status,
		"items":         snap.Items,
		"discount":      snap.Discount,
		"discount_type": snap.DiscountType,
		"template_id":   snap.TemplateID,
		"scope":         snap.Scope,
		"conditions":    snap.Conditions,
		"hourly_rate":   snap.HourlyRate,
		"company_name":  snap.CompanyName,
	})

	var updated model.Quote
	db.DB.Where("id = ?", quoteID).First(&updated)
	c.JSON(http.StatusOK, updated)
}

func createQuoteVersion(quote model.Quote, userID, changeNote string) {
	var count int64
	db.DB.Model(&model.QuoteVersion{}).Where("quote_id = ?", quote.ID).Count(&count)

	db.DB.Model(&model.QuoteVersion{}).
		Where("quote_id = ? AND is_active = true", quote.ID).
		Update("is_active", false)

	version := model.QuoteVersion{
		ID:            uuid.NewString(),
		QuoteID:       quote.ID,
		UserID:        userID,
		VersionNumber: int(count) + 1,
		IsActive:      true,
		ChangeNote:    changeNote,
		Snapshot: model.QuoteSnapshot{
			ClientName:   quote.ClientName,
			Date:         quote.Date,
			Notes:        quote.Notes,
			Status:       quote.Status,
			Items:        quote.Items,
			Discount:     quote.Discount,
			DiscountType: quote.DiscountType,
			TemplateID:   quote.TemplateID,
			Scope:        quote.Scope,
			Conditions:   quote.Conditions,
			HourlyRate:   quote.HourlyRate,
			CompanyName:  quote.CompanyName,
		},
	}
	db.DB.Create(&version)
}
