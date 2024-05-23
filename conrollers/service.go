package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"practise/config"
	"practise/logs"
	"practise/models"
	"strings"
	"time"
)

func CreateUser(c *gin.Context) {

	var req models.User
	if err := c.ShouldBind(&req); err != nil {
		logs.LogToFile("Error binding request: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.Name) == "" {
		// Log the error to a file
		logs.LogToFile("All fields are required")
		// Return a response with status Bad Request (400)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "All fields are required"})
		return
	}

	var existingUser models.User
	// Check if a user with the same email already exists
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		// If user with the same email already exists, return error
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Email already exists",
		})
		return
	}
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return // This stops further execution
	}
	req.Password = hashedPassword

	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return // This stops further execution
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User saved successfully",
		"data":    req})
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func DeleteUser(c *gin.Context) {
	// Get the user ID from the URL parameter
	userID := c.Param("id")

	// Query the database to find the user by ID
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete the user from the database
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully",
	})
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type LoginRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

var jwtSecret = []byte("2022@abu")

func generateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func LoginUser(c *gin.Context) {
	var req LoginRequest

	// Bind form request to struct
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	// Check if a user with the given email exists
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Check if the provided password matches the stored hashed password
	if !checkPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Successful login
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged in successfully",
		"token":   token,
	})
}
