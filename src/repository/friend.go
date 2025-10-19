package repository

import "fmt"

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