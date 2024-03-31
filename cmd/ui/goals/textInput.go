package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// #1e1e2e great bg color
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8948")).Bold(true)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle.Copy()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	focusIndex int
	goals      []textinput.Model
}

func InitialModel() model {
	m := model{
		goals: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.goals {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Goal 1"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Goal 2"
		case 2:
			t.Placeholder = "Goal 3"
		}

		m.goals[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.goals) {
				var x = []byte{}

				for i := 0; i < len(m.goals); i++ {
					b := []byte("Goal " + fmt.Sprint(i+1) + " " + m.goals[i].Value() + "\n")
					for j := 0; j < len(b); j++ {
						x = append(x, b[j])
					}
				}
				os.WriteFile("db.txt", x, 0644)
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.goals) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.goals)
			}

			cmds := make([]tea.Cmd, len(m.goals))
			for i := 0; i <= len(m.goals)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.goals[i].Focus()
					m.goals[i].PromptStyle = focusedStyle
					m.goals[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.goals[i].Blur()
				m.goals[i].PromptStyle = noStyle
				m.goals[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.goals))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.goals {
		m.goals[i], cmds[i] = m.goals[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.goals {
		b.WriteString(m.goals[i].View())
		if i < len(m.goals)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.goals) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
