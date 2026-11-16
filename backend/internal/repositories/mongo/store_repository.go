package mongorepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

type StoreRepository struct {
	coll *mongo.Collection
}

func NewStoreRepository(db *mongo.Database) *StoreRepository {
	return &StoreRepository{coll: db.Collection("stores")}
}

func (r *StoreRepository) Create(ctx context.Context, store models.Store) (models.Store, error) {
	doc := storeDoc{
		OwnerID: store.OwnerID, Name: store.Name, Slug: store.Slug, Description: store.Description,
		LogoURL: store.LogoURL, Currency: store.Currency, Published: store.Published,
		Storefront: storefrontToDoc(store.Storefront), Draft: storefrontToDoc(store.Draft),
		DesignNote: nil, CreatedAt: utc(store.CreatedAt), UpdatedAt: utc(store.UpdatedAt),
	}
	if store.DesignNote != nil {
		note := noteToDoc(*store.DesignNote)
		doc.DesignNote = &note
	}
	result, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return models.Store{}, mapErr(err)
	}
	store.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return store, nil
}

func (r *StoreRepository) FindByOwner(ctx context.Context, ownerID string) (models.Store, error) {
	return r.find(ctx, bson.M{"ownerId": ownerID})
}

func (r *StoreRepository) FindBySlug(ctx context.Context, slug string) (models.Store, error) {
	return r.find(ctx, bson.M{"slug": slug})
}

func (r *StoreRepository) UpdateProfile(ctx context.Context, store models.Store) error {
	oid, err := objectID(store.ID)
	if err != nil {
		return err
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid, "ownerId": store.OwnerID}, bson.M{"$set": bson.M{
		"name": store.Name, "slug": store.Slug, "description": store.Description,
		"logoUrl": store.LogoURL, "currency": store.Currency, "updatedAt": utc(store.UpdatedAt),
	}})
	if err != nil {
		return mapErr(err)
	}
	if result.MatchedCount == 0 {
		return interfaces.ErrNotFound
	}
	return nil
}

func (r *StoreRepository) SetPublished(ctx context.Context, ownerID, storeID string, published bool, at time.Time) error {
	return r.touch(ctx, ownerID, storeID, bson.M{"published": published, "updatedAt": utc(at)})
}

func (r *StoreRepository) UpdateDraft(ctx context.Context, ownerID, storeID string, draft models.StorefrontConfig, note *models.DesignNote, at time.Time) error {
	set := bson.M{"draft": storefrontToDoc(draft), "updatedAt": utc(at)}
	if note != nil {
		set["designNote"] = noteToDoc(*note)
	}
	return r.touch(ctx, ownerID, storeID, set)
}

func (r *StoreRepository) PublishDraft(ctx context.Context, ownerID, storeID string, at time.Time) error {
	oid, err := objectID(storeID)
	if err != nil {
		return err
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid, "ownerId": ownerID}, mongo.Pipeline{
		bson.D{{Key: "$set", Value: bson.M{"storefront": "$draft", "updatedAt": utc(at)}}},
	})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return interfaces.ErrNotFound
	}
	return nil
}

func (r *StoreRepository) touch(ctx context.Context, ownerID, storeID string, set bson.M) error {
	oid, err := objectID(storeID)
	if err != nil {
		return err
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid, "ownerId": ownerID}, bson.M{"$set": set})
	if err != nil {
		return mapErr(err)
	}
	if result.MatchedCount == 0 {
		return interfaces.ErrNotFound
	}
	return nil
}

func (r *StoreRepository) find(ctx context.Context, filter bson.M) (models.Store, error) {
	var doc storeDoc
	if err := r.coll.FindOne(ctx, filter).Decode(&doc); err != nil {
		return models.Store{}, mapErr(err)
	}
	return doc.model(), nil
}
