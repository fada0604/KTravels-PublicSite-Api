package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ktravels-publicsite-api/internal/features/exchange_rate/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

const collectionName = "exchange_rates"

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
		Keys:    bson.D{{Key: "event_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "currency_from", Value: 1},
			{Key: "currency_to", Value: 1},
			{Key: "valid_from", Value: -1},
		},
	})
}

func (r *MongoRepository) Insert(ctx context.Context, rate *domain.ExchangeRate) error {
	doc := toDocument(rate)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: event_id %s already processed", sharederrors.ErrConflict, rate.EventID)
		}
		return fmt.Errorf("failed to insert exchange rate: %w", err)
	}
	return nil
}

func (r *MongoRepository) ExistsByEventID(ctx context.Context, eventID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"event_id": eventID})
	if err != nil {
		return false, fmt.Errorf("failed to check event_id %s: %w", eventID, err)
	}
	return count > 0, nil
}

func (r *MongoRepository) FindLatestByPair(ctx context.Context) (map[string]float64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$sort", Value: bson.D{{Key: "valid_from", Value: -1}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "currency_from", Value: "$currency_from"},
				{Key: "currency_to", Value: "$currency_to"},
			}},
			{Key: "rate", Value: bson.D{{Key: "$first", Value: "$rate"}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate latest rates: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[string]float64)
	for cursor.Next(ctx) {
		var item struct {
			ID struct {
				CurrencyFrom string `bson:"currency_from"`
				CurrencyTo   string `bson:"currency_to"`
			} `bson:"_id"`
			Rate float64 `bson:"rate"`
		}
		if err := cursor.Decode(&item); err != nil {
			continue
		}
		result[item.ID.CurrencyFrom+":"+item.ID.CurrencyTo] = item.Rate
	}

	return result, cursor.Err()
}

func (r *MongoRepository) FindPaginated(ctx context.Context, page, pageSize int) ([]domain.ExchangeRate, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count exchange rates: %w", err)
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "valid_from", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find exchange rates: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []exchangeRateDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode exchange rates: %w", err)
	}

	rates := make([]domain.ExchangeRate, len(docs))
	for i, doc := range docs {
		rates[i] = *fromDocument(&doc)
	}

	return rates, total, nil
}
