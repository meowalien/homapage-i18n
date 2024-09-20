package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/oliveagle/jsonpath"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"homapage-i18n/mongodb"
	"net/http"
)

func GetI18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		lng := c.Param("lng")
		ns := c.Param("ns")

		collection := mongodb.GetCollection("i18n", lng)

		var jsonData bson.M
		projection := bson.M{"_id": 0, "content": 1}

		err := collection.FindOne(context.TODO(), bson.M{"_id": ns}, options.FindOne().SetProjection(projection)).Decode(&jsonData)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				logrus.Infof("No data found in collection %s", ns)
				c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
			} else {
				logrus.Errorf("Failed to find document %s: %v", ns, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}

		content, err := jsonpath.JsonPathLookup(jsonData, "$.content")
		if err != nil {
			logrus.Errorf("Failed to find content in document %s: %v", ns, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		jsonData, ok := content.(bson.M)
		if !ok {
			logrus.Errorf("Failed to convert content to map[string]interface{}")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, jsonData)
	}
}
