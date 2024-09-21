package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"homapage-i18n/mongodb"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Read the MongoDB URI from the environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	// Connect to MongoDB using the URI
	mongodb.ConnectDB(mongoURI)

	// Iterate over all folders under i18n directory
	rootDir := "./i18n"
	folders, err := os.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range folders {
		if folder.IsDir() {
			collectionName := folder.Name()
			dirPath := filepath.Join(rootDir, collectionName)

			// Insert JSON files from this folder into the corresponding collection
			err := insertJSONFiles(collectionName, dirPath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("Successfully inserted all files into corresponding collections.")
}

// insertJSONFiles reads JSON files from the directory and inserts them into the specified collection
func insertJSONFiles(collectionName, dirPath string) error {
	collection := mongodb.GetCollection("i18n", collectionName)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".json" {
			fileNameWithoutExt := file.Name()[:len(file.Name())-len(ext)]

			// delete the document with the same _id
			_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": fileNameWithoutExt})
			if err != nil {
				return err
			}

			filePath := filepath.Join(dirPath, file.Name())

			// Open and read the JSON file
			jsonFile, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer jsonFile.Close()

			byteValue, err := io.ReadAll(jsonFile)
			if err != nil {
				return err
			}

			// Convert the byte data to a map for MongoDB insertion
			var documentContent bson.M
			if err := json.Unmarshal(byteValue, &documentContent); err != nil {
				return err
			}

			var document = bson.M{}

			document["_id"] = fileNameWithoutExt
			document["content"] = documentContent

			// Insert the documentContent into the collection
			_, err = collection.InsertOne(context.TODO(), document)
			if err != nil {
				return err
			}

			fmt.Printf("Inserted %s into %s collection\n", file.Name(), collectionName)
		}
	}

	return nil
}
