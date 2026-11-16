package mongorepo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

type OrderRepository struct {
	coll *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{coll: db.Collection("orders")}
}

func (r *OrderRepository) Create(ctx context.Context, order models.Order) (models.Order, error) {
	result, err := r.coll.InsertOne(ctx, newOrderDoc(order))
	if err != nil {
		return models.Order{}, mapErr(err)
	}
	order.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return order, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, storeID, id string) (models.Order, error) {
	oid, err := objectID(id)
	if err != nil {
		return models.Order{}, err
	}
	var doc orderDoc
	if err := r.coll.FindOne(ctx, bson.M{"_id": oid, "storeId": storeID}).Decode(&doc); err != nil {
		return models.Order{}, mapErr(err)
	}
	return doc.model(), nil
}

func (r *OrderRepository) List(ctx context.Context, storeID string, filter interfaces.OrderFilter) ([]models.Order, error) {
	query := bson.M{"storeId": storeID}
	if filter.Status != "" {
		query["status"] = filter.Status
	} else if filter.ExcludeStatus != "" {
		query["status"] = bson.M{"$ne": filter.ExcludeStatus}
	}
	if !filter.Since.IsZero() {
		query["createdAt"] = bson.M{"$gte": filter.Since.UTC()}
	}
	cursor, err := r.coll.Find(ctx, query, options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(limitOf(filter.Limit, 50, 200)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	orders := []models.Order{}
	for cursor.Next(ctx) {
		var doc orderDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		orders = append(orders, doc.model())
	}
	return orders, cursor.Err()
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, storeID, id string, status models.OrderStatus, at time.Time) (models.Order, error) {
	oid, err := objectID(id)
	if err != nil {
		return models.Order{}, err
	}
	var doc orderDoc
	err = r.coll.FindOneAndUpdate(ctx,
		bson.M{"_id": oid, "storeId": storeID, "status": bson.M{"$ne": models.OrderCancelled}},
		bson.M{"$set": bson.M{"status": status, "updatedAt": utc(at)}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc)
	if err == nil {
		return doc.model(), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return models.Order{}, err
	}
	if _, findErr := r.FindByID(ctx, storeID, id); findErr != nil {
		return models.Order{}, findErr
	}
	return models.Order{}, interfaces.ErrConflict
}

func (r *OrderRepository) Summary(ctx context.Context, storeID string, since time.Time, exclude models.OrderStatus) (int64, int, error) {
	cursor, err := r.coll.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: salesMatch(storeID, since, exclude)}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":     nil,
			"revenue": bson.M{"$sum": "$totalCents"},
			"orders":  bson.M{"$sum": 1},
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
		Revenue int64 `bson:"revenue"`
		Orders  int   `bson:"orders"`
	}
	if err := cursor.Decode(&row); err != nil {
		return 0, 0, err
	}
	return row.Revenue, row.Orders, nil
}

func (r *OrderRepository) RevenueByDay(ctx context.Context, storeID string, since time.Time, exclude models.OrderStatus) ([]models.RevenuePoint, error) {
	cursor, err := r.coll.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: salesMatch(storeID, since, exclude)}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":     bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt", "timezone": "UTC"}},
			"revenue": bson.M{"$sum": "$totalCents"},
			"orders":  bson.M{"$sum": 1},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	points := []models.RevenuePoint{}
	for cursor.Next(ctx) {
		var row struct {
			Date    string `bson:"_id"`
			Revenue int64  `bson:"revenue"`
			Orders  int    `bson:"orders"`
		}
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		points = append(points, models.RevenuePoint{Date: row.Date, RevenueCents: row.Revenue, Orders: row.Orders})
	}
	return points, cursor.Err()
}

func (r *OrderRepository) BestSellers(ctx context.Context, storeID string, since time.Time, exclude models.OrderStatus, limit int) ([]models.BestSeller, error) {
	cursor, err := r.coll.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: salesMatch(storeID, since, exclude)}},
		bson.D{{Key: "$unwind", Value: "$items"}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":      "$items.productId",
			"name":     bson.M{"$last": "$items.name"},
			"quantity": bson.M{"$sum": "$items.quantity"},
			"revenue":  bson.M{"$sum": "$items.lineTotalCents"},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "quantity", Value: -1}, {Key: "revenue", Value: -1}}}},
		bson.D{{Key: "$limit", Value: limitOf(limit, 5, 20)}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	sellers := []models.BestSeller{}
	for cursor.Next(ctx) {
		var row struct {
			ID       string `bson:"_id"`
			Name     string `bson:"name"`
			Quantity int    `bson:"quantity"`
			Revenue  int64  `bson:"revenue"`
		}
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		sellers = append(sellers, models.BestSeller{ProductID: row.ID, Name: row.Name, Quantity: row.Quantity, RevenueCents: row.Revenue})
	}
	return sellers, cursor.Err()
}

func salesMatch(storeID string, since time.Time, exclude models.OrderStatus) bson.M {
	match := bson.M{
		"storeId":   storeID,
		"createdAt": bson.M{"$gte": since.UTC()},
	}
	if exclude != "" {
		match["status"] = bson.M{"$ne": exclude}
	}
	return match
}
