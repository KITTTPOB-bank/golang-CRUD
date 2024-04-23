package dbconfig

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/KITTTPOB-bank/Carnan_backend/pkg/Key"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBconnect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var mongopath = Key.GetMongoDB()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongopath))

	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client

}
