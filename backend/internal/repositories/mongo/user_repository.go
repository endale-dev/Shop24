package mongorepo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"shop24/internal/domain/models"
)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection("users")}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	doc := userDoc{
		Name: user.Name, Email: user.Email, PasswordHash: user.PasswordHash,
		CreatedAt: utc(user.CreatedAt), UpdatedAt: utc(user.UpdatedAt),
	}
	result, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return models.User{}, mapErr(err)
	}
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	return r.find(ctx, bson.M{"email": email})
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (models.User, error) {
	oid, err := objectID(id)
	if err != nil {
		return models.User{}, err
	}
	return r.find(ctx, bson.M{"_id": oid})
}

func (r *UserRepository) find(ctx context.Context, filter bson.M) (models.User, error) {
	var doc userDoc
	if err := r.coll.FindOne(ctx, filter).Decode(&doc); err != nil {
		return models.User{}, mapErr(err)
	}
	return doc.model(), nil
}
