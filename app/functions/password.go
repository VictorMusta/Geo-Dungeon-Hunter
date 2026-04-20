package functions

import (
	"crypto/rand"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt ...
func HashAndSalt(pwd string) ([]byte, error) {
	// Increased cost for better protection against brute-force
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return hash, err
}

/* 
	Called in the auth controller.
	Is used to check if the password is correct with the user's password in the database.
*/
func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GeneratePassword ...
func GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var (
		lowerCharSet   = "abcdefghijklmnopqrstuvwxyz"
		upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialCharSet = "!@#$%&*"
		numberSet      = "0123456789"
		allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
	)

	var password strings.Builder

	// Secure random index helper
	secureRandInt := func(max int) int {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
		return int(nBig.Int64())
	}

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		password.WriteString(string(specialCharSet[secureRandInt(len(specialCharSet))]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		password.WriteString(string(numberSet[secureRandInt(len(numberSet))]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		password.WriteString(string(upperCharSet[secureRandInt(len(upperCharSet))]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		password.WriteString(string(allCharSet[secureRandInt(len(allCharSet))]))
	}

	// Shuffle using Fisher-Yates and crypto/rand
	inRune := []rune(password.String())
	for i := len(inRune) - 1; i > 0; i-- {
		j := secureRandInt(i + 1)
		inRune[i], inRune[j] = inRune[j], inRune[i]
	}
	
	return string(inRune)
}
