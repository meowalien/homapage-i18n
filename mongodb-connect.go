package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the placeholder with your Atlas connection string
const uri = "mongodb://homepage-i18n:kingkingjin@mongodb-cluster-0.mongodb-cluster-svc.mongodb-cluster.svc.cluster.local:27017,mongodb-cluster-1.mongodb-cluster-svc.mongodb-cluster.svc.cluster.local:27017,mongodb-cluster-2.mongodb-cluster-svc.mongodb-cluster.svc.cluster.local:27017"

func main() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Define the directories containing the JSON files
	directories := []string{"./i18n/en", "./i18n/zh-Hant"}

	for _, dir := range directories {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".json") {
				filePath := filepath.Join(dir, file.Name())
				jsonFile, err := os.Open(filePath)
				if err != nil {
					panic(err)
				}
				defer jsonFile.Close()

				byteValue, err := ioutil.ReadAll(jsonFile)
				if err != nil {
					panic(err)
				}

				var jsonData map[string]interface{}
				if err := json.Unmarshal(byteValue, &jsonData); err != nil {
					panic(err)
				}

				// Insert the JSON data into the i18n database and a specific collection
				collectionName := fmt.Sprintf("%s-%s", filepath.Base(dir), strings.TrimSuffix(file.Name(), ".json"))
				collection := client.Database("i18n").Collection(collectionName)
				insertResult, err := collection.InsertOne(context.TODO(), jsonData)
				if err != nil {
					panic(err)
				}

				fmt.Printf("Inserted document into collection %s with ID: %v\n", collectionName, insertResult.InsertedID)
			}
		}
	}
}
