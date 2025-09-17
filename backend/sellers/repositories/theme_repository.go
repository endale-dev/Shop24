package repositories

import (
	"context"
	"time"

	"github.com/Endale2/DRPS/config"
	"github.com/Endale2/DRPS/sellers/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ThemeRepository struct {
	col *mongo.Collection
}

func NewThemeRepository() *ThemeRepository {
	col := config.GetCollection(config.GetDB().Name(), "themes")
	return &ThemeRepository{col: col}
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
	filter := bson.M{"store_id": storeID}
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

func (r *ThemeRepository) SetActive(ctx context.Context, storeID primitive.ObjectID, themeID primitive.ObjectID) error {
	// Unset previous
	_, _ = r.col.UpdateMany(ctx, bson.M{"store_id": storeID, "is_active": true}, bson.M{"$set": bson.M{"is_active": false}})
	// Set new
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": themeID, "store_id": storeID}, bson.M{"$set": bson.M{"is_active": true, "updated_at": time.Now()}})
	return err
}
