/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"database/sql"
	"fmt"
	"jobnbackpack/check/cmd"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/tursodatabase/go-libsql"
)

func main() {
	connectDB()
	cmd.Execute()
}

func connectDB() {
	dbName := goDotEnvVariable("DB_NAME")
	primaryUrl := goDotEnvVariable("PRIMARY_URL")
	authToken := goDotEnvVariable("AUTH_TOKEN")

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
	)
	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	queryGoals(db)
}

type Goal struct {
	ID          int
	Description string
	Date        time.Time
}

func queryGoals(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM goals")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var goals []Goal

	for rows.Next() {
		var goal Goal

		if err := rows.Scan(&goal.ID, &goal.Description); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		goals = append(goals, goal)
		fmt.Println(goal.ID, goal.Description)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		print("Error loading .env file")
	}

	return os.Getenv(key)
}
