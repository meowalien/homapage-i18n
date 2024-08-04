package routes

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"homapage-i18n/mongodb"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getI18nJSON(c *gin.Context) {
	lng := c.Param("lng")
	ns := c.Param("ns")

	// Debugging output to verify parameters
	logrus.Debugf("lng: %s, ns: %s\n", lng, ns)

	collectionName := fmt.Sprintf("%s-%s", lng, ns)
	collection := mongodb.GetCollection("i18n", collectionName)

	var jsonData map[string]interface{}
	err := collection.FindOne(context.TODO(), bson.D{}).Decode(&jsonData)
	if err != nil {
		logrus.Errorf("Failed to find data in collection %s: %v", collectionName, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	c.JSON(http.StatusOK, jsonData)
}
