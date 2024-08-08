package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"homapage-i18n/mongodb"
	"net/http"
)

func DeleteI18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		lng := c.Param("lng")
		ns := c.Param("ns")

		collectionName := fmt.Sprintf("%s-%s", lng, ns)
		collection := mongodb.GetCollection("i18n", collectionName)

		err := collection.Drop(context.TODO())
		if err != nil {
			logrus.Errorf("Failed to drop collection %s: %v", collectionName, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Collection deleted"})
	}
}
