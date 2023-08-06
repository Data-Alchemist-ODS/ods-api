package controllers

import (
	//default modules
	"context"
	"fmt"
	"strings"

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/config"
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/entity"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"

	//third party modules
	api "github.com/sashabaranov/go-openai"
)

type NotesController interface {
	//GET HANDLER
	GetAllNotes(c *fiber.Ctx) error

	//POST HANDLER
	CreateNewNotes(c *fiber.Ctx) error

	//DELETE HANDLER
	DeleteNotesHistory(c *fiber.Ctx) error
}

type notesController struct {}

func NewNotesController() NotesController {
	return &notesController{}
}

func (controller *notesController) GetAllNotes(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "Notes")

	var notes []entity.Notes

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed on finding notes",
			"status": fiber.StatusInternalServerError,
		})
	}

	err = cursor.All(context.Background(), &notes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed when decoding notes",
			"status": fiber.StatusInternalServerError,
		})
	}
	return c.JSON(fiber.Map{
		"message": "success on finding notes",
		"status": fiber.StatusOK,
		"data": notes,
	})
}

func (controller *notesController) CreateNewNotes(c *fiber.Ctx) error {
    inputData := new(request.NotesReq)
    if err := c.BodyParser(inputData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "failed to parse request body",
            "status":  fiber.StatusBadRequest,
            "error":   err.Error(),
        })
    }

    notesSections := strings.Split(inputData.TextInput, "\n\n")

    var processedNotes []entity.Notes

    for _, section := range notesSections {
        lines := strings.Split(section, "\n")
        if len(lines) > 0 {
            dateAndDescription := lines[0]
            description := strings.Join(lines[1:], "\n")
            processedNotes = append(processedNotes, entity.Notes{
                Date:        dateAndDescription,
                Description: description,
            })
        }
    }

    client := api.NewClient(config.LoadAPIKey())

    for i, note := range processedNotes {
        prompt := fmt.Sprintf("Input:\n%s\n\nGenerate an analyze based on the following description:", note.Description)

        resp, err := client.CreateChatCompletion(
            context.Background(),
            api.ChatCompletionRequest{
                Model: api.GPT3Dot5Turbo,
                Messages: []api.ChatCompletionMessage{
                    {
                        Role:    api.ChatMessageRoleSystem,
                        Content: "You are a helpful assistant that analyze text based on note description.",
                    },
                    {
                        Role:    api.ChatMessageRoleUser,
                        Content: prompt,
                    },
                },
            },
        )
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "failed to create new notes",
                "status":  fiber.StatusInternalServerError,
                "error":   err.Error(),
            })
        }

        processedNotes[i].Analysis = resp.Choices[0].Message.Content
    }

    db := database.ConnectDB()
    defer db.Disconnect(context.Background())

    collection := database.GetCollection(database.GetDB(), "Notes")

    for _, note := range processedNotes {
        _, err := collection.InsertOne(context.Background(), note)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "failed to insert notes",
                "status":  fiber.StatusInternalServerError,
                "error":   err.Error(),
            })
        }
    }

    return c.JSON(fiber.Map{
        "message": "success inserting notes to database",
        "status":  fiber.StatusOK,
        "data":    processedNotes,
    })
}

func (controller *notesController) DeleteNotesHistory(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "failed to parse id",
            "status":  fiber.StatusBadRequest,
            "error":   err.Error(),
        })
    }

    db := database.ConnectDB()
    defer db.Disconnect(context.Background())

    client := database.GetDB()
    collection := database.GetCollection(client, "Notes")

    result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "failed to delete notes",
            "status":  fiber.StatusInternalServerError,
            "error":   err.Error(),
        })
    }

    if result.DeletedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "notes not found",
            "status":  fiber.StatusNotFound,
        })
    }

    return c.JSON(fiber.Map{
        "message": "success deleting notes",
        "status":  fiber.StatusOK,
    })
}
