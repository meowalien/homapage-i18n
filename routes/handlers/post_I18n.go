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

func PostI18n(c *gin.Context) {
	user_id, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}
	fmt.Println("user_id", user_id)

	lng := c.Param("lng")
	ns := c.Param("ns")

	collectionName := fmt.Sprintf("%s-%s", lng, ns)
	collection := mongodb.GetCollection("i18n", collectionName)

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

	filter := bson.D{} // Match all documents
	opts := options.Replace().SetUpsert(true)
	replaceResult, err := collection.ReplaceOne(context.TODO(), filter, jsonData, opts)
	if err != nil {
		logrus.Errorf("Failed to replace data in collection %s: %v", collectionName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to replace data"})
		return
	}

	if replaceResult.MatchedCount > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Document replaced", "matchedCount": replaceResult.MatchedCount})
	} else {
		c.JSON(http.StatusCreated, gin.H{"message": "Document inserted", "upsertedId": replaceResult.UpsertedID})
	}
}
