package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/cespare/xxhash/v2" // Import package for XXHash algorithm
)

// Data structure
type Data struct {
	// Flexible using map interface
	Attributes map[string]interface{}
}

// Read data from CSV file
func readData() [][]string {
	// Open the file
	f, err := os.Open("test.csv")
	if err != nil {
		fmt.Println(err)
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
	// Read data
	records := readData()

	// Take key
	chooseCol := takeKey(records)

	fmt.Print("How much sharder:")
	fmt.Scan(&numofShard)
	numShard := numofShard

	performSharding(records, chooseCol, numShard)
}
