package utils

import (
	"context"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserWithNotes(db *mongo.Database, userID primitive.ObjectID) (*models.UserWithNotes, error) {
	collection := db.Collection("users")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": userID,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "notes",
				"localField":   "noteIds",
				"foreignField": "_id",
				"as":           "notes",
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []models.UserWithNotes
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &results[0], nil
}

func GetAllUserNotes(db *mongo.Database, userID primitive.ObjectID) ([]models.Note, error) {
	collection := db.Collection("users")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": userID,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "notes",
				"localField":   "noteIds",
				"foreignField": "_id",
				"as":           "notes",
			},
		},
		{
			"$project": bson.M{
				"notes": 1,
				"_id":   0,
			},
		},
		{
			"$unwind": "$notes",
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$notes",
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var notes []models.Note
	if err = cursor.All(context.Background(), &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

func AddNoteToUser(db *mongo.Database, userID primitive.ObjectID, noteID primitive.ObjectID) error {
	collection := db.Collection("users")

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$addToSet": bson.M{"noteIds": noteID}},
	)

	return err
}

func RemoveNoteFromUser(db *mongo.Database, userID primitive.ObjectID, noteID primitive.ObjectID) error {
	collection := db.Collection("users")

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$pull": bson.M{"noteIds": noteID}},
	)

	return err
}
