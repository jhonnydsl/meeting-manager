package dtos

type AddFriendInput struct {
	FriendID int `json:"friend_id" binding:"required"`
}