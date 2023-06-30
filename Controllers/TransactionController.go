package Controllers

import (
	"encoding/csv"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cespare/xxhash/v2"
)

func GetALLTranscation (c *fiber.Ctx) error {
	coll := database.GetCollection("")
}