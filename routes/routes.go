package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"homapage-i18n/role"
	"homapage-i18n/routes/handlers"
	"homapage-i18n/routes/middleware"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		logrus.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Configure CORS
	r.Use(middleware.Cors())

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	i18nGroup := r.Group("/i18n")
	{
		i18nGroup.GET("/:lng/:ns", handlers.GetI18n())
		i18nGroup.POST("/:lng/:ns", middleware.ParseToken(), middleware.CheckTokenRole(role.Admin), middleware.ParseUserID(), handlers.PostI18n())
		i18nGroup.DELETE("/:lng/:ns", middleware.ParseToken(), middleware.CheckTokenRole(role.Admin), middleware.ParseUserID(), handlers.DeleteI18n())
	}

	return r
}
