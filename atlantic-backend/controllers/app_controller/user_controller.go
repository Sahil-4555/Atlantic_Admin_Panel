package appcontrollers

import (
	"context"
	"net/http"
	"time"

	"github/sahil/atlantic-backend/configs"
	"github/sahil/atlantic-backend/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "app_users")
var validate_user = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.AppUsers
	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	user.Createdat = now.Format("2006-01-02 15:04:05") // INDIAN STANDARD TIME
	user.Updatedat = now.Format("2006-01-02 15:04:05") // INDIAN STANDARD TIME
	defer cancel()

	//validate_user the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	//use the validator library to validate_user required fields
	if validationErr := validate_user.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": validationErr.Error()})
	}

	newUser := models.AppUsers{
		Id:        primitive.NewObjectID(),
		Uid:       user.Uid,
		Email:     user.Email,
		Photourl:  user.Photourl,
		Name:      user.Name,
		Createdat: user.Createdat,
		Updatedat: user.Updatedat,
	}

	err := userCollection.FindOne(ctx, bson.M{"uid": user.Uid}).Decode(&user)

	if err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"message": "user already exists"})
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"user":    result,
	})
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("uid")
	var user models.AppUsers
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"uid": userId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"user":    user,
	})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("uid")
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"uid": userId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error()})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			fiber.Map{"message": "User with specified ID Not Found!"},
		)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "success"})
}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("uid")
	var user models.AppUsers
	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	user.Updatedat = now.Format("2006-01-02 15:04:05") // INDIAN STANDARD TIME
	defer cancel()

	//validate_user the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	//use the validator library to validate_user required fields
	if validationErr := validate_user.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": validationErr.Error()})
	}

	update := bson.M{
		"uid":       user.Uid,
		"email":     user.Email,
		"photourl":  user.Photourl,
		"name":      user.Name,
		"updatedat": user.Updatedat,
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"uid": userId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	var updatedUser models.AppUsers

	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"uid": user.Uid}).Decode(&updatedUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	} else {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "User Not Found"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"user":    updatedUser,
	})
}

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.AppUsers
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.AppUsers
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"message": "success",
			"user":    users,
		},
	)
}
