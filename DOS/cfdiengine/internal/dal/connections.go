package dal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Sets up a MongoDB connection with a pool of connections
func SetUpConnMongoDB(mcli **mongo.Client, uri string) error {
	// Define connection pool options
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(20).                  // Maximum number of connections in the pool
		SetMinPoolSize(5).                   // Minimum number of connections in the pool
		SetMaxConnIdleTime(10 * time.Minute) // Connection idle time before closing

	// Connect to MongoDB with a timeout
	ctxConn, cancelConn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConn()

	cli, err := mongo.Connect(ctxConn, clientOptions)
	if err != nil {
		return err
	}

	// Set up a timeout for the ping
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()

	// Ping the MongoDB server to ensure the connection is established
	if err = cli.Ping(ctxPing, readpref.Primary()); err != nil {
		return err
	}

	// Set the client if the connection and ping succeed
	*mcli = cli

	return nil
}
