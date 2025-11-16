package services

import (
	"time"

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

func (service *MeetingService) UpdateMeeting(input dtos.UpdateMeeting, meetingID, ownerID int) (dtos.MeetingOutput, error) {
	if meetingID == 0 {
		return dtos.MeetingOutput{}, utils.BadRequestError("invalid meeting id")
	}

	layout := "02/01/2006 15:04"
	start, err := time.Parse(layout, input.StartTime)
	if err != nil {
		return dtos.MeetingOutput{}, utils.BadRequestError("invalid start_time")
	}
	end, err := time.Parse(layout, input.EndTime)
	if err != nil {
		return dtos.MeetingOutput{}, utils.BadRequestError("invalid end_time")
	}

	hasConflict, err := service.MeetingRepo.HasConflict(ownerID, start, end, meetingID)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	if hasConflict {
		return dtos.MeetingOutput{}, utils.ConflictError("schedule conflict: there is already a meeting during this time range")
	}

	/*meeting := dtos.MeetingOutput{
		ID: meetingID,
		Title: input.Title,
		Description: input.Description,gi
		StartTime: start,
		EndTime: end,
		OwnerID: ownerID,
	}*/

	return service.MeetingRepo.UpdateMeeting(input, meetingID, ownerID, start, end)
}

func (service *MeetingService) DeleteMeeting(id, ownerID int) error {
	return service.MeetingRepo.DeleteMeeting(id, ownerID)
}