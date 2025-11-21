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

func ValidateMeetingInput(meeting dtos.Meeting) error {
	if meeting.Title == "" {
		return errors.New("title must not empty")
	}

	if meeting.Status != "iniciada" && meeting.Status != "finalizada" && meeting.Status != "cancelada" && meeting.Status != "agendada" {
		return errors.New("status invalid")
	}

	if meeting.StartTime.IsZero() || meeting.EndTime.IsZero() {
		return errors.New("start_time and end_time is required")
	}

	if meeting.StartTime.After(meeting.EndTime) {
		return errors.New("start_time cannot be after end_time")
	}

	return nil
}