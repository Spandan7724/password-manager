package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var initialStyle = lipgloss.NewStyle().
	Padding(1).
	Border(lipgloss.RoundedBorder()).
	Background(lipgloss.Color("#3a00fa")).
	Foreground(lipgloss.Color("#FFFFFF")).Bold(true).
	BorderForeground(lipgloss.Color("#FFFFFF")).Bold(true)

var menuTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ff006e")).Bold(true).
	Padding(0, 1)

var menuBorderStyle = lipgloss.NewStyle().
	Padding(1).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#3a00fa"))

var promptStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ff006e"))

type initialModel struct{}

func (m initialModel) Init() tea.Cmd {
	return nil
}

func (m initialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			return mainModel{}, tea.Quit
		}
	}
	return m, nil
}

func (m initialModel) View() string {
	return initialStyle.Render("Password Manager\n\nPress Enter to continue")
}

type mainModel struct{}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m mainModel) View() string {
	return ""
}

// Custom text input function
func prompt(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptStyle.Render(message))
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Main menu
func mainMenu() {
	for {
		menuItems := []string{
			"1. View saved passwords",
			"2. Add a new password",
			"3. Update an existing password",
			"4. Delete a password",
			"5. Exit",
		}

		var menuContent strings.Builder
		for _, item := range menuItems {
			menuContent.WriteString(menuTextStyle.Render(item) + "\n\n")
		}

		menu := menuBorderStyle.Render(menuContent.String())

		fmt.Println(initialStyle.Render(" Password Manager "))
		fmt.Println(menu)

		choice := prompt("Choose an option: ")

		switch choice {
		case "1":
			viewPasswordEntries()
		case "2":
			inputs := newTextInputModel()
			p := tea.NewProgram(inputs)
			if err := p.Start(); err != nil {
				fmt.Printf("Error running program: %v", err)
				os.Exit(1)
			}

			website := inputs.textInputs[0].Value()
			username := inputs.textInputs[1].Value()
			password := inputs.textInputs[2].Value()
			addPasswordEntry(website, username, password)

		case "3":
			website := prompt("Enter website: ")
			username := prompt("Enter username: ")
			newPassword := prompt("Enter new password: ")
			updatePasswordEntry(website, username, newPassword)
		case "4":
			website := prompt("Enter website: ")
			username := prompt("Enter username: ")
			deletePasswordEntry(website, username)
		case "5":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func Run() {
	p := tea.NewProgram(initialModel{})
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting program: %v", err)
		os.Exit(1)
	}

	for {
		username := prompt("Enter username: ")

		passwordModel := newPasswordInputModel()
		passwordProgram := tea.NewProgram(passwordModel)
		if err := passwordProgram.Start(); err != nil {
			fmt.Printf("Error running password input: %v", err)
			os.Exit(1)
		}
		password := passwordModel.Value()

		if userExists(username) {
			if authenticate(username, password) {
				currentUser = username
				break
			} else {
				fmt.Println("Invalid credentials, please try again.")
			}
		} else {
			fmt.Println("User not found. Creating new user.")
			if err := createUser(username, password); err != nil {
				fmt.Println("Error creating user:", err)
				continue
			}
			fmt.Println("User created successfully. Please log in.")
		}
	}

	mainMenu()
}

// New password input model for login
type passwordInputModel struct {
	textinput.Model
	done         bool
	maskPassword bool
}

func newPasswordInputModel() passwordInputModel {
	ti := textinput.NewModel()
	ti.Placeholder = "Enter master password (Ctrl+T to toggle mask): "
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.EchoMode = textinput.EchoPassword // Set echo mode to password for the password field
	ti.EchoCharacter = '*'               // Set the character to use for masking
	return passwordInputModel{Model: ti, maskPassword: true}
}

func (m passwordInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m passwordInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.done = true
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+t":
			// Toggle password masking
			m.maskPassword = !m.maskPassword
			if m.maskPassword {
				m.Model.EchoMode = textinput.EchoPassword
			} else {
				m.Model.EchoMode = textinput.EchoNormal
			}
		}
	}

	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)
	return m, cmd
}

func (m passwordInputModel) View() string {
	if m.done {
		return "Finished"
	}
	return m.Model.View()
}

func (m passwordInputModel) Value() string {
	return m.Model.Value()
}
