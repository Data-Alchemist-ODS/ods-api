package controllers

import (
	//default modules

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//mongoDB modules

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/database"
)

// UserController is a contract what this controller can do
type UserController interface {
	//GET HANDLER
	//POST HANDLER
	//PATCH HANDLER
	//DELETE HANDLER
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
	defer db .Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "User")

	var users []entity.User

	cursor, err := collection.Find(context.Background(), options.Find())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get user",
			"status": fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	if err := cursor.All(context.Background(), &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to decode user",
			"status": fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success get all user",
		"status": fiber.StatusOK,
		"records": users,
	})
}