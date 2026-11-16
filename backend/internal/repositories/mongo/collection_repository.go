package mongorepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

type collectionDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	StoreID     string             `bson:"storeId"`
	Name        string             `bson:"name"`
	Slug        string             `bson:"slug"`
	Description string             `bson:"description"`
	ImageURL    string             `bson:"imageUrl"`
	ProductIDs  []string           `bson:"productIds"`
	Published   bool               `bson:"published"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

func newCollectionDoc(collection models.Collection) collectionDoc {
	ids := collection.ProductIDs
	if ids == nil {
		ids = []string{}
	}
	return collectionDoc{
		StoreID: collection.StoreID, Name: collection.Name, Slug: collection.Slug,
		Description: collection.Description, ImageURL: collection.ImageURL, ProductIDs: ids,
		Published: collection.Published, CreatedAt: utc(collection.CreatedAt), UpdatedAt: utc(collection.UpdatedAt),
	}
}

func (d collectionDoc) model() models.Collection {
	ids := d.ProductIDs
	if ids == nil {
		ids = []string{}
	}
	return models.Collection{
		ID: d.ID.Hex(), StoreID: d.StoreID, Name: d.Name, Slug: d.Slug,
		Description: d.Description, ImageURL: d.ImageURL, ProductIDs: ids, Published: d.Published,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type CollectionRepository struct {
	coll *mongo.Collection
}

func NewCollectionRepository(db *mongo.Database) *CollectionRepository {
	return &CollectionRepository{coll: db.Collection("collections")}
}

func (r *CollectionRepository) Create(ctx context.Context, collection models.Collection) (models.Collection, error) {
	result, err := r.coll.InsertOne(ctx, newCollectionDoc(collection))
	if err != nil {
		return models.Collection{}, mapErr(err)
	}
	collection.ID = result.InsertedID.(primitive.ObjectID).Hex()
	if collection.ProductIDs == nil {
		collection.ProductIDs = []string{}
	}
	return collection, nil
}

func (r *CollectionRepository) Update(ctx context.Context, collection models.Collection) (models.Collection, error) {
	oid, err := objectID(collection.ID)
	if err != nil {
		return models.Collection{}, err
	}
	ids := collection.ProductIDs
	if ids == nil {
		ids = []string{}
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid, "storeId": collection.StoreID}, bson.M{"$set": bson.M{
		"name": collection.Name, "slug": collection.Slug, "description": collection.Description,
		"imageUrl": collection.ImageURL, "productIds": ids, "published": collection.Published,
		"updatedAt": utc(collection.UpdatedAt),
	}})
	if err != nil {
		return models.Collection{}, mapErr(err)
	}
	if result.MatchedCount == 0 {
		return models.Collection{}, interfaces.ErrNotFound
	}
	collection.ProductIDs = ids
	return collection, nil
}

func (r *CollectionRepository) Delete(ctx context.Context, storeID, id string) error {
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

func (r *CollectionRepository) FindByID(ctx context.Context, storeID, id string) (models.Collection, error) {
	oid, err := objectID(id)
	if err != nil {
		return models.Collection{}, err
	}
	var doc collectionDoc
	if err := r.coll.FindOne(ctx, bson.M{"_id": oid, "storeId": storeID}).Decode(&doc); err != nil {
		return models.Collection{}, mapErr(err)
	}
	return doc.model(), nil
}

func (r *CollectionRepository) FindBySlug(ctx context.Context, storeID, slug string) (models.Collection, error) {
	var doc collectionDoc
	if err := r.coll.FindOne(ctx, bson.M{"storeId": storeID, "slug": slug}).Decode(&doc); err != nil {
		return models.Collection{}, mapErr(err)
	}
	return doc.model(), nil
}

func (r *CollectionRepository) List(ctx context.Context, storeID string, filter interfaces.CollectionFilter) ([]models.Collection, error) {
	query := bson.M{"storeId": storeID}
	if filter.Published != nil {
		query["published"] = *filter.Published
	}
	cursor, err := r.coll.Find(ctx, query, options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(limitOf(filter.Limit, 100, 100)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	items := []models.Collection{}
	for cursor.Next(ctx) {
		var doc collectionDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		items = append(items, doc.model())
	}
	return items, cursor.Err()
}

func (r *CollectionRepository) DetachProduct(ctx context.Context, storeID, productID string) error {
	_, err := r.coll.UpdateMany(ctx, bson.M{"storeId": storeID}, bson.M{"$pull": bson.M{"productIds": productID}})
	return err
}
