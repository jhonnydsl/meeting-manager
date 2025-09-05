package services

import (
	"errors"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
)

type MeetingService struct {
	MeetingRepo *repository.MeetingRepository
}

// CreateMeeting validates business rules and delegates creation to the repository.
func (service *MeetingService) CreateMeeting(meeting dtos.Meeting, ownerID int) (dtos.MeetingOutput, error) {
	// Validate required fields.
	if meeting.Title == "" || meeting.StartTime.IsZero() || meeting.EndTime.IsZero() {
		return dtos.MeetingOutput{}, errors.New("erro ao criar reunião, favor preencher campos obrigatorios")
	}

	// Ensure start_time is before end_time.
	if meeting.StartTime.After(meeting.EndTime) {
		return dtos.MeetingOutput{}, errors.New("start_time não pode ser depois do end_time")
	}

	// Check for scheduling conflicts.
	hasConflict, err := service.MeetingRepo.HasConflict(ownerID, meeting.StartTime, meeting.EndTime)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	if hasConflict {
		return dtos.MeetingOutput{}, errors.New("conflito de horario: já existe uma reunião nesse intervalo")
	}

	// Delegate to repository to save the meeting.
	return service.MeetingRepo.CreateMeeting(meeting, ownerID)
}

// GetAllMeetings retrieves all meetings for a given owner.
func (service *MeetingService) GetAllMeetings(ownerID int) ([]dtos.MeetingOutput, error) {
	return service.MeetingRepo.GetAllMeetings(ownerID)
}

// UpdateMeeting validates data, checks conflicts, and updates an existing meeting.
func (service *MeetingService) UpdateMeeting(meeting dtos.MeetingOutput, ownerID int) (dtos.MeetingOutput, error) {
	// Validate required fields.
	if meeting.Title == "" || meeting.StartTime.IsZero() || meeting.EndTime.IsZero() {
		return dtos.MeetingOutput{}, errors.New("erro ao atualizar reunião, favor preencher campos obrigatorios")
	}

	// Ensure start_time is before end_time.
	if meeting.StartTime.After(meeting.EndTime) {
		return dtos.MeetingOutput{}, errors.New("start_time não pode ser depois do end_time")
	}

	// Check for conflicts excluding the current meeting ID.
	hasConflict, err := service.MeetingRepo.HasConflict(ownerID, meeting.StartTime, meeting.EndTime, meeting.ID)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	if hasConflict {
		return dtos.MeetingOutput{}, errors.New("conflito de horario: já existe uma reunião nesse intervalo")
	}

	// Delegate update to repository
	return service.MeetingRepo.UpdateMeeting(meeting, ownerID)
}

// DeleteMeeting delegates deletion to the repository.
func (service *MeetingService) DeleteMeeting(id, ownerID int) error {
	return service.MeetingRepo.DeleteMeeting(id, ownerID)
}