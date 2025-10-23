package dtos

type AddFriendInput struct {
	FriendID int `json:"friend_id" binding:"required"`
}

type FriendOutput struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AcceptedFriend struct {
	ID int `json:"id"`
}