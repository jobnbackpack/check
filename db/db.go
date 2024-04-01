package db

import (
	"database/sql"
	"fmt"
	"jobnbackpack/check/util"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

func ConnectDB() *sql.DB {
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

	return sql.OpenDB(connector)
}

type Goal struct {
	ID          int
	Description string
	Date        string
	Complete    int
}

func InitGoalsTable(db *sql.DB) {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS goals ( id INTEGER PRIMARY KEY AUTOINCREMENT, description TEXT NOT NULL, complete INTEGER NOT NULL, date TEXT NOT NULL)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
}

func InsertGoal(db *sql.DB, goal Goal) {
	_, err := db.Query(fmt.Sprintf("INSERT INTO goals (description, complete, date) VALUES(%s, %d, %s)", goal.Description, goal.Complete, goal.Date))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
}

func QueryGoals(db *sql.DB) {
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
