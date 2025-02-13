package helper

import (
	"context"
	"strings"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetOneDoc[T any](db *mongo.Database, collection string, filter bson.M) (doc T, err error) {
	err = db.Collection(collection).FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		return
	}
	return
}

func GetParam(c *fiber.Ctx) string {
	path := c.Path()
	return path[strings.LastIndex(path, "/")+1:]
}