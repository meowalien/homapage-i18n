package routes

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"homapage-i18n/mongodb"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getI18nJSON(c *gin.Context) {
	lng := c.Param("lng")
	ns := c.Param("ns")

	// Debugging output to verify parameters
	logrus.Infof("lng: %s, ns: %s\n", lng, ns)

	collectionName := fmt.Sprintf("%s-%s", lng, ns)
	collection := mongodb.GetCollection("i18n", collectionName)

	var jsonData map[string]interface{}
	err := collection.FindOne(context.TODO(), bson.D{}).Decode(&jsonData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Infof("No data found in collection %s", collectionName)
			c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		} else {
			logrus.Errorf("Failed to find data in collection %s: %v", collectionName, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, jsonData)
}

func postI18nJSON(c *gin.Context) {
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
