package helper

import (
    "context"
     "errors"
     "regexp"
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

func ValidateUpdateInput(username, name string) error {
    if len(username) < 3 || len(name) < 3 {
        return errors.New("Username and Name must be at least 3 characters long")
    }
    if match, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username); !match {
        return errors.New("Username can only contain alphanumeric characters and underscores")
    }
    return nil
}

func ValidateLoginInput(username, password string) error {
    re := regexp.MustCompile("^[a-zA-Z0-9_]+$")
    if !re.MatchString(username) {
        return errors.New("invalid username format")
    }

    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    
    return nil
}

func UpdateAdminPassword(db *mongo.Database, col string, adminID string, newPassword string) error {
    collection := db.Collection(col)

    filter := bson.M{"_id": adminID}
    update := bson.M{
        "$set": bson.M{
            "password": newPassword,
        },
    }

    _, err := collection.UpdateOne(context.Background(), filter, update)
    return err
}

