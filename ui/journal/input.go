package journal

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type JournalModel struct {
	Textarea textarea.Model
	err      error
}

func InitialModel() JournalModel {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time..."
	ti.Focus()

	return JournalModel{
		Textarea: ti,
		err:      nil,
	}
}

func (m JournalModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m JournalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.Textarea.Focused() {
				m.Textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.Textarea.Focused() {
				cmd = m.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.Textarea, cmd = m.Textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m JournalModel) View() string {
	return fmt.Sprintf(
		"Tell me a story.\n\n%s\n\n%s",
		m.Textarea.View(),
		"(ctrl+c to quit)",
	) + "\n\n"
}
