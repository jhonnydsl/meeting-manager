package services

import (
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils"
)

type InvitationService struct {
	InvitRepo *repository.InvitationRepository
}

func (service *InvitationService) CreateInvitation(invitation dtos.InvitationInput, senderID int) (dtos.InvitationOutput, error) {
	if invitation.ReuniaoID <= 0 || invitation.ReceiverID <= 0 {
		return dtos.InvitationOutput{}, utils.BadRequestError("reuniaoID e receiverID must be greater than zero")
	}

	ownerID, err := service.InvitRepo.GetOwnerID(invitation.ReuniaoID)
	if err != nil {
		return dtos.InvitationOutput{}, err
	}

	if ownerID != senderID {
		return dtos.InvitationOutput{}, utils.BadRequestError("only the meeting owner can send invitations")
	}

	return service.InvitRepo.CreateInvitation(invitation, senderID)
}

func (service *InvitationService) GetAllInvitations(senderID int) ([]dtos.InvitationOutput, error) {
	return service.InvitRepo.GetAllInvitations(senderID)
}

func (service *InvitationService) GetReceiver(receiverID int) ([]dtos.InvitationOutput, error) {
	return service.InvitRepo.GetReceiver(receiverID)
}

func (service *InvitationService) DeleteInvitation(invitationID int, ownerID int) error {
	return service.InvitRepo.DeleteInvitation(invitationID, ownerID)
}

func (service *InvitationService) UpdateInvitationStatus(invitationID int, status string) error {
	return service.InvitRepo.UpdateInvitationStatus(invitationID, status)
}

func (service *InvitationService) ReturnUserByEmail(userID int) (string, error) {
	return service.InvitRepo.ReturnUserByEmail(userID)
}