package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ktravels-publicsite-api/internal/features/currency/domain"
)

const collectionName = "currencies"

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	r := &MongoRepository{collection: db.Collection(collectionName)}
	r.ensureIndexes(context.Background())
	return r
}

func (r *MongoRepository) ensureIndexes(ctx context.Context) {
	r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}

func (r *MongoRepository) Upsert(ctx context.Context, currency *domain.Currency) error {
	doc := toDocument(currency)
	filter := bson.M{"_id": doc.Code}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert currency %s: %w", currency.Code, err)
	}
	return nil
}

func (r *MongoRepository) FindAll(ctx context.Context) ([]domain.Currency, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find currencies: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []currencyDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode currencies: %w", err)
	}

	currencies := make([]domain.Currency, len(docs))
	for i, doc := range docs {
		currencies[i] = *fromDocument(&doc)
	}
	return currencies, nil
}
