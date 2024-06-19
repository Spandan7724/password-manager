package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textInputModel struct {
	textInputs   []textinput.Model
	cursor       int
	choices      []string
	err          error
	done         bool
	maskPassword bool // Field to toggle masking
}

func newTextInputModel() textInputModel {
	ti := textinput.NewModel()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	// Create a list of text inputs for website, username, and password
	textInputs := []textinput.Model{ti, textinput.NewModel(), textinput.NewModel()}
	textInputs[2].EchoMode = textinput.EchoPassword // Set echo mode to password for the password field
	textInputs[2].EchoCharacter = '*'               // Set the character to use for masking

	return textInputModel{
		textInputs:   textInputs,
		choices:      []string{"Website", "Username", "Password (Ctrl+T to toggle mask)"},
		maskPassword: true, // Start with password masking enabled
	}
}

func (m textInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.done = true
			return m, tea.Quit
		case "enter":
			if m.cursor < len(m.textInputs)-1 {
				m.cursor++
				m.textInputs[m.cursor].Focus()
			} else {
				m.done = true
				return m, tea.Quit
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.textInputs[m.cursor].Focus()
			}
		case "down":
			if m.cursor < len(m.textInputs)-1 {
				m.cursor++
				m.textInputs[m.cursor].Focus()
			}
		case "ctrl+t":
			// Toggle password masking
			m.maskPassword = !m.maskPassword
			if m.maskPassword {
				m.textInputs[2].EchoMode = textinput.EchoPassword
			} else {
				m.textInputs[2].EchoMode = textinput.EchoNormal
			}
		}
	}

	var cmd tea.Cmd
	m.textInputs[m.cursor], cmd = m.textInputs[m.cursor].Update(msg)
	return m, cmd
}

func (m textInputModel) View() string {
	if m.done {
		return "Finished"
	}

	var b strings.Builder

	for i := range m.textInputs {
		fmt.Fprintf(&b, "%s\n%s\n\n", m.choices[i], m.textInputs[i].View())
	}

	return b.String()
}
