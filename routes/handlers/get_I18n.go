package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"homapage-i18n/mongodb"
	"net/http"
)

func GetI18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		lng := c.Param("lng")
		ns := c.Param("ns")

		// Debugging output to verify parameters
		//logrus.Infof("lng: %s, ns: %s\n", lng, ns)

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
}
