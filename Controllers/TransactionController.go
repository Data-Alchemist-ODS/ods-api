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

func GetALLTranscation (ctx *fiber.Ctx) error {
	var transaction []
}