package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oliveagle/jsonpath"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"homapage-i18n/mongodb"
	"net/http"
)

func ListI18nArticles() gin.HandlerFunc {
	return func(c *gin.Context) {
		lng := c.Param("lng")
		nsPrefix := c.Param("ns_prefix")
		collection := mongodb.GetCollection("homepage", "i18n")

		documentID := fmt.Sprintf("^%s-%s\\.", lng, nsPrefix)
		filter := bson.M{"_id": bson.M{"$regex": documentID}}
		projection := bson.M{"_id": 0, "content.title": 1}

		cursor, err := collection.Find(context.TODO(), filter, options.Find().SetProjection(projection))
		if err != nil {
			logrus.Errorf("Failed to find documents: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer cursor.Close(context.TODO())

		var results []interface{}
		for cursor.Next(context.Background()) {
			var result bson.M
			if err = cursor.Decode(&result); err != nil {
				logrus.Errorf("Failed to decode document: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}

			content, jsonPathLookupError := jsonpath.JsonPathLookup(result, "$.content.title")
			if jsonPathLookupError != nil {
				logrus.Errorf("Failed to find content.title in document %s: %v", result, jsonPathLookupError)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}

			results = append(results, content)
		}

		//return 404 if no data found
		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
