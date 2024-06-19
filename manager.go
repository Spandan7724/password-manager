package main

import (
"encoding/base64"
"fmt"
"log"
"math/rand"

"github.com/charmbracelet/bubbles/list"
)

var currentUser string
var masterKey []byte

func authenticate(username, password string) bool {
storedHash, err := getUserHash(username)
if err != nil {
return false
}
if checkPasswordHash(password, storedHash) {
masterKey = deriveKey(password)
return true
}
return false
}

func createUser(username, password string) error {
hash, err := hashPassword(password)
if err != nil {
return err
}
return saveUser(username, hash)
}

func generateSalt() (string, error) {
salt := make([]byte, 16)
_, err := rand.Read(salt)
if err != nil {
return "", err
}
return base64.StdEncoding.EncodeToString(salt), nil
}

func addPasswordEntry(website, username, password string) {
salt, err := generateSalt()
if err != nil {
log.Println("Error generating salt:", err)
return
}

saltedPassword := salt + password
encryptedPassword, err := encrypt(saltedPassword, masterKey)
if err != nil {
	log.Println("Error encrypting password:", err)
	return
}

err = savePassword(currentUser, website, username, encryptedPassword)
if err != nil {
	log.Println("Error saving password:", err)
} else {
	fmt.Println("Password saved successfully!")
}
}

func viewPasswordEntries() {
entries, err := getPasswords(currentUser)
if err != nil {
log.Println("Error retrieving passwords:", err)
return
}

items := make([]list.Item, len(entries))
for i, entry := range entries {
	decryptedPassword, err := decrypt(entry.Password, masterKey)
	if err != nil {
		log.Println("Error decrypting password:", err)
		continue
	}
	saltSize := 24 // length of base64 encoded 16 byte salt
	saltedPassword := decryptedPassword[saltSize:]
	items[i] = item{
		title:       entry.Website,
		description: fmt.Sprintf("Username: %s, Password: %s", entry.Username, saltedPassword),
	}
}

runBubbleTea(items)
}

func updatePasswordEntry(website, username, newPassword string) {
salt, err := generateSalt()
if err != nil {
log.Println("Error generating salt:", err)
return
}
saltedPassword := salt + newPassword
encryptedPassword, err := encrypt(saltedPassword, masterKey)
if err != nil {
	log.Println("Error encrypting password:", err)
	return
}

err = updatePassword(currentUser, website, username, encryptedPassword)
if err != nil {
	log.Println("Error updating password:", err)
} else {
	fmt.Println("Password updated successfully!")
}
}

func deletePasswordEntry(website, username string) {
err := deletePassword(currentUser, website, username)
if err != nil {
log.Println("Error deleting password:", err)
} else {
fmt.Println("Password deleted successfully!")
}
}

func userExists(username string) bool {
_, err := getUserHash(username)
return err == nil
}

