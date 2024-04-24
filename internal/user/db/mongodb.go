package db

import (
	"context"
	"errors"
	"fmt"

	"example.com/m/internal/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
}

func (d *db) Create(ctx context.Context, user user.CreateUserDTO) (string, error) {
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("Failed to create user due to error: %v", err)
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to convert objectid to hex, oid: %s", oid)
}

func (d *db) FindOne(ctx context.Context, id string) (u user.CreateUserDTO, err error) {
	// oid, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return u, fmt.Errorf("failed to convert hex to objectid, hex: %s", id)
	// }

	filter := bson.M{"guid": id}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, fmt.Errorf("not found")
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to err: %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user(id:^%s)from db due to err: %v", id, err)
	}
	return u, nil
}

func (d *db) FindRefresh(ctx context.Context, refToken string) (u user.CreateUserDTO, err error) {
	filter := bson.M{"refreshtoken": refToken}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, fmt.Errorf("not found")
		}
		return u, fmt.Errorf("failed to find one user by refreshtoken: %s due to err: %v", refToken, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user(refreshtoken:^%s)from db due to err: %v", refToken, err)
	}
	return u, nil
}

func (d *db) Update(ctx context.Context, user user.CreateUserDTO) error {
	// oid, err := primitive.ObjectIDFromHex(user.GUID)
	// if err != nil {
	// 	return fmt.Errorf("failed to convert user ID to objectid, ID: %s", user.GUID)
	// }

	filter := bson.M{"guid": user.GUID}
	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user, error: %v", err)
	}

	var updateUserObj bson.M
	if err = bson.Unmarshal(userBytes, &updateUserObj); err != nil {
		return fmt.Errorf("failed to unmarshal user bytes, error: %v", err)
	}

	delete(updateUserObj, "guid")

	update := bson.M{
		"$set": &updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query, error: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to objectid, ID: %s", id)
	}

	filter := bson.M{"guid": oid}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query, error: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

func NewStorage(database *mongo.Database, collection string) user.Storage {
	return &db{
		collection: database.Collection(collection),
	}
}
