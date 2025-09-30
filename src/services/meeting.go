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
		return dtos.MeetingOutput{}, errors.New("error creating meeting, please fill in the required fields")
	}

	if meeting.StartTime.After(meeting.EndTime) {
		return dtos.MeetingOutput{}, errors.New("start_time cannot be after end_time")
	}

	hasConflict, err := service.MeetingRepo.HasConflict(ownerID, meeting.StartTime, meeting.EndTime)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	if hasConflict {
		return dtos.MeetingOutput{}, errors.New("schedule conflict: there is already a meeting in this time slot")
	}

	return service.MeetingRepo.CreateMeeting(meeting, ownerID)
}

func (service *MeetingService) GetAllMeetings(ownerID int) ([]dtos.MeetingOutput, error) {
	return service.MeetingRepo.GetAllMeetings(ownerID)
}

func (service *MeetingService) UpdateMeeting(meeting dtos.MeetingOutput, ownerID int) (dtos.MeetingOutput, error) {
	if meeting.Title == "" || meeting.StartTime.IsZero() || meeting.EndTime.IsZero() {
		return dtos.MeetingOutput{}, errors.New("error updating meeting, please fill in the required fields")
	}

	if meeting.StartTime.After(meeting.EndTime) {
		return dtos.MeetingOutput{}, errors.New("start_time cannot be after end_time")
	}

	hasConflict, err := service.MeetingRepo.HasConflict(ownerID, meeting.StartTime, meeting.EndTime, meeting.ID)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	if hasConflict {
		return dtos.MeetingOutput{}, errors.New("schedule conflict: there is already a meeting during this time range")
	}

	return service.MeetingRepo.UpdateMeeting(meeting, ownerID)
}

func (service *MeetingService) DeleteMeeting(id, ownerID int) error {
	return service.MeetingRepo.DeleteMeeting(id, ownerID)
}