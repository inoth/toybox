package base

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{"yes", "no"}

type TextModel struct {
	textInput textinput.Model
	err       error

	Message string
	Value   string
}

func NewTextModel(message string) TextModel {
	ti := textinput.New()
	ti.Placeholder = "helloworld"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return TextModel{
		textInput: ti,
		Message:   message,
		err:       nil,
	}
}

func (m *TextModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *TextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.Value = m.textInput.Value()
	return m, cmd
}

func (m *TextModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n",
		m.Message,
		m.textInput.View(),
		"(esc to quit)",
	)
}

type SelectModel struct {
	cursor int
	Choice string

	Message string
	Value   string
}

func (m *SelectModel) Init() tea.Cmd { return nil }

func (m *SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.Choice = choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m *SelectModel) View() string {
	s := strings.Builder{}
	s.WriteString("📂 Do you want to override the folder ?\n\n")
	s.WriteString("Delete the existing folder and create the project.\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
