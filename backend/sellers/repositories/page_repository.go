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

type PageRepository struct {
	col *mongo.Collection
}

func NewPageRepository() *PageRepository {
	col := config.GetCollection(config.GetDB().Name(), "pages")
	return &PageRepository{col: col}
}

func (r *PageRepository) Create(ctx context.Context, p *models.Page) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.col.InsertOne(ctx, p)
	return err
}

func (r *PageRepository) Update(ctx context.Context, id primitive.ObjectID, p *models.Page) error {
	p.UpdatedAt = time.Now()
	_, err := r.col.ReplaceOne(ctx, bson.M{"_id": id}, p)
	return err
}

func (r *PageRepository) FindByStoreAndPath(ctx context.Context, storeID primitive.ObjectID, path string) (*models.Page, error) {
	var out models.Page
	err := r.col.FindOne(ctx, bson.M{"store_id": storeID, "path": path}).Decode(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PageRepository) PublishSnapshot(ctx context.Context, id primitive.ObjectID, snapshot map[string]interface{}) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"snapshot": snapshot, "status": "published", "updated_at": time.Now()}})
	return err
}

