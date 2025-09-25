package services

import (
	"errors"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
)

type MeetingService struct {
	MeetingRepo *repository.MeetingRepository
}

func (service *MeetingService) CreateMeeting(meeting dtos.Meeting, ownerID int) (dtos.MeetingOutput, error) {
	if meeting.Title == "" || meeting.StartTime.IsZero() || meeting.EndTime.IsZero() {
		return dtos.MeetingOutput{}, errors.New("erro ao criar reunião, favor preencher campos obrigatorios")
	}

	if meeting.StartTime.After(meeting.EndTime) {
		return dtos.MeetingOutput{}, errors.New("start_time não pode ser depois do end_time")
	}

	hasConflict, err := service.MeetingRepo.HasConflict(ownerID, meeting.StartTime, meeting.EndTime)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	if hasConflict {
		return dtos.MeetingOutput{}, errors.New("conflito de horario: já existe uma reunião nesse intervalo")
	}

	return service.MeetingRepo.CreateMeeting(meeting, ownerID)
}

func (service *MeetingService) GetAllMeetings(ownerID int) ([]dtos.MeetingOutput, error) {
	return service.MeetingRepo.GetAllMeetings(ownerID)
}

func (service *MeetingService) UpdateMeeting(meeting dtos.MeetingOutput, ownerID int) (dtos.MeetingOutput, error) {
	if meeting.Title == "" || meeting.StartTime.IsZero() || meeting.EndTime.IsZero() {
		return dtos.MeetingOutput{}, errors.New("erro ao atualizar reunião, favor preencher campos obrigatorios")
	}

	if meeting.StartTime.After(meeting.EndTime) {
		return dtos.MeetingOutput{}, errors.New("start_time não pode ser depois do end_time")
	}

	hasConflict, err := service.MeetingRepo.HasConflict(ownerID, meeting.StartTime, meeting.EndTime, meeting.ID)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	if hasConflict {
		return dtos.MeetingOutput{}, errors.New("conflito de horario: já existe uma reunião nesse intervalo")
	}

	return service.MeetingRepo.UpdateMeeting(meeting, ownerID)
}

func (service *MeetingService) DeleteMeeting(id, ownerID int) error {
	return service.MeetingRepo.DeleteMeeting(id, ownerID)
}