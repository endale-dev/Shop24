package mongorepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}
	return client, nil
}

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	unique := func(name string, keys bson.D) mongo.IndexModel {
		return mongo.IndexModel{Keys: keys, Options: options.Index().SetUnique(true).SetName(name)}
	}
	index := func(name string, keys bson.D) mongo.IndexModel {
		return mongo.IndexModel{Keys: keys, Options: options.Index().SetName(name)}
	}
	if _, err := db.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
		unique("email", bson.D{{Key: "email", Value: 1}}),
		{Keys: bson.D{{Key: "googleSub", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true).SetName("googleSub")},
		index("resetHash", bson.D{{Key: "resetHash", Value: 1}}),
	}); err != nil {
		return err
	}
	if _, err := db.Collection("stores").Indexes().CreateMany(ctx, []mongo.IndexModel{
		unique("slug", bson.D{{Key: "slug", Value: 1}}),
		unique("ownerId", bson.D{{Key: "ownerId", Value: 1}}),
	}); err != nil {
		return err
	}
	if _, err := db.Collection("products").Indexes().CreateMany(ctx, []mongo.IndexModel{
		unique("store_sku", bson.D{{Key: "storeId", Value: 1}, {Key: "sku", Value: 1}}),
		index("store_created", bson.D{{Key: "storeId", Value: 1}, {Key: "createdAt", Value: -1}}),
		index("store_published", bson.D{{Key: "storeId", Value: 1}, {Key: "published", Value: 1}}),
	}); err != nil {
		return err
	}
	if _, err := db.Collection("orders").Indexes().CreateMany(ctx, []mongo.IndexModel{
		index("store_orders", bson.D{{Key: "storeId", Value: 1}, {Key: "createdAt", Value: -1}}),
		index("store_status", bson.D{{Key: "storeId", Value: 1}, {Key: "status", Value: 1}}),
	}); err != nil {
		return err
	}
	_, err := db.Collection("collections").Indexes().CreateMany(ctx, []mongo.IndexModel{
		unique("collection_slug", bson.D{{Key: "storeId", Value: 1}, {Key: "slug", Value: 1}}),
		index("collection_created", bson.D{{Key: "storeId", Value: 1}, {Key: "createdAt", Value: -1}}),
	})
	return err
}
