package controllers

import (
	//default modules
	"context"

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/entity"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
)

// UserController is a contract what this controller can do
type UserController interface {
	//GET HANDLER
	GetAllUser(c *fiber.Ctx) error
	GetOneUser(c *fiber.Ctx) error

	//POST HANDLER
	RegisterUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error

	//PUT HANDLER
	UpdateUser(c *fiber.Ctx) error

	//DELETE HANDLER
	DeleteOneUser(c *fiber.Ctx) error
}

// userController is a struct that represent the UserController contract
type userController struct{}

// NewUserController is the constructor
func NewUserController() UserController {
	return &userController{}
}

//GET REQUEST CONTROLLER
//Get All User
func (controller *userController) GetAllUser(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "User")

	var users []entity.User

	cursor, err := collection.Find(context.Background(), options.Find())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get user",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	if err := cursor.All(context.Background(), &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to decode user",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success get all user",
		"status":  fiber.StatusOK,
		"records": users,
	})
}

//Get One User From Database By Id Params
func (controller *userController) GetOneUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid id format",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "User")

	var user entity.User

	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "user not found",
				"status":  fiber.StatusNotFound,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get user",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success get one user",
		"status":  fiber.StatusOK,
		"record":  user,
	})
}

//POST REQUEST CONTROLLER
//Create User For Register
func (controller *userController) RegisterUser(c *fiber.Ctx) error {
	var request request.UserCreateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse json",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "User")

	// CHECK IF EMAIL ALREADY EXIST
	userExist := collection.FindOne(context.Background(), bson.M{"email": request.Email})
	if userExist.Err() == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email already exist",
			"status":  fiber.StatusBadRequest,
		})
	}

	user := entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
	user.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create user",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success create user",
		"status":  fiber.StatusOK,
		"record":  user,
	})
}

//Function For User Login
func (controller *userController) LoginUser(c *fiber.Ctx) error {
	var request request.UserLoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse json",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "User")

	var user entity.User

	err := collection.FindOne(context.Background(), bson.M{"email": request.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid email or password",
				"status":  fiber.StatusUnauthorized,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to find user",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	if user.Password != request.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid email or password",
			"status":  fiber.StatusUnauthorized,
		})
	}

	return c.JSON(fiber.Map{
		"user_id": user.ID.Hex(),
		"message": "success login user",
		"status":  fiber.StatusOK,
	})
}

//PUT REQUEST CONTROLLER
//Update User By Requested Id Params
func (controller *userController) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid id format",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "User")

	var request request.UserCreateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse json",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	update := bson.M{
		"$set": bson.M{},
	}

	if request.Name != "" {
		update["$set"].(bson.M)["name"] = request.Name
	}

	if request.Email != "" {
		update["$set"].(bson.M)["email"] = request.Email
	}

	if request.Password != "" {
		update["$set"].(bson.M)["password"] = request.Password
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update user",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success update user",
		"status":  fiber.StatusOK,
		"record":  update,
	})
}

//DELETE REQUEST CONTROLLER
//Delete One User By Id Params
func (controller *userController) DeleteOneUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid id format",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "User")

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete user",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "user not found",
			"status":  fiber.StatusNotFound,
		})
	}

	return c.JSON(fiber.Map{
		"message": "success delete user",
		"status":  fiber.StatusOK,
	})
}
