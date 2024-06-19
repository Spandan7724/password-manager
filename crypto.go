package main

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"crypto/sha256"
"encoding/base64"
"errors"
"io"
)

// deriveKey derives a 32-byte key from the user's master password
func deriveKey(password string) []byte {
hash := sha256.Sum256([]byte(password))
return hash[:]
}

// encrypt encrypts the plaintext using the provided key
func encrypt(plaintext string, key []byte) (string, error) {
block, err := aes.NewCipher(key)
if err != nil {
return "", err
}

gcm, err := cipher.NewGCM(block)
if err != nil {
	return "", err
}

nonce := make([]byte, gcm.NonceSize())
if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	return "", err
}

ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts the ciphertext using the provided key
func decrypt(ciphertext string, key []byte) (string, error) {
data, err := base64.StdEncoding.DecodeString(ciphertext)
if err != nil {
return "", err
}
block, err := aes.NewCipher(key)
if err != nil {
	return "", err
}

gcm, err := cipher.NewGCM(block)
if err != nil {
	return "", err
}

nonceSize := gcm.NonceSize()
if len(data) < nonceSize {
	return "", errors.New("ciphertext too short")
}

nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
if err != nil {
	return "", err
}

return string(plaintext), nil
}

