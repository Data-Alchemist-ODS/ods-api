package helpers

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"

	"github.com/Data-Alchemist-ODS/ods-api/repositories"
	"github.com/go-sql-driver/mysql"
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

	result, err := db.Query("SELECT * FROM `fortune500_2018_2022`")

	if err != nil {
		log.Fatal("failed to execute query", err)
	}

	fmt.Println(result.Scan())

	return nil
}
