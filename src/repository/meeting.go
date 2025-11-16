package repository

import (
	"time"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils"
)

type MeetingRepository struct{}

func (r *MeetingRepository) CreateMeeting(meeting dtos.Meeting, ownerId int) (dtos.MeetingOutput, error) {
	query := `INSERT INTO reunioes (title, description, start_time, end_time, owner_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, description, start_time, end_time, owner_id, created_at`
	var createdMeeting dtos.MeetingOutput

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
		return dtos.MeetingOutput{}, utils.InternalServerError("error creating meeting")
	}

	return createdMeeting, nil
}

func (r *MeetingRepository) GetAllMeetings(ownerID int) ([]dtos.MeetingOutput, error) {
	query := `SELECT id,title, description, start_time, end_time, owner_id, created_at FROM reunioes WHERE owner_id = $1`

	var lista []dtos.MeetingOutput

	rows, err := DB.Query(query, ownerID)
	if err != nil {
		return  nil, utils.InternalServerError("error fetching meeting")
	}
	defer rows.Close()

	for rows.Next() {
		var u dtos.MeetingOutput

		err = rows.Scan(&u.ID, &u.Title, &u.Description, &u.StartTime, &u.EndTime, &u.OwnerID, &u.CreatedAt)
		if err != nil {
			return nil, utils.InternalServerError("error fetching meeting")
		}

		lista = append(lista, u)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.InternalServerError("error fetching meeting")
	}

	return lista, nil
}

func (r *MeetingRepository) UpdateMeeting(meetingInput dtos.UpdateMeeting, meetingID int, ownerID int, start time.Time, end time.Time) (dtos.MeetingOutput, error) {
    query := `
        UPDATE reunioes
        SET title = $1, description = $2, start_time = $3, end_time = $4
        WHERE id = $5 AND owner_id = $6
        RETURNING id, title, description, start_time, end_time, owner_id, created_at
    `

    var updated dtos.MeetingOutput

    err := DB.QueryRow(
        query,
        meetingInput.Title,
        meetingInput.Description,
        start,
        end,
        meetingID,
        ownerID,
    ).Scan(
        &updated.ID,
        &updated.Title,
        &updated.Description,
        &updated.StartTime,
        &updated.EndTime,
        &updated.OwnerID,
        &updated.CreatedAt,
    )

    if err != nil {
        return dtos.MeetingOutput{}, utils.InternalServerError("error updating meeting")
    }

    return updated, nil
}

func (r *MeetingRepository) DeleteMeeting(id, ownerID int) error {
	query := `DELETE FROM reunioes WHERE id = $1 AND owner_id = $2`

	res, err := DB.Exec(query, id, ownerID)
	if err != nil {
		return utils.InternalServerError("error deleting meeting")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.InternalServerError("error deleting meeting")
	}

	if rowsAffected == 0 {
		return utils.NotFoundError("no meeting found to delete")
	}

	return nil
}

func (r *MeetingRepository) HasConflict(ownerID int, start, end time.Time, excludeID ...int) (bool, error) {
	query := `SELECT COUNT(*) FROM reunioes WHERE owner_id = $1 AND (
		(start_time < $3 AND end_time > $2) -- intervalo se sobrepõe
	)`
	args := []interface{}{ownerID, start, end}

	if len(excludeID) > 0 {
		query += " AND id != $4"
		args = append(args, excludeID[0])
	}

	var count int 
	err := DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, utils.InternalServerError("error checking meeting conflicts")
	}

	return count > 0, nil
}