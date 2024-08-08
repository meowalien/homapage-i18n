package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"homapage-i18n/role"
	"net/http"
)

func CheckTokenRole(expectedRole role.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exist := c.Get("jwt_claims")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "jwt_claims not found"})
			c.Abort()
			return
		}
		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "jwt_claims is invalid"})
			c.Abort()
			return
		}

		roleInClaims, exist := mapClaims["roles"]
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found in jwt_claims"})
			c.Abort()
			return
		}

		if roleArray, okRoleArray := roleInClaims.([]string); okRoleArray {
			logrus.Debugf("roleArray: %v", roleArray)
			for _, r := range roleArray {
				if r == expectedRole {
					c.Next()
					return
				}
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role is invalid"})
			c.Abort()
			return
		}
		c.Next()
	}
}
