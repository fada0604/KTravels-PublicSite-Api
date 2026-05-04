package infrastructure

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
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
		log.Printf("Failed to upsert %s: %v", service.ID, err)
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
		return fmt.Errorf("update error: %w", err)
	}

	return nil
}

func (r *MongoRepository) UpdateProviderLogo(ctx context.Context, providerID string, logo *domain.Image) error {
	filter := bson.M{"provider_id": providerID}

	var update bson.M
	if logo == nil {
		update = bson.M{
			"$unset": bson.M{
				"provider_logo": 1,
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				"provider_logo": logo,
			},
		}
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update provider logo for provider_id %s: %w", providerID, err)
	}

	return nil
}

func (r *MongoRepository) UpdateUnitImages(ctx context.Context, unitID string, images []domain.Image) error {
	filter := bson.M{"units.id": unitID}

	var update bson.M
	if images == nil {
		update = bson.M{
			"$unset": bson.M{
				"units.$.images": 1,
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				"units.$.images": images,
			},
		}
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update unit images for unit_id %s: %w", unitID, err)
	}

	return nil
}
