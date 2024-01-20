package db

import (
	"context"

	"github.com/HrvojeLesar/recommender/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoInstance struct {
	ctx      context.Context
	config   config.MongoConfig
	client   *mongo.Client
	database *mongo.Database
}

func Setup(ctx context.Context, config config.Config) (*MongoInstance, error) {
	clientOptions := options.Client().ApplyURI(config.Mongo.Uri())

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	database := client.Database(config.Mongo.Database())

	inst := MongoInstance{
		ctx:      ctx,
		client:   client,
		database: database,
		config:   config.Mongo,
	}

	return &inst, nil
}
