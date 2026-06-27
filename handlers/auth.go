package handlers

import (
	"net/http"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type registerInput struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register creates a new user account and returns a JWT.
func Register(c *gin.Context) {
	var in registerInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var existing models.User
	if err := database.DB.Where("email = ?", in.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "email already registered"})
		return
	}

	hashed, err := utils.HashPassword(in.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to hash password"})
		return
	}

	user := models.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: hashed,
		Role:     "user",
		Active:   true,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create user"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email, user.Role)
	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

// Login authenticates a user and returns a JWT.
func Login(c *gin.Context) {
	var in loginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", in.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{"message": "account is disabled"})
		return
	}

	if !utils.CheckPassword(user.Password, in.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email, user.Role)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// Me returns the currently authenticated user.
func Me(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
