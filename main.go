package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
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

	r.Run(":8080")
}

func getI18nJSON(c *gin.Context) {
	lng := c.Param("lng")
	ns := c.Param("ns")

	// Debugging output to verify parameters
	fmt.Printf("lng: %s, ns: %s\n", lng, ns)

	filePath := fmt.Sprintf("i18n/%s/%s.json", lng, ns)
	fmt.Println("filePath: ", filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}

	c.JSON(http.StatusOK, jsonData)
}
