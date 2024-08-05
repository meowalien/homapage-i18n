package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ParseUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		jet_claim, exist := c.Get("jwt_claims")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "jwt_claims not found"})
			c.Abort()
			return
		}

		claims, ok := jet_claim.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "jwt_claims is invalid"})
			c.Abort()
			return
		}
		userString := claims["user_id"]
		if userString == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in jwt_claims"})
			c.Abort()
			return
		}

		userID, ok := userString.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in jwt_claims"})
			c.Abort()
		}
		c.Set("user_id", userID)
		c.Next()
	}
}
