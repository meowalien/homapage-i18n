package mongodb

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://homepage-i18n:kingkingjin@mongodb-cluster-0.mongodb-cluster-svc.mongodb-cluster.svc.cluster.local:27017,mongodb-cluster-1.mongodb-cluster-svc.mongodb-cluster.svc.cluster.local:27017,mongodb-cluster-2.mongodb-cluster-svc.mongodb-cluster.svc.cluster.local:27017"

var client *mongo.Client

func ConnectDB() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		logrus.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Fatalf("Failed to ping MongoDB: %v", err)
	}

	logrus.Println("Connected to MongoDB!")
}

func DisconnectDB() {
	if err := client.Disconnect(context.TODO()); err != nil {
		logrus.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
	logrus.Println("Disconnected from MongoDB!")
}

func GetCollection(database, collectionName string) *mongo.Collection {
	return client.Database(database).Collection(collectionName)
}
