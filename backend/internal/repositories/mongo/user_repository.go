package mongorepo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection("users")}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	doc := userDoc{
		Name: user.Name, Email: user.Email, PasswordHash: user.PasswordHash, GoogleSub: user.GoogleSub,
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

func (r *UserRepository) FindByGoogleSub(ctx context.Context, sub string) (models.User, error) {
	if sub == "" {
		return models.User{}, interfaces.ErrNotFound
	}
	return r.find(ctx, bson.M{"googleSub": sub})
}

func (r *UserRepository) FindByResetHash(ctx context.Context, hash string) (models.User, error) {
	if hash == "" {
		return models.User{}, interfaces.ErrNotFound
	}
	return r.find(ctx, bson.M{"resetHash": hash})
}

func (r *UserRepository) Update(ctx context.Context, user models.User) error {
	oid, err := objectID(user.ID)
	if err != nil {
		return err
	}
	set := bson.M{
		"name": user.Name, "email": user.Email, "passwordHash": user.PasswordHash, "updatedAt": utc(user.UpdatedAt),
	}
	unset := bson.M{}
	if user.GoogleSub != "" {
		set["googleSub"] = user.GoogleSub
	} else {
		unset["googleSub"] = ""
	}
	if user.ResetHash != "" {
		set["resetHash"] = user.ResetHash
		set["resetExpires"] = utc(user.ResetExpires)
	} else {
		unset["resetHash"] = ""
		unset["resetExpires"] = ""
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set, "$unset": unset})
	if err != nil {
		return mapErr(err)
	}
	if result.MatchedCount == 0 {
		return interfaces.ErrNotFound
	}
	return nil
}

func (r *UserRepository) find(ctx context.Context, filter bson.M) (models.User, error) {
	var doc userDoc
	if err := r.coll.FindOne(ctx, filter).Decode(&doc); err != nil {
		return models.User{}, mapErr(err)
	}
	return doc.model(), nil
}
