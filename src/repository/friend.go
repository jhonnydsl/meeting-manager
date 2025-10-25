package repository

import (
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils"
)

type FriendRepository struct{}

func (r *FriendRepository) AddFriend(userID, friendID int) error {
	existsQuery := `
	SELECT COUNT(*) FROM friends
	WHERE (user_id = $1 AND friend_id = $2)
	OR (user_id = $2 AND friend_id = $1)
	`
	var count int
	err := DB.QueryRow(existsQuery, userID, friendID).Scan(&count)
	if err != nil {
		return utils.BadRequestError("error sending request")
	}
	if count > 0 {
		return utils.ConflictError("existing friendship or invitation")
	}
	
	query := `INSERT INTO friends (user_id, friend_id, status) VALUES ($1, $2, 'pending');`

	_, err = DB.Exec(query, userID, friendID)
	if err != nil {
		return utils.BadRequestError("error sending request")
	}

	return nil
}

func (r *FriendRepository) GetFriends(userID int) ([]dtos.FriendOutput, error) {
	query := `
	SELECT u.id, u.name
	FROM friends f
	JOIN users u ON (u.id = f.friend_id AND f.user_id = $1) 
    OR (u.id = f.user_id AND f.friend_id = $1)
	WHERE f.status = 'accepted';`
	var friends []dtos.FriendOutput

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, utils.BadRequestError("error searching for friends")
	}
	defer rows.Close()

	for rows.Next() {
		var i dtos.FriendOutput 

		err := rows.Scan(&i.ID, &i.Name)
		if err != nil {
			return nil, utils.InternalServerError("error searching for friends")
		}

		friends = append(friends, i)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.InternalServerError("error searching for friends")
	}

	return friends, nil
}

func (r *FriendRepository) GetFriendsPending(userID int) ([]dtos.FriendOutput, error) {
	query := `
	SELECT u.id, u.name
	FROM friends f
	JOIN users u ON u.id = f.user_id
	WHERE f.friend_id = $1 AND f.status = 'pending'
	;`
	var pendings []dtos.FriendOutput

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, utils.BadRequestError("it was not possible to search requests with the data received")
	}
	defer rows.Close()

	for rows.Next() {
		var i dtos.FriendOutput

		err := rows.Scan(&i.ID, &i.Name)
		if err != nil {
			return nil, utils.InternalServerError("error listing pending requests")
		}

		pendings = append(pendings, i)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.InternalServerError("error listing pending requests")
	}

	return pendings, nil
}

func (r *FriendRepository) AcceptFriend(friendID, userID int) error {
	query := `
	UPDATE friends SET status = 'accepted' WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	;`

	_, err := DB.Exec(query, userID, friendID)
	if err != nil {
		return utils.BadRequestError("error accepting request")
	}

	return nil
}

func (r *FriendRepository) RefuseFriend(friendID, userID int) error {
	query := `
	UPDATE friends
	SET status = 'rejected'
	WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	;`

	_, err := DB.Exec(query, userID, friendID)
	if err != nil {
		return utils.BadRequestError("error refusing request")
	}

	return nil
}