package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
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
	doc := toDocument(service)
	filter := bson.M{"_id": doc.ID}
	opts := options.Replace().SetUpsert(true)

	_, err := r.collection.ReplaceOne(ctx, filter, doc, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert provider service %s: %w", service.ID, err)
	}

	return nil
}

func (r *MongoRepository) UpdateStatus(ctx context.Context, id string, status int) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update status for provider service %s: %w", id, err)
	}

	return nil
}

func (r *MongoRepository) UpdateProviderLogo(ctx context.Context, providerID string, logo *domain.Image) error {
	filter := bson.M{"provider_id": providerID}

	var update bson.M
	if logo == nil {
		update = bson.M{"$unset": bson.M{"provider_logo": 1}}
	} else {
		update = bson.M{"$set": bson.M{"provider_logo": toImageDocument(*logo)}}
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
		update = bson.M{"$unset": bson.M{"units.$.images": 1}}
	} else {
		update = bson.M{"$set": bson.M{"units.$.images": toImageDocuments(images)}}
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update unit images for unit_id %s: %w", unitID, err)
	}

	return nil
}

func (r *MongoRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete provider service %s: %w", id, err)
	}

	return nil
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*domain.ProviderService, error) {
	filter := bson.M{"_id": id}

	var doc ProviderServiceDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: provider service %s", sharederrors.ErrNotFound, id)
		}
		return nil, fmt.Errorf("failed to find provider service %s: %w", id, err)
	}

	return fromDocument(&doc), nil
}
