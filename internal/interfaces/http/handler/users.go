package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 临时处理器（待实现具体业务逻辑）
func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - to be implemented"})
}

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Register endpoint - to be implemented"})
}

func RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token endpoint - to be implemented"})
}

func ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List users endpoint - to be implemented"})
}

func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get user endpoint - to be implemented"})
}

func UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update user endpoint - to be implemented"})
}

func DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete user endpoint - to be implemented"})
}
