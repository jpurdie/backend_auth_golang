package Auth

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"unicode"
)

func VerifyPassword(s string) bool {
	var (
		hasMinLen = false
		hasMaxLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
		//hasSpecial = false
	)
	if len(s) >= 10 {
		hasMinLen = true
	}
	if len(s) <= 64 {
		hasMaxLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
			//case unicode.IsPunct(char) || unicode.IsSymbol(char):
			//hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasMaxLen
}

func GenerateInviteTokenHash(str string) string {
	myHash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(myHash[:])
}

func GenerateInviteToken() string {
	n := 64
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
