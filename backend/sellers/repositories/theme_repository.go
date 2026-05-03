package repositories

import (
	"context"
	"time"

	"github.com/Endale2/DRPS/config"
	"github.com/Endale2/DRPS/sellers/models"
	sharedModels "github.com/Endale2/DRPS/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ThemeRepository struct {
	col      *mongo.Collection
	shopsCol *mongo.Collection
}

func NewThemeRepository() *ThemeRepository {
	col := config.GetCollection(config.GetDB().Name(), "themes")
	shopsCol := config.GetCollection(config.GetDB().Name(), "shops")
	return &ThemeRepository{
		col:      col,
		shopsCol: shopsCol,
	}
}

func (r *ThemeRepository) Create(ctx context.Context, t *models.Theme) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	res, err := r.col.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	t.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ThemeRepository) FindByStore(ctx context.Context, storeID primitive.ObjectID) ([]models.Theme, error) {
	filter := bson.M{"shop_id": storeID}
	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var out []models.Theme
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *ThemeRepository) FindActiveByShopID(ctx context.Context, shopID primitive.ObjectID) (*sharedModels.Theme, error) {
	filter := bson.M{"shop_id": shopID, "isActive": true}
	var theme sharedModels.Theme
	err := r.col.FindOne(ctx, filter).Decode(&theme)
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func (r *ThemeRepository) FindActiveByShopSlug(ctx context.Context, shopSlug string) (*sharedModels.Theme, error) {
	// First, find the shop by slug
	var shop sharedModels.Shop
	err := r.shopsCol.FindOne(ctx, bson.M{"slug": shopSlug}).Decode(&shop)
	if err != nil {
		return nil, err
	}

	// Then find the active theme for that shop
	return r.FindActiveByShopID(ctx, shop.ID)
}

func (r *ThemeRepository) SetActive(ctx context.Context, storeID primitive.ObjectID, themeID primitive.ObjectID) error {
	// Unset previous
	_, _ = r.col.UpdateMany(ctx, bson.M{"shop_id": storeID, "isActive": true}, bson.M{"$set": bson.M{"isActive": false}})
	// Set new
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": themeID, "shop_id": storeID}, bson.M{"$set": bson.M{"isActive": true, "updatedAt": time.Now()}})
	return err
}

func (r *ThemeRepository) Update(ctx context.Context, themeID primitive.ObjectID, update interface{}) error {
	filter := bson.M{"_id": themeID}
	updateDoc := bson.M{
		"$set": update,
	}
	_, err := r.col.UpdateOne(ctx, filter, updateDoc)
	return err
}
