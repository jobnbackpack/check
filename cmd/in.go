/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"jobnbackpack/check/db"
	"jobnbackpack/check/ui/goals"
	"jobnbackpack/check/ui/journal"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// inCmd represents the in command
var inCmd = &cobra.Command{
	Use:   "in",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1e1e2e")).
			Background(lipgloss.Color("#ff8948")).
			PaddingLeft(1).
			PaddingRight(1)
		var runJournal bool

		var database *sql.DB

		fmt.Println(style.Align(lipgloss.Center).Bold(true).MarginTop(1).Render(time.Now().Format("2006-01-02")))
		fmt.Println(style.MarginBottom(1).Render("My main Goals for today:"))
		m, err := tea.NewProgram(goals.InitialModel()).Run()
		if err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
		// Assert the final tea.Model to our local model and print the choice.
		if m, ok := m.(goals.GoalsInputModel); ok && m.Goals[0].Value() != "" {
			database = db.ConnectDB()
			// defer database.Close()

			db.InitGoalsTable(database)
			db.InsertGoal(database, db.Goal{Description: m.Goals[0].Value(), Complete: 0, Date: time.Now().Format("2006-01-02")})

			fmt.Printf("\n---\nYour goals are: %s, %s and %s!\n",
				m.Goals[0].Value(),
				m.Goals[1].Value(),
				m.Goals[2].Value())
			runJournal = m.Journal
		}

		if runJournal {
			m, err = tea.NewProgram(journal.InitialModel()).Run()
			if err != nil {
				fmt.Printf("could not start program: %s\n", err)
				os.Exit(1)
			}
			if m, ok := m.(journal.JournalModel); ok && m.Textarea.Value() != "" {
				fmt.Printf("\n---\nYour Journal: %s\n", m.Textarea.Value())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(inCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
