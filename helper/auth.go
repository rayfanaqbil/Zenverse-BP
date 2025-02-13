package helper

import (
    "context"
    "errors"
    "github.com/gofiber/fiber/v2"
    "regexp"
    "go.mongodb.org/mongo-driver/bson"
    "time"
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
        return errors.New("username and name must be at least 3 characters long")
    }
    if match, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username); !match {
        return errors.New("username can only contain alphanumeric characters and underscores")
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

func UpdateAdminPassword(db *mongo.Database, collectionName, adminID, newHashedPassword string) error {
    collection := db.Collection(collectionName)

    objID, err := primitive.ObjectIDFromHex(adminID)
    if err != nil {
        return err
    }

    filter := bson.M{"_id": objID} 
    update := bson.M{"$set": bson.M{"password": newHashedPassword, "updated_at": time.Now()}}

    _, err = collection.UpdateOne(context.TODO(), filter, update)
    return err
}

func GetSecretFromHeader(c *fiber.Ctx) string {
	secret := c.Get("secret") 
	if secret == "" {
		secret = c.Get("Secret")
	}
	return secret
}
