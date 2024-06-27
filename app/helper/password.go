package helper

import (
	"crypto/rand"
	"math/big"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	passwordLength = 16
	letters        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits         = "0123456789"
	specialChars   = "!@#$%^&*()-_=+[]{}|;:,.<>?/`~"
)

// generates a strong password of the specified length
func GenerateStrongPassword() (string, error) {
	var password strings.Builder

	// Ensure password contains at least one character from each set
	charsets := []string{letters, digits, specialChars}

	for _, charset := range charsets {
		char, err := generateRandomChar(charset)
		if err != nil {
			return "", err
		}
		password.WriteByte(char)
	}

	// Fill the rest of the password length with random characters from all sets
	allChars := letters + digits + specialChars
	for password.Len() < passwordLength {
		char, err := generateRandomChar(allChars)
		if err != nil {
			return "", err
		}
		password.WriteByte(char)
	}

	// Convert the password to a slice to shuffle the characters
	passwordSlice := []byte(password.String())
	for i := range passwordSlice {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordSlice))))
		if err != nil {
			return "", err
		}
		passwordSlice[i], passwordSlice[j.Int64()] = passwordSlice[j.Int64()], passwordSlice[i]
	}

	return string(passwordSlice), nil
}

func generateRandomChar(charset string) (byte, error) {
	max := big.NewInt(int64(len(charset)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return charset[n.Int64()], nil
}

// calculates the strength of the password and returns a score from 0 to 100
func CheckPasswordStrength(password string) int {
	var score int

	// Criteria
	length := len(password)
	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	// Length score
	if length >= 8 {
		score += 25
	} else {
		score += length * 3 // up to 21
	}

	// Uppercase letter score
	if hasUpper {
		score += 15
	}

	// Lowercase letter score
	if hasLower {
		score += 15
	}

	// Digit score
	if hasDigit {
		score += 20
	}

	// Special character score
	if hasSpecial {
		score += 25
	}

	// Ensure the score is capped at 100
	if score > 100 {
		score = 100
	}

	return score
}

// encrypt entered password
func EncryptPassword(password string) (string, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// verifies the encrypted password
func ValidateEncryptedPassword(hashedPassword, enteredPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(enteredPassword))
	return err == nil
}
