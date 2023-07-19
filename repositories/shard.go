package repositories

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"github.com/cespare/xxhash/v2"
)

// Data structure
type Data struct {
	// Flexible using map interface
	Attributes map[string]interface{}
}

func check_file_format(filename string) string {
	contentType := GetFileContentType(filename)

	return contentType
}

func GetFileContentType(filename string) string {

	extension := filepath.Ext(filename)

	switch extension {
	case ".csv":
		return "text/csv"
	case ".json":
		return "application/json"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".pdf":
		return "application/pdf"
	case "":
		return "file format not detected"
	default:
		return "application/octet-stream"
	}
}

func readJSON(filename string, c *fiber.Ctx) [][]string {
	f := filename

	content, err := ioutil.ReadFile(f)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to read file",
			"status":  fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	var jsonData []map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to unmarshal json file",
			"status":  fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	columns := make([]string, 0)
	for key := range jsonData[0] {
		columns = append(columns, key)
	}

	data := make([][]string, 0)
	data = append(data, columns)

	for _, item := range jsonData {
		row := make([]string, 0)
		for _, col := range columns {
			value, ok := item[col].(string)
			if !ok {
				continue
			}
			row = append(row, value)
		}

		data = append(data, row)
	}

	return data

}

// Read data from CSV file
func readCSV(filename string, c *fiber.Ctx) [][]string {
	// Open the file
	f, err := os.Open(filename)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to open file",
			"status":  fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	// Close the file
	defer f.Close()

	// Read the data
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to read csv file",
			"status":  fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	return records
}

// Take shard key from the user
// Take shard key from the user
func takeKey(key string, records [][]string, c *fiber.Ctx) (string, error) {
	// Identified column on CSV data
	if len(records) == 0 {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No data found in your file",
			"status":  fiber.StatusBadRequest,
		})
		return "", fmt.Errorf("No data found in your file")
	}

	columns := records[0]

	// See if input is correct
	found := false
	for _, col := range columns {
		if col == key {
			found = true
			break
		}
	}

	if !found {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "shard key not found in available column",
			"status":  fiber.StatusBadRequest,
		})
		return "", fmt.Errorf("Shard key not found in available column")
	}

	return key, nil
}

// Hash the shard key
func hashKey(shardKey string, numShards int) uint32 {
	hash := xxhash.Sum64String(shardKey)

	return uint32(hash % uint64(numShards))
}

func count(Database []string, c *fiber.Ctx) (int, error) {
	if len(Database) < 2 {
		err := fiber.NewError(fiber.StatusBadRequest, "chosen databases must be 2 or more")
		return 0, err
	}

	return len(Database), nil
}


// Perform Horizontal Sharding
func Sharding(records [][]string, shardKey string, numShards int) {
	columns := records[0]
	for _, rec := range records[1:] {
		model := Data{
			Attributes: make(map[string]interface{}),
		}

		shardKeyIndex := 0
		for i, col := range columns {
			if col == shardKey {
				shardKeyIndex = i
				break
			}
		}

		shardKey := rec[shardKeyIndex]

		for i, value := range rec {
			model.Attributes[columns[i]] = value
		}

		shardIndex := hashKey(shardKey, numShards)

		fmt.Printf("Database %d: %v\n", shardIndex, model.Attributes)
	}
}

func HorizontalSharding(Data string, ShardKey string, Database []string, c *fiber.Ctx) error {
	fileFormat := check_file_format(Data)
	if fileFormat == "text/csv" {
		records := readCSV(Data, c)

		key, err := takeKey(ShardKey, records, c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to take shard key",
				"status":  fiber.StatusBadRequest,
				"error": err.Error(),
			})
		}

		databases, err := count(Database, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to count database",
				"status":  fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}

		Sharding(records, key, databases)
	}

	if fileFormat == "application/json" {
		records := readJSON(Data, c)

		key, err := takeKey(ShardKey, records, c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to take shard key",
				"status":  fiber.StatusBadRequest,
				"error": err.Error(),
			})
		}

		databases, err := count(Database, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to count database",
				"status":  fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}

		Sharding(records, key, databases)
	}

	return c.JSON(fiber.Map{
		"message": "sharding is done",
	})
}		


//Perform Vertical Sharding 
func ShardingTwo(records [][]string, shardKey string, numShards int){
	columns := records[0]
	numColumns := len(columns)
	numColumnsPerDB := (numColumns - 1) / numShards // Exclude the shard key column

	for i := 1; i <= numShards; i++ {
		startIndex := (i-1)*numColumnsPerDB + 1 // Start from index 1 to exclude the shard key column
		endIndex := i * numColumnsPerDB

		dbRecords := make([][]string, 0)
		dbRecords = append(dbRecords, columns[:1]) // Include the shard key column

		for _, rec := range records[1:] {
			dbRec := make([]string, 0)
			dbRec = append(dbRec, rec[0]) // Include the shard key value

			for j := startIndex; j <= endIndex; j++ {
				dbRec = append(dbRec, rec[j])
			}

			dbRecords = append(dbRecords, dbRec)
		}

		fmt.Printf("Database %d: %v\n", i, dbRecords)
	}
}

func VerticalSharding(Data string, ShardKey string, Databases []string, c *fiber.Ctx) error {
	fileFormat := check_file_format(Data)
	if fileFormat == "text/csv" {
		records := readCSV(Data, c)

		key, err := takeKey(ShardKey, records, c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to take shard key",
				"status":  fiber.StatusBadRequest,
				"error":   err.Error(),
			})
		}

		databases, err := count(Databases, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to count databases",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		ShardingTwo(records, key, databases)
	}

	if fileFormat == "application/json" {
		records := readJSON(Data, c)

		key, err := takeKey(ShardKey, records, c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to take shard key",
				"status":  fiber.StatusBadRequest,
				"error":   err.Error(),
			})
		}

		databases, err := count(Databases, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to count databases",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		ShardingTwo(records, key, databases)
	}

	return c.JSON(fiber.Map{
		"message": "vertical sharding is done",
	})
}