# password-manager

A terminal-based password manager application written in Go for managing passwords securely. This application uses encryption to securely store your passwords and provides a simple interface for managing them.

## Features

- Secure storage of passwords using AES encryption
- User authentication with hashed master passwords
- Interactive terminal UI with masking for password inputs
- Easy-to-use menu for viewing, adding, updating, and deleting passwords
- Toggle password masking on/off for password input fields

## Getting Started

### Prerequisites

- Go (version 1.16+)
- SQLite3

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/Spandan7724/password-manager.git
   ```
2. Change to the project directory   
   ```sh
   cd password-manager
   ```
3. Install the required dependencies:
   ```sh
   go get -u github.com/charmbracelet/bubbletea
   go get -u github.com/charmbracelet/bubbles
   go get -u github.com/charmbracelet/lipgloss
   go get -u github.com/eiannone/keyboard
   go get -u github.com/mattn/go-sqlite3
   go get -u github.com/pterm/pterm
   ```   
### Running the Application
    
1. Run Directly         
    ```sh
    go run .
    ```
2. Or, Run the executable from the project directory 

## Acknowledgments

*   [Bubble Tea](https://github.com/charmbracelet/bubbletea)
*   [Lip Gloss](https://github.com/charmbracelet/lipgloss)
*   [Pterm](https://github.com/pterm/pterm)
*   [Keyboard](https://github.com/eiannone/keyboard)
*   [Go-SQLite3](https://github.com/mattn/go-sqlite3)

