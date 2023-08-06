package helpers

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/Data-Alchemist-ODS/ods-api/repositories"
	"github.com/go-sql-driver/mysql"
	"github.com/sashabaranov/go-openai"
)

func SaveToTiDB(data []repositories.Data, serverName, user, password, database string) error {
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		// ServerName: "gateway01.eu-central-1.prod.aws.tidbcloud.com",
		ServerName: serverName,
	})

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:4000)/%s?tls=tidb", user, password, serverName, database)
	// db, err := sql.Open("mysql", "4MXeBRmXXzc7uqt.root:<your_password>@tcp(gateway01.eu-central-1.prod.aws.tidbcloud.com:4000)/test?tls=tidb")
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	prompt := fmt.Sprintf(
		"Task:Generate SQL statement to create a table."+
			"Instructions:"+
			"Use only the provided variable names."+
			"Do not use any other variables that are not provided."+
			"Variables:"+
			"%s"+
			"Note: Do not include any explanations or apologies in your responses."+
			"Do not respond to any questions that might ask anything else than for you to construct an SQL statement."+
			"Do not include any text except the generated SQL statement."+
			"Table name: shard", reflect.ValueOf(data[0].Attributes).MapKeys(),
	)

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return err
	}

	fmt.Println(resp.Choices[0].Message.Content)

	_, err = db.Exec("DROP TABLE IF EXISTS shard")
	if err != nil {
		log.Fatal("failed to execute query", err)
	}

	_, err = db.Exec(resp.Choices[0].Message.Content)
	if err != nil {
		log.Fatal("failed to execute query", err)
	}

	return nil
}
