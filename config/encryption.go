package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)


func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := 0; i < padding; i++ {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// PKCS7 unpadding function
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	padding := data[len(data)-1]
	if padding > byte(blockSize) {
		return nil, fmt.Errorf("invalid padding size")
	}
	return data[:len(data)-int(padding)], nil
}

func EncryptID(id string) (string, error) {
	key := []byte(os.Getenv("AES_KEY")) 
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Apply PKCS7 padding
	idBytes := []byte(id)
	idBytes = pkcs7Pad(idBytes, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(idBytes))
	iv := ciphertext[:aes.BlockSize]

	// Generate a random IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt the ID
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], idBytes)

	// Return base64-encoded string
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptID(encryptedID string) (string, error) {
	key := []byte(os.Getenv("AES_KEY")) 
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedID)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt the ID
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	// Apply PKCS7 unpadding
	decryptedID, err := pkcs7Unpad(ciphertext, aes.BlockSize)
	if err != nil {
		return "", err
	}

	// Return the decrypted ID as a string
	return string(decryptedID), nil
}
