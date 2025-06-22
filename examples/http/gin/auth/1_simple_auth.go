// File: main.go
package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecret")

// GenerateJWT creates a new JWT token for a given user ID
func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": "admin",
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWTMiddleware checks for a valid JWT token in the Authorization header
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["sub"])
			c.Set("role", claims["role"])
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		// In real world: validate user credentials here
		token, err := GenerateJWT("12345")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	protected := r.Group("/api")
	protected.Use(JWTMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.MustGet("user_id").(string)
			role := c.MustGet("role").(string)
			c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role})
		})
	}

	r.Run(":8080")
}
