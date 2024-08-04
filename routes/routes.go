package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		logrus.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://meowalien.com", "http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	i18nGroup := r.Group("/i18n")
	{
		i18nGroup.GET("/:lng/:ns", getI18nJSON)
	}

	return r
}
