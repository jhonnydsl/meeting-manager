package repository

import (
	"errors"
	"time"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
)

type MeetingRepository struct{}

// CreateMeeting inserts a new meeting into the database and returns the created meeting.
func (r *MeetingRepository) CreateMeeting(meeting dtos.Meeting, ownerId int) (dtos.MeetingOutput, error) {
	query := `INSERT INTO reunioes (title, description, start_time, end_time, owner_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, description, start_time, end_time, owner_id, created_at`
	var createdMeeting dtos.MeetingOutput

	// Execute the query and scan the returned values into createdMeeting.
	err := DB.QueryRow(query, meeting.Title, meeting.Description, meeting.StartTime, meeting.EndTime, ownerId).Scan(
		&createdMeeting.ID,
		&createdMeeting.Title,
		&createdMeeting.Description,
		&createdMeeting.StartTime,
		&createdMeeting.EndTime,
		&createdMeeting.OwnerID,
		&createdMeeting.CreatedAt,
	)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	return createdMeeting, nil
}

// GetAllMeetings retrieves all meetings from a specific owner.
func (r *MeetingRepository) GetAllMeetings(ownerID int) ([]dtos.MeetingOutput, error) {
	query := `SELECT id,title, description, start_time, end_time, owner_id, created_at FROM reunioes WHERE owner_id = $1`

	var lista []dtos.MeetingOutput

	rows, err := DB.Query(query, ownerID)
	if err != nil {
		return  nil, err
	}
	defer rows.Close()

	// Iterate through rows and append each meeting to the list.
	for rows.Next() {
		var u dtos.MeetingOutput

		err = rows.Scan(&u.ID, &u.Title, &u.Description, &u.StartTime, &u.EndTime, &u.OwnerID, &u.CreatedAt)
		if err != nil {
			return nil, err
		}

		lista = append(lista, u)
	}

	// Check for iteration errors.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lista, nil
}

// UpdateMeeting updates an existing meeting owned by the given user.
func (r *MeetingRepository) UpdateMeeting(meeting dtos.MeetingOutput ,ownerID int) (dtos.MeetingOutput, error) {
	query := `UPDATE reunioes SET title = $1, description = $2, start_time = $3, end_time = $4 WHERE id = $5 AND owner_id = $6
	RETURNING title, description, start_time, end_time, owner_id, created_at`

	var updateMeeting dtos.MeetingOutput

	// Execute update and return updated meeting data.
	err := DB.QueryRow(query, meeting.Title, meeting.Description, meeting.StartTime, meeting.EndTime, meeting.ID, ownerID).Scan(
		&updateMeeting.Title,
		&updateMeeting.Description,
		&updateMeeting.StartTime,
		&updateMeeting.EndTime,
		&updateMeeting.OwnerID,
		&updateMeeting.CreatedAt,
	)
	if err != nil {
		return dtos.MeetingOutput{}, err
	}

	return updateMeeting, nil
}

// DeleteMeeting removes a meeting by ID if it belongs to the given owner.
func (r *MeetingRepository) DeleteMeeting(id, ownerID int) error {
	query := `DELETE FROM reunioes WHERE id = $1 AND owner_id = $2`

	res, err := DB.Exec(query, id, ownerID)
	if err != nil {
		return err
	}

	// Check if any rows were deleted.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("nenhuma reunião encontrada para deletar")
	}

	return nil
}

// HasConflict checks if a meeting time range overlaps with another meeting of the same owner.
func (r *MeetingRepository) HasConflict(ownerID int, start, end time.Time, excludeID ...int) (bool, error) {
	query := `SELECT COUNT(*) FROM reunioes WHERE owner_id = $1 AND (
		(start_time < $3 AND end_time > $2) -- intervalo se sobrepõe
	)`
	args := []interface{}{ownerID, start, end}

	// Optionally exclude a meeting ID (useful when updating).
	if len(excludeID) > 0 {
		query += " AND id != $4"
		args = append(args, excludeID[0])
	}

	var count int 
	err := DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}