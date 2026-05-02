package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
)

const collectionName = "provider_services"

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *MongoRepository) Upsert(ctx context.Context, service *domain.ProviderService) error {
	filter := map[string]interface{}{"_id": service.ID}
	opts := options.Replace().SetUpsert(true)

	_, err := r.collection.ReplaceOne(ctx, filter, service, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) UpdateStatus(ctx context.Context, id string, status int) error {
	filter := map[string]interface{}{"_id": id}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"status": status,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
