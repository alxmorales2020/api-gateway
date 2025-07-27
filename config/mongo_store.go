package config

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRouteStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoRouteStore creates a RouteStore backed by MongoDB
func NewMongoRouteStore(cfg *MongoDBConfig) (RouteStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)

	// Set credentials if provided
	if cfg.Username != "" && cfg.Password != "" {
		clientOpts.SetAuth(options.Credential{
			Username: cfg.Username,
			Password: cfg.Password,
		})
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// Ping to ensure connection is valid
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	dbName := cfg.Database
	if dbName == "" {
		dbName = "apigateway"
	}
	collName := cfg.Collection
	if collName == "" {
		collName = "routes"
	}

	collection := client.Database(dbName).Collection(collName)
	return &MongoRouteStore{
		client:     client,
		collection: collection,
	}, nil
}

// LoadRoutes fetches all route documents from MongoDB
func (m *MongoRouteStore) LoadRoutes() ([]RouteConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var routes []RouteConfig
	for cursor.Next(ctx) {
		var route RouteConfig
		if err := cursor.Decode(&route); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	return routes, nil
}

// SaveRoute inserts a new route into MongoDB
func (m *MongoRouteStore) SaveRoute(route RouteConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.collection.InsertOne(ctx, route)
	return err
}

// DeleteRoute removes a route by path
func (m *MongoRouteStore) DeleteRoute(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := m.collection.DeleteOne(ctx, bson.M{"path": path})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("route not found")
	}
	return nil
}
