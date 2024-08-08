package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"homapage-i18n/token"
	"net/http"
)

func ParseToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			logrus.Errorf("Failed to get token from cookie: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
			c.Abort()
			return
		}
		claims, err := token.VerifyToken(tokenString)
		if err != nil {
			logrus.Errorf("Failed to verify token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
			c.Abort()
			return
		}
		logrus.Debug("claims: ", claims)
		c.Set("jwt_claims", claims)
		c.Next()
	}
}
