package ui

import (
	"fmt"
	"jobnbackpack/check/cmd/api"
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

type GoalsInputModel struct {
	focusIndex int
	Goals      []textinput.Model
}

func InitialModel() GoalsInputModel {
	m := GoalsInputModel{
		Goals: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.Goals {
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

		m.Goals[i] = t
	}

	return m
}

func (m GoalsInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m GoalsInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if s == "enter" && m.focusIndex == len(m.Goals) {
				api.WriteToFile(m.Goals)
				return m, tea.Quit
			}

			// submit with journal
			if s == "enter" && m.focusIndex == len(m.Goals)+1 {
				api.WriteToFile(m.Goals)
				//TODO: open next view
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.Goals) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.Goals)
			}

			cmds := make([]tea.Cmd, len(m.Goals))
			for i := 0; i <= len(m.Goals)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.Goals[i].Focus()
					m.Goals[i].PromptStyle = focusedStyle
					m.Goals[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.Goals[i].Blur()
				m.Goals[i].PromptStyle = noStyle
				m.Goals[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *GoalsInputModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Goals))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Goals {
		m.Goals[i], cmds[i] = m.Goals[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m GoalsInputModel) View() string {
	var b strings.Builder

	for i := range m.Goals {
		b.WriteString(m.Goals[i].View())
		if i < len(m.Goals)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.Goals) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
