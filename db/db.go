package db

import (
	"database/sql"
	"fmt"
	"jobnbackpack/check/util"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

func ConnectDB() {
	dbName := util.GoDotEnvVariable("DB_NAME")
	primaryUrl := util.GoDotEnvVariable("PRIMARY_URL")
	authToken := util.GoDotEnvVariable("AUTH_TOKEN")

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

	initGoalsTable(db)
	queryGoals(db)
}

type Goal struct {
	ID          int
	Description string
	Date        string
	Complete    bool
}

func initGoalsTable(db *sql.DB) {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS tasks ( id INTEGER PRIMARY KEY AUTOINCREMENT, description TEXT NOT NULL, complete INTEGER NOT NULL, date TEXT NOT NULL)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
}

func queryGoals(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM tasks")
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
