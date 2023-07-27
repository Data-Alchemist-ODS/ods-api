package repositories

import (
	//default modules
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//third party modules
	"github.com/cespare/xxhash/v2"
)

// Data structure
type (
	Data struct {
		// Flexible using map interface
		Attributes map[string]interface{}
	}

	DataSharded map[int][]Data
)

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
	file, err := c.FormFile(filename)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error when opening file",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	f, err := file.Open()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error when opening file",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to read file",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	var jsonData []map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to unmarshal json file",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
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
	file, err := c.FormFile(filename)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to open file",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	f, err := file.Open()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to open file",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
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
			"error":   err.Error(),
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

// Perform Horizontal Sharding
func shardForHorizontal(records [][]string, shardKey string, numShards int) [][]Data {
	columns := records[0]

	// create an array of [][]Data with length of numShards
	result := make([][]Data, numShards)

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
		result[int(shardIndex)] = append(result[int(shardIndex)], model)
	}

	return result
}

func HorizontalSharding(Data string, ShardKey string, Database []string, c *fiber.Ctx) ([][]Data, error) {

	file, err := c.FormFile(Data)

	if err != nil {
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"message": "failed to open file using form file",
		// 	"status":  fiber.StatusInternalServerError,
		// 	"error":   err.Error(),
		// })
		return nil, err
	}

	fileFormat := file.Header.Get("Content-Type")

	fmt.Println("File Format:", fileFormat)

	if fileFormat == "text/csv" {

		records := readCSV(Data, c)

		key, err := takeKey(ShardKey, records, c)

		if err != nil {
			log.Default().Println("failed to take shard key:", err.Error())
			// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			// 	"message": "failed to take shard key",
			// 	"status":  fiber.StatusBadRequest,
			// 	"error":   err.Error(),
			// })
			return nil, err
		}

		databases := len(Database)
		fmt.Println("databases:", databases)
		if databases < 2 {
			// log.Default().Println("databases must be more than 1:", err.Error())
			// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			// 	"message": "databases must be more than 1",
			// 	"status":  fiber.StatusBadRequest,
			// })
			return nil, err
		}

		return shardForHorizontal(records, key, databases), nil
	}

	if fileFormat == "application/json" {
		records := readJSON(Data, c)

		key, err := takeKey(ShardKey, records, c)
		if err != nil {
			// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			// 	"message": "failed to take shard key",
			// 	"status":  fiber.StatusBadRequest,
			// 	"error":   err.Error(),
			// })
			return nil, err
		}

		databases := len(Database)
		if databases < 2 {
			// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			// 	"message": "databases must be more than 1",
			// 	"status":  fiber.StatusBadRequest,
			// })
			return nil, fmt.Errorf("databases must be more than 1")
		}

		return shardForHorizontal(records, key, databases), nil
	}

	// return c.JSON(fiber.Map{
	// 	"message": "shard is done",
	// 	"status":  fiber.StatusOK,
	// })
	return nil, fmt.Errorf("file format not supported")
}

// //THE LOGIC IS FALSE HERE!!!
// // Vertical Sharding
// func shardForVertical(records [][]string, shardKey string, numShards int) {
// 	shardsMap := make([]map[string]string, numShards)

// 	for i := 0; i < numShards; i++ {
// 		shardsMap[i] = make(map[string]string)
// 	}

// 	// Find the index of the shardKey column
// 	shardKeyIndex := -1
// 	for i, columnName := range records[0] {
// 		if columnName == shardKey {
// 			shardKeyIndex = i
// 			break
// 		}
// 	}

// 	if shardKeyIndex == -1 {
// 		fmt.Println("Shard key not found in the columns.")
// 		return
// 	}

// 	// Iterate over rows (skip header)
// 	for rowIndex := 1; rowIndex < len(records); rowIndex++ {
// 		row := records[rowIndex]
// 		if len(row) <= shardKeyIndex {
// 			fmt.Println("Invalid record format. Shard key index out of range.")
// 			return
// 		}
// 		shardKey := row[shardKeyIndex]

// 		// Create the shard value using the shard key and other column values
// 		shardValue := ""
// 		for colIndex, columnName := range records[0] {
// 			if colIndex == shardKeyIndex {
// 				shardValue = fmt.Sprintf("%s: %s", columnName, shardKey)
// 			} else {
// 				shardValue += fmt.Sprintf(" %s: %s", columnName, row[colIndex])
// 			}
// 		}

// 		// Calculate the shard index based on the shard key using consistent hashing
// 		shardIndex := consistentHash(shardKey, numShards)

// 		shardsMap[shardIndex][shardValue] = fmt.Sprintf("map[%s]", shardValue)
// 	}

// 	// Output
// 	for i, shardMap := range shardsMap {
// 		fmt.Printf("Shard %d:\n", i)
// 		for _, value := range shardMap {
// 			fmt.Println(value)
// 		}
// 		fmt.Println()
// 	}
// }

// // Custom function for consistent hashing
// func consistentHash(key string, numShards int) int {
// 	hash := fnv.New32a()
// 	_, _ = hash.Write([]byte(key))
// 	hashValue := int(hash.Sum32())
// 	shardIndex := hashValue % numShards
// 	if shardIndex < 0 {
// 		shardIndex = -shardIndex
// 	}
// 	return shardIndex
// }

// func VerticalSharding(Data string, ShardKey string, Databases []string, c *fiber.Ctx) error {
// 	fileFormat := check_file_format(Data)
// 	if fileFormat == "text/csv" {
// 		records := readCSV(Data, c)

// 		key, err := takeKey(ShardKey, records, c)
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"status": fiber.StatusBadRequest,
// 				"error":  err.Error(),
// 			})
// 		}

// 		databases, err := count(Databases)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "failed to count databases",
// 				"status":  fiber.StatusInternalServerError,
// 				"error":   err.Error(),
// 			})
// 		}

// 		shardForVertical(records, key, databases)
// 	}

// 	if fileFormat == "application/json" {
// 		records := readJSON(Data, c)

// 		key, err := takeKey(ShardKey, records, c)
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"message": "failed to take shard key",
// 				"status":  fiber.StatusBadRequest,
// 				"error":   err.Error(),
// 			})
// 		}

// 		databases, err := count(Databases)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "failed to count databases",
// 				"status":  fiber.StatusInternalServerError,
// 				"error":   err.Error(),
// 			})
// 		}

// 		shardForVertical(records, key, databases)
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "vertical sharding is done",
// 	})
// }
