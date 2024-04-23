package config

import (
	"context"
	"os"

	"github.com/apexwoot/mr_scraper/internal/app"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewConnection(connUri string) (*app.DB, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(ComposeDBConnectionString()).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))

	return &app.DB{Client: client, FocusedCollection: collection}, nil
}
