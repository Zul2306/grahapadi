package controllers

import (
	"inventory-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// In-memory users store (replace with database)
var users = []models.User{
	{ID: 1, Name: "Administrator", Email: "admin@inventory.com", Role: "admin", CreatedAt: time.Now()},
}

// GetUsers returns all users
func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": len(users),
	})
}

// GetUser returns a single user by ID
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	for _, u := range users {
		if u.ID == uint(id) {
			c.JSON(http.StatusOK, gin.H{"data": u})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

// CreateUserRequest holds data for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin staff"`
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check duplicate email
	for _, u := range users {
		if u.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
	}

	newUser := models.User{
		ID:        uint(len(users) + 1),
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Role,
		CreatedAt: time.Now(),
	}
	users = append(users, newUser)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    newUser,
	})
}

// DeleteUser removes a user by ID
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	for i, u := range users {
		if u.ID == uint(id) {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}
