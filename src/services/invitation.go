package services

import (
	"errors"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
)

type InvitationService struct {
	InvitRepo *repository.InvitationRepository
}

// CreateInvitation checks if IDs are valid and ensures that only
// the meeting owner can send invitations.
func (service *InvitationService) CreateInvitation(invitation dtos.InvitationInput, senderID int) (dtos.InvitationOutput, error) {
	if invitation.ReuniaoID <= 0 || invitation.ReceiverID <= 0 {
		return dtos.InvitationOutput{}, errors.New("reuniaoID e receiverID devem ser maiores que zero")
	}

	// Fetch the meeting owner ID from the database.
	ownerID, err := service.InvitRepo.GetOwnerID(invitation.ReuniaoID)
	if err != nil {
		return dtos.InvitationOutput{}, err
	}

	// Validate that the sender is the meeting owner.
	if ownerID != senderID {
		return dtos.InvitationOutput{}, errors.New("somente o dono da reunião pode enviar convites")
	}

	// Delegate creation to repository
	return service.InvitRepo.CreateInvitation(invitation, senderID)
}

// GetAllInvitations returns all invitations created by a specific sender.
func (service *InvitationService) GetAllInvitations(senderID int) ([]dtos.InvitationOutput, error) {
	return service.InvitRepo.GetAllInvitations(senderID)
}

// GetReceiver returns all invitations received by a specific user.
func (service *InvitationService) GetReceiver(receiverID int) ([]dtos.InvitationOutput, error) {
	return service.InvitRepo.GetReceiver(receiverID)
}

// DeleteInvitation removes an invitation if the requester is the meeting owner.
func (service *InvitationService) DeleteInvitation(invitationID int, ownerID int) error {
	return service.InvitRepo.DeleteInvitation(invitationID, ownerID)
}

// UpdateInvitationStatus updates the status of an invitation (e.g., accepted or declined).
func (service *InvitationService) UpdateInvitationStatus(invitationID int, status string) error {
	return service.InvitRepo.UpdateInvitationStatus(invitationID, status)
}

// ReturnUserByEmail retrieves the email of a user based on their ID.
func (service *InvitationService) ReturnUserByEmail(userID int) (string, error) {
	return service.InvitRepo.ReturnUserByEmail(userID)
}