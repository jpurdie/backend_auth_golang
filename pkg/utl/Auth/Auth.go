package auth

import (
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

