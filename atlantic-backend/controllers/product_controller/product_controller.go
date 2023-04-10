package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"net/http"
	"time"

	"github/sahil/atlantic-backend/configs"
	"github/sahil/atlantic-backend/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	resize "github.com/nfnt/resize"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "products_atlantic")
var validate_user = validator.New()

func CreateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var product models.Product
	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	product.Createdat = now.Format("2006-01-02 15:04:05") // INDIAN STANDARD TIME
	product.Updatedat = now.Format("2006-01-02 15:04:05") // INDIAN STANDARD TIME
	defer cancel()

	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to fetch Data",
		})
	}
	// imageFile, err := ioutil.ReadAll(c.Body("image"))
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed To Read Image",
		})
	}
	src, err := file.Open()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Error while reading image",
		})
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Error while Decoding the image from file to image.Image",
		})
	}

	// Compress image if it is larger than 1MB
	maxSize := int64(1 * 1024 * 1024)
	if file.Size > maxSize {
		newImg := resize.Resize(0, 500, img, resize.Lanczos3)
		img = newImg
	}

	// Encode the image to []byte
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Error while Encode the image to []byte",
		})
	}
	imgBase64Str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	product.Image = imgBase64Str

	// Create Product object
	newProduct := models.Product{
		Id:          primitive.NewObjectID(),
		Productid:   product.Productid,
		Image:       product.Image,
		Title:       product.Title,
		Price:       product.Price,
		Size:        product.Size,
		Description: product.Description,
		Color:       product.Color,
		Createdat:   product.Createdat,
		Updatedat:   product.Updatedat,
	}

	result, err := userCollection.InsertOne(ctx, newProduct)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error Inserting User Into MongoDB",
		})
	}

	// Return response with inserted User ID
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"id":      result.InsertedID,
	})
}

func GetAProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	productID := c.Params("productID")
	var product models.Product
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"productid": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "User Not Found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error Fetching User From MongoDB",
		})
	}

	//Return response with product and image
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"product": product,
	})
}

func GetAProductPhoto(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	productID := c.Params("productID")
	var product models.Product
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"productid": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "User Not Found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error Fetching User From MongoDB",
		})
	}

	imgBase64Bytes, err := base64.StdEncoding.DecodeString(product.Image)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while Decode the base64 encoded string to []byte",
		})
	}

	// Create an io.Reader from the decoded []byte
	imgReader := bytes.NewReader(imgBase64Bytes)

	// Decode the image from io.Reader to image.Image
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while Decode the image from io.Reader to image.Image",
		})
	}

	// Encode the image to []byte
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, img)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while Encode the image to []byte",
		})
	}

	// Set the response header to indicate that the response is an image
	c.Set("Content-Type", "image/png")

	//Send the response with the image
	return c.Send(buffer.Bytes())

}

func UpdateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productID := c.Params("productID")

	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		return c.Status(http.StatusNotFound).JSON(fiber.Map{
	// 			"message": "User Not Found",
	// 		})
	// 	}
	// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 		"message": "Error Fetching User From MongoDB",
	// 	})
	// }

	var product models.Product

	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	product.Updatedat = now.Format("2006-01-02 15:04:05") // INDIAN STANDARD TIME

	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to fetch Data",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"productId":   product.Productid,
			"title":       product.Title,
			"price":       product.Price,
			"size":        product.Size,
			"description": product.Description,
			"color":       product.Color,
			"updatedat":   product.Updatedat,
		},
	}

	// Check if an image file is attached
	file, err := c.FormFile("image")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Error While Reading Image",
			})
		}
		defer src.Close()

		img, _, err := image.Decode(src)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Error While Decoding The Image From File To Image.",
			})
		}

		// Compress image if it is larger than 1MB
		maxSize := int64(1 * 1024 * 1024)
		if file.Size > maxSize {
			newImg := resize.Resize(0, 500, img, resize.Lanczos3)
			img = newImg
		}

		// Encode the image to []byte
		buffer := new(bytes.Buffer)
		err = jpeg.Encode(buffer, img, nil)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Error while Encode the image to []byte",
			})
		}
		imgBase64Str := base64.StdEncoding.EncodeToString(buffer.Bytes())
		update["$set"].(bson.M)["image"] = imgBase64Str
	}

	filter := bson.M{"productid": productID}
	// if err != nil {
	// 	return c.Status(http.StatusNotFound).JSON(fiber.Map{
	// 		"message": "User Not Found",
	// 	})
	// }
	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error Updating User In MongoDB",
		})
	}

	var updateProduct_ models.Product

	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"productid": productID}).Decode(&updateProduct_)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error While Fetching Updated User"})
		}

	} else {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "User Not Found!"})
	}

	// Return response with updated User ID
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    updateProduct_,
	})
}

func DeleteAProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	productID := c.Params("productID")
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"productid": productID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "User With Specified ID Not Found!"},
		)
	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "success"},
	)
}

func GetAllProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var products []models.Product
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleProduct models.Product
		if err = results.Decode(&singleProduct); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		products = append(products, singleProduct)
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"message": "success",
			"product": products,
		},
	)
}

// func UpdateUser(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	userId := c.Params("userId")
// 	defer cancel()

// 	// Get User ID from request params

// 	objId, _ := primitive.ObjectIDFromHex(userId)

// 	// Get User object from MongoDB
// 	var user models.User
// 	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
// 	if err != nil {
// 		return c.Status(http.StatusNotFound).JSON(fiber.Map{
// 			"message": "User not found",
// 		})
// 	}

// 	// Parse request body into a new User object
// 	var updatedUser_ models.User
// 	if err := c.BodyParser(&updatedUser_); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Failed to fetch data from request body",
// 		})
// 	}

// 	// Update User fields if provided in the request body
// 	if updatedUser_.Username != "" {
// 		user.Username = updatedUser_.Username
// 	}
// 	if updatedUser_.Name != "" {
// 		user.Name = updatedUser_.Name
// 	}
// 	if updatedUser_.Email != "" {
// 		user.Email = updatedUser_.Email
// 	}
// 	if updatedUser_.Phone != "" {
// 		user.Phone = updatedUser_.Phone
// 	}
// 	if updatedUser_.Password != "" {
// 		user.Password = updatedUser_.Password
// 	}

// 	// Update User image if provided in the request body
// 	if updatedUser_.Image != "" {

// 		// Decode the base64-encoded image string to image.Image
// 		imgBytes, err := base64.StdEncoding.DecodeString(updatedUser_.Image)
// 		if err != nil {
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 				"message": "Failed to decode image from base64 string",
// 			})
// 		}
// 		imgReader := bytes.NewReader(imgBytes)
// 		img, _, err := image.Decode(imgReader)
// 		if err != nil {
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 				"message": "Failed to decode image bytes to image",
// 			})
// 		}

// 		// Encode the image to JPEG and store the base64-encoded string in User object
// 		buffer := new(bytes.Buffer)
// 		err = jpeg.Encode(buffer, img, nil)
// 		if err != nil {
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 				"message": "Failed to encode image to JPEG format",
// 			})
// 		}
// 		user.Image = base64.StdEncoding.EncodeToString(buffer.Bytes())
// 	}

// 	// Update the Updatedat field with the current time
// 	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
// 	user.Updatedat = now.Format("2006-01-02 15:04:05")

// 	// Update User object in MongoDB
// 	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": user})
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to update user in MongoDB",
// 		})
// 	}

// 	var updateUser models.User

// 	if result.MatchedCount == 1 {
// 		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updateUser)
// 		if err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Error While Fetching Updated User"})
// 		}

// 	} else {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "User Not Found!"})
// 	}
// 	// Return response with number of documents updated
// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"data":              updateUser,
// 		"message":           "User updated successfully",
// 		"documents_updated": result.ModifiedCount,
// 	})
// }
