package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "passwords.db"

type PasswordEntry struct {
	Website  string
	Username string
	Password string
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, hash TEXT);
	CREATE TABLE IF NOT EXISTS passwords (username TEXT, website TEXT, site_username TEXT, password TEXT, PRIMARY KEY (username, website, site_username));
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
}

// saveUser saves a new user in the database
func saveUser(username, hash string) error {
	_, err := db.Exec("INSERT INTO users (username, hash) VALUES (?, ?)", username, hash)
	return err
}

// getUserHash retrieves the hash for a given username
func getUserHash(username string) (string, error) {
	var hash string
	err := db.QueryRow("SELECT hash FROM users WHERE username = ?", username).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// savePassword saves a new password entry in the database
func savePassword(user, website, siteUser, password string) error {
	_, err := db.Exec("INSERT INTO passwords (username, website, site_username, password) VALUES (?, ?, ?, ?)", user, website, siteUser, password)
	return err
}

// getPasswords retrieves all password entries for a given user
func getPasswords(user string) ([]PasswordEntry, error) {
	rows, err := db.Query("SELECT website, site_username, password FROM passwords WHERE username = ?", user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []PasswordEntry
	for rows.Next() {
		var entry PasswordEntry
		if err := rows.Scan(&entry.Website, &entry.Username, &entry.Password); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// updatePassword updates an existing password entry in the database
func updatePassword(user, website, siteUser, newPassword string) error {
	_, err := db.Exec("UPDATE passwords SET password = ? WHERE username = ? AND website = ? AND site_username = ?", newPassword, user, website, siteUser)
	return err
}

// deletePassword deletes a password entry from the database
func deletePassword(user, website, siteUser string) error {
	_, err := db.Exec("DELETE FROM passwords WHERE username = ? AND website = ? AND site_username = ?", user, website, siteUser)
	return err
}
