package utils

import (
	"errors"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
)

func ValidateMeetingOutput(meeting dtos.MeetingOutput) error {
	if meeting.Title == "" {
		return errors.New("title must not empty")
	}

	if meeting.StartTime.IsZero() || meeting.EndTime.IsZero() {
		return errors.New("start_time and end_time is required")
	}

	if meeting.StartTime.After(meeting.EndTime) {
		return errors.New("start_time cannot be after end_time")
	}

	return nil
}