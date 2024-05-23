package router

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/golang-jwt/jwt"
	controllers "practise/conrollers"
	"strings"
)

// JWT secret key (ensure this is securely stored)
var jwtSecret = []byte("2022@abu")

// AuthMiddleware Middleware to authenticate requests using JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user email from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		email := claims["email"].(string)
		c.Set("email", email)

		c.Next()
	}
}
func ProtectedEndpoint(c *gin.Context) {
	email, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"email":  email,
	})
}

func MapUrl(r *gin.Engine) {
	r.POST("/user/register", controllers.CreateUser)
	r.POST("/user/login", controllers.LoginUser)
	r.DELETE("/user/:id", controllers.DeleteUser)

	authorized := r.Group("/").Use(AuthMiddleware())
	authorized.GET("/protected", ProtectedEndpoint)
}
