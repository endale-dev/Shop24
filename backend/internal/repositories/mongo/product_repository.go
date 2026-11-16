package mongorepo

import (
	"context"
	"errors"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

type ProductRepository struct {
	coll *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{coll: db.Collection("products")}
}

func (r *ProductRepository) Create(ctx context.Context, product models.Product) (models.Product, error) {
	result, err := r.coll.InsertOne(ctx, newProductDoc(product))
	if err != nil {
		return models.Product{}, mapErr(err)
	}
	product.ID = result.InsertedID.(primitive.ObjectID).Hex()
	if product.Images == nil {
		product.Images = []string{}
	}
	return product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product models.Product) (models.Product, error) {
	oid, err := objectID(product.ID)
	if err != nil {
		return models.Product{}, err
	}
	images := product.Images
	if images == nil {
		images = []string{}
	}
	names := product.OptionNames
	if names == nil {
		names = []string{}
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid, "storeId": product.StoreID}, bson.M{"$set": bson.M{
		"name": product.Name, "description": product.Description, "priceCents": product.PriceCents,
		"images": images, "category": product.Category, "sku": product.SKU, "stock": product.Stock,
		"optionNames": names, "variants": variantDocs(product.Variants),
		"published": product.Published, "updatedAt": utc(product.UpdatedAt),
	}})
	if err != nil {
		return models.Product{}, mapErr(err)
	}
	if result.MatchedCount == 0 {
		return models.Product{}, interfaces.ErrNotFound
	}
	product.Images = images
	return product, nil
}

func (r *ProductRepository) Delete(ctx context.Context, storeID, id string) error {
	oid, err := objectID(id)
	if err != nil {
		return err
	}
	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": oid, "storeId": storeID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return interfaces.ErrNotFound
	}
	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, storeID, id string) (models.Product, error) {
	oid, err := objectID(id)
	if err != nil {
		return models.Product{}, err
	}
	var doc productDoc
	if err := r.coll.FindOne(ctx, bson.M{"_id": oid, "storeId": storeID}).Decode(&doc); err != nil {
		return models.Product{}, mapErr(err)
	}
	return doc.model(), nil
}

