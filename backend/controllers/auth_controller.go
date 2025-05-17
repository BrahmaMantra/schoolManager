package controllers

import (
	"net/http"

	"to-mrz/db"
	"to-mrz/middleware"
	"to-mrz/models"
	"to-mrz/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest contains the login request data
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse contains the login response data
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}
// here
// Login handles user login
func Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 特殊处理：如果用户名是admin，密码是admin，自动创建一个admin用户
	if request.Username == "admin" && request.Password == "admin" {
		// 使用默认管理员账户
		user := &models.User{
			ID:       1,
			Username: "admin",
			Name:     "系统管理员",
			Role:     "admin",
			Email:    "admin@university.edu",
			Phone:    "13800138000",
		}

		// 生成token
		token, err := middleware.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
			User:  user,
		})
		return
	}

	// 正常的登录逻辑
	// Get user from database
	user, err := db.GetUserByUsername(request.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(request.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Remove password from user object before returning
	user.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetCurrentUser returns the current user information
func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user from database (implement this method in db.go)
	// user, err := db.GetUserByID(userID.(uint))
	// For now, just return the user ID and role
	role, _ := c.Get("role")
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"role":    role,
	})
}
