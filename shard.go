package main

import (
	"encoding/csv"
	"encoding/json"

	"path/filepath"
	"io/ioutil"
	"fmt"
	"os"

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

func GetFileContentType (filename string) string {

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

func readJSON(filename string) [][]string {
	f := filename

	content, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	var jsonData []map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	if err != nil {
		panic(err)
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
func readCSV(filename string) [][]string {
	// Open the file
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	// Close the file
	defer f.Close()

	// Read the data
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	return records
}

// Take shard key from the user
func takeKey(records [][]string) string {
	// Identified column on CSV data
	if len(records) == 0 {
		fmt.Println("No data")
	}

	columns := records[0]

	// Show the available column
	for _, col := range columns {
		fmt.Println(col)
	}

	// Take the column name
	var chooseCol string
	fmt.Print("Select column to be your sharding key: ")
	fmt.Scan(&chooseCol)
	fmt.Println(chooseCol)

	// See if input is correct
	found := false
	for _, col := range columns {
		if col == chooseCol {
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Column not found")
	} else {
		fmt.Println(found)
		return chooseCol
	}

	panic("Something Wrong")
}

// Hash the shard key
func hashKey(shardKey string, numShards int) uint32 {
	hash := xxhash.Sum64String(shardKey)

	return uint32(hash % uint64(numShards))
}

// Perform sharding
func performSharding(records [][]string, shardKey string, numShards int) {
	fmt.Println("Perform sharding...")

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

		fmt.Printf("Shard %d: %v\n", shardIndex, model.Attributes)
	}
}

// Main function
func main() {

	var numofShard int

	file := "test.csv"

	checked_file := check_file_format(file)

	if checked_file == "text/csv" {
		records := readCSV(file)
		
		chooseCol := takeKey(records)

		fmt.Print("How much sharder:")
		fmt.Scan(&numofShard)
		numShard := numofShard

		performSharding(records, chooseCol, numShard)

	} else if checked_file == "application/json" {
		records := readJSON(file)

		chooseCol := takeKey(records)

		fmt.Print("How much sharder:")
		fmt.Scan(&numofShard)
		numShard := numofShard

		performSharding(records, chooseCol, numShard)

	} else {
		fmt.Println("format file not supported!!!")
	}
}
