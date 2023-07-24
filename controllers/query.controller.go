package query

import (
	//default modules
	"context"
	"fmt"
	"log"

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//third-party modules
	api "github.com/sashabaranov/go-openai"

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/config"
	"github.com/Data-Alchemist-ODS/ods-api/entity"
	"github.com/Data-Alchemist-ODS/ods-api/request"
)

//The contract
type NaturalQueryController interface {
	CreateNaturalQuery(c *fiber.Ctx) error
}

//Represent the contract
type naturalqueryController struct{}

//The constructer
func NewNaturalQueryController() NaturalQueryController {
	return &naturalqueryController{}
}

//Create a new natural query
func (controller *naturalqueryController) CreateNaturalQuery(c *fiber.Ctx) error {
	req := new()
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse request body",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	client := api.NewClient(config.LoadAPIKey())
	resp, err := client.CreateChatCompletion(
		context.Background(),
		api.ChatCompletionRequest{
			Model: api.GPT3Dot5Turbo,
			Messages: []api.ChatCompletionMessage{
				Role: api.ChatMessageRoleUser,
				Content: req.Prompt,
			},
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create natural query",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"status":  fiber.StatusOK,
		"query": entity.queryResp(
			Response: resp.Choices[0].Message.Content,
		),
	})
}