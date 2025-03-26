package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func ValidateToken(tokenString string) (string, bool, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	fmt.Println(secretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return "", false, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiry, ok := claims["exp"].(float64)
		if !ok {
			return "", false, fmt.Errorf("missing expiry claim")
		}
		if time.Unix(int64(expiry), 0).Before(time.Now()) {
			return "", false, fmt.Errorf("token expired")
		}

		userID, ok := claims["staff_id"].(float64)
		if !ok {
			return "", false, fmt.Errorf("invalid token claims")
		}

		return strconv.FormatFloat(userID, 'f', 2, 64), true, nil
	}

	return "", false, fmt.Errorf("invalid token")
}

func TokenValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		userID, valid, err := ValidateToken(tokenString)
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
