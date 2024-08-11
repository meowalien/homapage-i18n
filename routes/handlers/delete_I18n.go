package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"homapage-i18n/mongodb"
	"net/http"
)

func DeleteI18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		lng := c.Param("lng")
		ns := c.Param("ns")

		collection := mongodb.GetCollection("homepage", "i18n")

		documentID := fmt.Sprintf("%s-%s", lng, ns)

		// Delete the specific document by its _id
		filter := bson.M{"_id": documentID}

		result, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			logrus.Errorf("Failed to delete document %s: %v", documentID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
			return
		}

		if result.DeletedCount == 0 {
			logrus.Infof("No document found with id %s", documentID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
	}
}
