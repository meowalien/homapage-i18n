package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"homapage-i18n/mongodb"
	"net/http"
)

func PostI18n() gin.HandlerFunc {
	return func(c *gin.Context) {

		lng := c.Param("lng")
		ns := c.Param("ns")

		collection := mongodb.GetCollection("homepage", "i18n")

		documentID := fmt.Sprintf("%s-%s", lng, ns)

		var jsonData map[string]interface{}
		if err := c.BindJSON(&jsonData); err != nil {
			logrus.Errorf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Remove _id field if it exists to avoid the immutable field error
		if _, ok := jsonData["_id"]; ok {
			delete(jsonData, "_id")
		}

		// Wrap jsonData inside "content" field
		updateData := bson.M{"content": jsonData}

		// Prepare the filter and update for the upsert operation
		filter := bson.M{"_id": documentID}
		update := bson.M{"$set": updateData}
		opts := options.Update().SetUpsert(true)

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			logrus.Errorf("Failed to update data in document %s: %v", documentID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data"})
			return
		}

		if updateResult.MatchedCount > 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Document updated", "matchedCount": updateResult.MatchedCount})
		} else {
			c.JSON(http.StatusCreated, gin.H{"message": "Document inserted", "upsertedId": updateResult.UpsertedID})
		}
	}
}
