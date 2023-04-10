package admin_controller

import (
	"context"
	"fmt"
	"github/sahil/atlantic-backend/configs"
	"github/sahil/atlantic-backend/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "admin_users")
var validate_user = validator.New()

// Find User By Email Search
func FindUserByEmail(ctx context.Context, email *string) (models.User, error) {
	var foundUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
	if err != nil {
		return foundUser, err
	}
	return foundUser, err
}

// Find User By ID Search
func FindUserByID(ctx context.Context, id *string) (models.User, error) {
	var foundUser models.User
	err := userCollection.FindOne(ctx, bson.M{"id": id}).Decode(&foundUser)
	if err != nil {
		return foundUser, err
	}
	return foundUser, err
}

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Printf(err.Error())
	}

	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	user.Createdat = now.Format("2006-01-02 15:04:05") // INDIAN STANDARD TIME
	defer cancel()

	//validate_user the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if validationErr := validate_user.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": validationErr.Error()})
	}

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)

	if err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"message": "user already exists with this credentials"})
	}

	// password := HashPassword(user.Password)
	// user.Password = password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	user.Password = string(hashedPassword)

	user.Id = primitive.NewObjectID()
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	// Return response with inserted User ID
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"user":    result,
	})
}

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"data": "Invalid JSON Provided"})
	}

	if validationErr := validate_user.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"data": validationErr.Error()})
	}
	// var user models.User
	foundUser, err := FindUserByEmail(ctx, &user.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Error Occured While Checking For The Email"})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Email Or Password"})
	}
	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"message":   "success",
	// 	"founduser": foundUser.Password,
	// 	"data":      user.Password,
	// })

	// passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
	// // return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// // 	"message": "success",
	// // 	"data":    passwordIsValid,
	// // })
	// if !passwordIsValid {
	// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": msg})
	// }
	// if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Incorrect Password"})
	// }

	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"message": "success",
	// 	"data":    err.Error(),
	// })

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    foundUser.Id.Hex(),
		ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24)), //1 day
	})

	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"message": "success",
	// 	"data":    foundUser.Id.Hex(),
	// })

	token, err := claims.SignedString([]byte(configs.SecretKey()))

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could Not Login"})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    foundUser,
		"token":   token,
	})
}

func User(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.SecretKey()), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "UnAuthenticated"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"message": "success",
	// 	"user":    claims.Issuer,
	// })

	var updatedUser models.User
	_id, err := primitive.ObjectIDFromHex(claims.Issuer)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}
	err = userCollection.FindOne(ctx, bson.M{"id": _id}).Decode(&updatedUser)

	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"user":    updatedUser,
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
