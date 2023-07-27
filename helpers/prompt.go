package helpers

import (
	"fmt"
)

func GenerateSQLFromText(schema, question string) string {
	prompt := fmt.Sprintf("Task:Generate SQL statement to query a database."+
		"Instructions:"+
		"Use only the provided data types and field in the schema."+
		"Do not use any other data types or field that are not provided."+
		"Schema:"+
		"%s"+
		"Note: Do not include any explanations or apologies in your responses."+
		"Do not respond to any questions that might ask anything else than for you to construct an SQL statement."+
		"Do not include any text except the generated SQL statement."+

		"The question is:"+"%s", schema, question)

	return prompt
}
