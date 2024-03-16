package utils

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 2)
	return string(bytes), err
}

func CheckPassword(password string) error {
    if utf8.RuneCountInString(password) < 6 {
        return errors.New("пароль должен содержать не менее 6 символов")
    }

    containsLower := false
    containsUpper := false
    containsDigit := false

    for _, char := range password {
        switch {
        case unicode.IsLower(char):
            containsLower = true
        case unicode.IsUpper(char):
            containsUpper = true
        case unicode.IsDigit(char):
            containsDigit = true
        }
    }

    if !containsLower || !containsUpper || !containsDigit {
        return errors.New("пароль должен содержать символы в верхнем и нижнем регистре, а также хотя бы одну цифру")
    }

    return nil
}