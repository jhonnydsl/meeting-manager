package services

import (
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils"
)

type MeetingService struct {
	MeetingRepo *repository.MeetingRepository
}

func (service *MeetingService) CreateMeeting(meeting dtos.Meeting, ownerID int) (dtos.MeetingOutput, error) {
	if err := utils.ValidateMeetingInput(meeting); err != nil {
		return dtos.MeetingOutput{}, err
	}

	hasConflict, err := service.MeetingRepo.HasConflict(ownerID, meeting.StartTime, meeting.EndTime)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	if hasConflict {
		return dtos.MeetingOutput{}, utils.ConflictError("schedule conflict: there is already a meeting in this time slot")
	}

	return service.MeetingRepo.CreateMeeting(meeting, ownerID)
}

func (service *MeetingService) GetAllMeetings(ownerID int) ([]dtos.MeetingOutput, error) {
	return service.MeetingRepo.GetAllMeetings(ownerID)
}

func (service *MeetingService) UpdateMeeting(meeting dtos.MeetingOutput, ownerID int) (dtos.MeetingOutput, error) {
	if err := utils.ValidateMeetingOutput(meeting); err != nil {
		return dtos.MeetingOutput{}, err
	}

	hasConflict, err := service.MeetingRepo.HasConflict(ownerID, meeting.StartTime, meeting.EndTime, meeting.ID)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	if hasConflict {
		return dtos.MeetingOutput{}, utils.ConflictError("schedule conflict: there is already a meeting during this time range")
	}

	return service.MeetingRepo.UpdateMeeting(meeting, ownerID)
}

func (service *MeetingService) DeleteMeeting(id, ownerID int) error {
	return service.MeetingRepo.DeleteMeeting(id, ownerID)
}