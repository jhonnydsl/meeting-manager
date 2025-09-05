package repository

import (
	"errors"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
)

type InvitationRepository struct{}

// CreateInvitation inserts a new invitation into the database and returns it.
func (r *InvitationRepository) CreateInvitation(invitation dtos.InvitationInput, senderID int) (dtos.InvitationOutput, error) {
	query := `INSERT INTO convites (reuniao_id, receiver_id, sender_id) VALUES ($1, $2, $3) RETURNING id, reuniao_id, receiver_id, status, created_at, sender_id`
	var createdInvitation dtos.InvitationOutput

	// Execute insert and scan result into createdInvitation.
	err := DB.QueryRow(query, invitation.ReuniaoID, invitation.ReceiverID, senderID).Scan(
		&createdInvitation.ID,
		&createdInvitation.ReuniaoID,
		&createdInvitation.ReceiverID,
		&createdInvitation.Status,
		&createdInvitation.CreatedAt,
		&createdInvitation.SenderID,
	)
	if err != nil {
		return dtos.InvitationOutput{}, err
	}

	return createdInvitation, nil
}

// GetAllInvitations retrieves invitations sent by a specific user.
func (r *InvitationRepository) GetAllInvitations(senderID int) ([]dtos.InvitationOutput, error) {
	query := `SELECT id, reuniao_id, receiver_id, status, created_at, sender_id FROM convites WHERE sender_id = $1`
	var lista []dtos.InvitationOutput

	rows, err := DB.Query(query, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through rows and build the list of invitations.
	for rows.Next() {
		var i dtos.InvitationOutput

		err = rows.Scan(&i.ID, &i.ReuniaoID, &i.ReceiverID, &i.Status, &i.CreatedAt, &i.SenderID)
		if err != nil {
			return nil, err
		}

		lista = append(lista, i)
	}

	return lista, nil
}

// GetReceiver retrieves invitations received by a specific user.
func (r *InvitationRepository) GetReceiver(receiverID int) ([]dtos.InvitationOutput, error) {
	query := `SELECT id, reuniao_id, receiver_id, status, created_at, sender_id FROM convites WHERE receiver_id = $1`
	var lista []dtos.InvitationOutput

	rows, err := DB.Query(query, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i dtos.InvitationOutput

		err = rows.Scan(&i.ID, &i.ReuniaoID, &i.ReceiverID, &i.Status, &i.CreatedAt, &i.SenderID)
		if err != nil {
			return nil, err
		}

		lista = append(lista, i)
	}

	return lista, nil
}

// GetOwnerID retrieves the owner ID of a given meeting.
func (r *InvitationRepository) GetOwnerID(meetingID int) (int, error) {
	query := `SELECT owner_id FROM reunioes WHERE id = $1`
	var ownerID int

	err := DB.QueryRow(query, meetingID).Scan(&ownerID)
	if err != nil {
		return 0, err
	}

	return ownerID, nil
}

// DeleteInvitation removes an invitation if the user is the owner of the meeting.
func (r *InvitationRepository) DeleteInvitation(invitationID int, ownerID int ) error {
	query := `DELETE FROM convites WHERE id = $1 AND reuniao_id IN (SELECT id FROM reunioes WHERE owner_id = $2)`

	result, err := DB.Exec(query, invitationID, ownerID)
	if err != nil {
		return err
	}

	// Ensure at least one row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("nenhum convite deletado: ou não existe ou você não é o dono da reunião")
	}

	return nil
}

// UpdateInvitationStatus updates the status of an invitation (e.g., sent/pending).
func (r *InvitationRepository) UpdateInvitationStatus(invitationID int, status string) error {
	query := `UPDATE convites SET status = $1 WHERE id = $2`
	
	_, err := DB.Exec(query, status, invitationID)
	return err
}

// ReturnUserByEmail retrieves the email of a user by ID.
func (r *InvitationRepository) ReturnUserByEmail(userID int) (string, error) {
	var email string
	query := `SELECT email FROM users WHERE id = $1`

	err := DB.QueryRow(query, userID).Scan(&email)
	if err != nil {
		return "", errors.New("usuário não encontrado")
	}

	return email, nil
}