func (r *ProductRepository) List(ctx context.Context, storeID string, filter interfaces.ProductFilter) ([]models.Product, error) {
	query := bson.M{"storeId": storeID}
	if filter.Published != nil {
		query["published"] = *filter.Published
	}
	clauses := make([]bson.M, 0, 2)
	if filter.LowStock {
		clauses = append(clauses, bson.M{"$or": []bson.M{
			{"variants": bson.M{"$elemMatch": bson.M{"stock": bson.M{"$lte": filter.LowStockAt}}}},
			{"variants.0": bson.M{"$exists": false}, "stock": bson.M{"$lte": filter.LowStockAt}},
		}})
	}
	if filter.Query != "" {
		pattern := regexp.QuoteMeta(filter.Query)
		clauses = append(clauses, bson.M{"$or": []bson.M{
			{"name": bson.M{"$regex": pattern, "$options": "i"}},
			{"sku": bson.M{"$regex": pattern, "$options": "i"}},
			{"variants.sku": bson.M{"$regex": pattern, "$options": "i"}},
		}})
	}
	if len(clauses) == 1 {
		for key, value := range clauses[0] {
			query[key] = value
		}
	} else if len(clauses) > 1 {
		query["$and"] = clauses
	}
	cursor, err := r.coll.Find(ctx, query, options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(limitOf(filter.Limit, 100, 200)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	products := []models.Product{}
	for cursor.Next(ctx) {
		var doc productDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		products = append(products, doc.model())
	}
	return products, cursor.Err()
}

func (r *ProductRepository) AddStock(ctx context.Context, storeID, id, variantID string, delta int) error {
	oid, err := objectID(id)
	if err != nil {
		return err
	}
	if variantID == "" {
		variantID = id
	}
	now := time.Now().UTC()
	match := bson.M{"id": variantID}
	if delta < 0 {
		match["stock"] = bson.M{"$gte": -delta}
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{
		"_id": oid, "storeId": storeID, "variants": bson.M{"$elemMatch": match},
	}, bson.M{
		"$inc": bson.M{"variants.$.stock": delta, "stock": delta},
		"$set": bson.M{"updatedAt": now},
	})
	if err != nil {
		return err
	}
	if result.MatchedCount == 1 {
		return nil
	}
	if variantID == id {
		legacy := bson.M{"_id": oid, "storeId": storeID, "variants.0": bson.M{"$exists": false}}
		if delta < 0 {
			legacy["stock"] = bson.M{"$gte": -delta}
		}
		result, err = r.coll.UpdateOne(ctx, legacy, bson.M{
			"$inc": bson.M{"stock": delta},
			"$set": bson.M{"updatedAt": now},
		})
		if err != nil {
			return err
		}
		if result.MatchedCount == 1 {
			return nil
		}
	}
	var doc productDoc
	err = r.coll.FindOne(ctx, bson.M{"_id": oid, "storeId": storeID}).Decode(&doc)
	if err != nil {
		return mapErr(err)
	}
	for _, variant := range doc.Variants {
		if variant.ID == variantID {
			return interfaces.ErrInsufficientStock
		}
	}
	if len(doc.Variants) == 0 {
		return interfaces.ErrInsufficientStock
	}
	return interfaces.ErrNotFound
}

func (r *ProductRepository) SetStock(ctx context.Context, storeID, id, variantID string, stock int, at time.Time) error {
	oid, err := objectID(id)
	if err != nil {
		return err
	}
	if variantID == "" {
		variantID = id
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$set", Value: bson.M{
			"variants": bson.M{"$map": bson.M{
				"input": bson.M{"$ifNull": bson.A{"$variants", bson.A{}}},
				"as":    "v",
				"in": bson.M{"$cond": bson.A{
					bson.M{"$eq": bson.A{"$$v.id", variantID}},
					bson.M{"$mergeObjects": bson.A{"$$v", bson.M{"stock": stock}}},
					"$$v",
				}},
			}},
			"updatedAt": utc(at),
		}}},
		bson.D{{Key: "$set", Value: bson.M{"stock": bson.M{"$sum": "$variants.stock"}}}},
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid, "storeId": storeID, "variants.id": variantID}, pipeline)
	if err != nil {
		return err
	}
	if result.MatchedCount == 1 {
		return nil
	}
	if variantID == id {
		result, err = r.coll.UpdateOne(ctx, bson.M{
			"_id": oid, "storeId": storeID, "variants.0": bson.M{"$exists": false},
		}, bson.M{"$set": bson.M{"stock": stock, "updatedAt": utc(at)}})
		if err != nil {
			return err
		}
		if result.MatchedCount == 1 {
			return nil
		}
	}
	err = r.coll.FindOne(ctx, bson.M{"_id": oid, "storeId": storeID}).Err()
	return mapErr(err)
}

func (r *ProductRepository) SKUInUse(ctx context.Context, storeID, sku, exceptProductID string) (bool, error) {
	query := bson.M{
		"storeId": storeID,
		"$or": []bson.M{
			{"sku": sku},
			{"variants.sku": sku},
		},
	}
	if exceptProductID != "" {
		oid, err := objectID(exceptProductID)
		if err != nil {
			return false, err
		}
		query["_id"] = bson.M{"$ne": oid}
	}
	err := r.coll.FindOne(ctx, query).Err()
	if err == nil {
		return true, nil
	}
	if errors.Is(mapErr(err), interfaces.ErrNotFound) {
		return false, nil
	}
	return false, err
}

func (r *ProductRepository) Totals(ctx context.Context, storeID string) (int, int, error) {
	cursor, err := r.coll.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"storeId": storeID}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"count": bson.M{"$sum": 1},
			"stock": bson.M{"$sum": "$stock"},
		}}},
	})
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)
	if !cursor.Next(ctx) {
		return 0, 0, cursor.Err()
	}
	var row struct {
		Count int `bson:"count"`
		Stock int `bson:"stock"`
	}
	if err := cursor.Decode(&row); err != nil {
		return 0, 0, err
	}
	return row.Count, row.Stock, nil
}
