package persistence

import (
	"context"
	"errors"

	"poc-auth-svc/internal/domain/entities"
	err_domain "poc-auth-svc/internal/domain/errors"
	"poc-auth-svc/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) repositories.UserRepository {
	return &mongoUserRepository{
		collection: db.Collection("users"),
	}
}

// Create implements repositories.UserRepository.
func (m *mongoUserRepository) Create(ctx context.Context, user *entities.User) error {
	_, err := m.collection.InsertOne(ctx, user)
	return err
}

// GetByEmail implements repositories.UserRepository.
func (m *mongoUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := m.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(err_domain.GetMessage(err_domain.UserNotFound))
		}
		return nil, err
	}
	return &user, nil
}

// GetByID implements repositories.UserRepository.
func (m *mongoUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	if err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(err_domain.GetMessage(err_domain.UserNotFound))
		}
		return nil, err
	}
	return &user, nil
}

// Update implements repositories.UserRepository.
func (m *mongoUserRepository) Update(ctx context.Context, user *entities.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete implements repositories.UserRepository.
func (m *mongoUserRepository) Delete(ctx context.Context, id string) error {
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
