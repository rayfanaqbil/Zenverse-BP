package helper

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/rayfanaqbil/zenverse-BE/v2/model"
)

func GetAdminByID(db *mongo.Database, col string, adminID string) (*model.Admin, error) {
    var admin model.Admin
    collection := db.Collection("Admin")

    objectID, err := primitive.ObjectIDFromHex(adminID)
    if err != nil {
        return nil, err
    }

    filter := bson.M{"_id": objectID}
    err = collection.FindOne(context.Background(), filter).Decode(&admin)
    if err != nil {
        return nil, err
    }

    return &admin, nil
}
