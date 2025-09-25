package utils

import (
	"errors"
	"unicode/utf8"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
)

func ValidateUserInput(user dtos.UserInput) error {
	if user.Email == "" || user.Name == "" {
		return errors.New("name and email must not be empty")
	}

	if utf8.RuneCountInString(user.Password) < 6 {
		return errors.New("password must contain at least 6 characters")
	}

	return nil
}