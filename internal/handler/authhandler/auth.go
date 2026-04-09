package authhandler

import (
	"net/http"

	"budgetgen/internal/auth"
	"budgetgen/internal/db"
	"budgetgen/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type registerInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input registerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing model.User
	if err := db.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "e-mail já cadastrado"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao processar senha"})
		return
	}

	user := model.User{
		ID:       uuid.NewString(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar usuário"})
		return
	}

	token, _ := auth.GenerateToken(user.ID)
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "name": user.Name, "email": user.Email},
	})
}

func Login(c *gin.Context) {
	var input loginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "e-mail não encontrado"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "senha incorreta"})
		return
	}

	token, _ := auth.GenerateToken(user.ID)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "name": user.Name, "email": user.Email},
	})
}
