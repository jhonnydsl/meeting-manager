package repository

import (
	"fmt"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
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
		return err
	}
	if count > 0 {
		return fmt.Errorf("existing friendship or invitation")
	}
	
	query := `INSERT INTO friends (user_id, friend_id, status) VALUES ($1, $2, 'pending');`

	_, err = DB.Exec(query, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (r *FriendRepository) GetFriends(userID int) ([]dtos.FriendOutput, error) {
	query := `
	SELECT u.id, u.name
	FROM friends f
	JOIN users u ON u.id = f.friend_id
	WHERE f.user_id = $1 AND f.status = 'accepted'
	;`
	var friends []dtos.FriendOutput

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i dtos.FriendOutput 

		err := rows.Scan(&i.ID, &i.Name)
		if err != nil {
			return nil, err
		}

		friends = append(friends, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}

func (r *FriendRepository) GetFriendsPending(userID int) ([]dtos.FriendOutput, error) {
	query := `
	SELECT u.id, u.name
	FROM friends f
	JOIN users u ON u.id = f.friend_id
	WHERE f.user_id = $1 AND f.status = 'pending'
	;`
	var pendings []dtos.FriendOutput

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i dtos.FriendOutput

		err := rows.Scan(&i.ID, &i.Name)
		if err != nil {
			return nil, err
		}

		pendings = append(pendings, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pendings, nil
}

func (r *FriendRepository) AcceptFriend(friendID, userID int) error {
	query := `
	UPDATE friends SET status = 'accepted' WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	;`

	_, err := DB.Exec(query, friendID, userID)
	if err != nil {
		return err
	}

	return nil
}