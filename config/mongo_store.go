package config

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (m *MongoRouteStore) SaveRoute(route *RouteConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if route.ID == "" {
		route.ID = uuid.NewString()
	}

	_, err := m.collection.InsertOne(ctx, route)
	return err
}

// DeleteRoute removes a route by its ID from MongoDB
func (m *MongoRouteStore) DeleteRoute(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Normalize
	id = strings.TrimSpace(id)

	// 1) Preferred: _id is a string UUID
	if res, err := m.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	} else if res.DeletedCount > 0 {
		return nil
	}

	// 2) Legacy: separate "id" field (older documents)
	if res, err := m.collection.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return err
	} else if res.DeletedCount > 0 {
		return nil
	}

	// 3) Very old: _id is an ObjectId
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		if res, err := m.collection.DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
			return err
		} else if res.DeletedCount > 0 {
			return nil
		}
	}

	return errors.New("route not found")
}
