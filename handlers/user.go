package handlers

import (
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// ListUsers returns a paginated, searchable list of users.
func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	query := database.DB.Model(&models.User{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", like, like)
	}

	var total int64
	query.Count(&total)

	var users []models.User
	offset := (page - 1) * limit
	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetUser returns a single user by ID.
func GetUser(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

type createUserInput struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
	Active   *bool  `json:"active"`
	StaffID  string `json:"staff_id"`
}

// CreateUser adds a new user.
func CreateUser(c *gin.Context) {
	var in createUserInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var existing models.User
	if err := database.DB.Where("email = ?", in.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "email already exists"})
		return
	}

	hashed, _ := utils.HashPassword(in.Password)
	role := in.Role
	if role == "" {
		role = "user"
	}
	active := true
	if in.Active != nil {
		active = *in.Active
	}

	user := models.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: hashed,
		Role:     role,
		Active:   active,
		StaffID:  in.StaffID,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

type updateUserInput struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	Active   *bool   `json:"active"`
	StaffID  *string `json:"staff_id"` // pointer: omit = unchanged, "" = unlink
}

// UpdateUser modifies an existing user.
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	var in updateUserInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if in.Name != "" {
		user.Name = in.Name
	}
	if in.Email != "" {
		user.Email = in.Email
	}
	if in.Role != "" {
		user.Role = in.Role
	}
	if in.Active != nil {
		user.Active = *in.Active
	}
	if in.StaffID != nil {
		user.StaffID = *in.StaffID
	}
	if in.Password != "" {
		hashed, _ := utils.HashPassword(in.Password)
		user.Password = hashed
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser removes a user by ID.
func DeleteUser(c *gin.Context) {
	if err := database.DB.Delete(&models.User{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// Stats returns simple dashboard metrics.
func Stats(c *gin.Context) {
	var total, active, admins int64
	database.DB.Model(&models.User{}).Count(&total)
	database.DB.Model(&models.User{}).Where("active = ?", true).Count(&active)
	database.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&admins)

	c.JSON(http.StatusOK, gin.H{
		"total_users":  total,
		"active_users": active,
		"admins":       admins,
		"inactive":     total - active,
	})
}
