package dtos

import "time"

type InvitationInput struct {
	ReuniaoID  int    `json:"reuniao_id" binding:"required"`
	ReceiverID int    `json:"receiver_id" binding:"required"`
}

type InvitationOutput struct {
	ID         int    `json:"id"`
	ReuniaoID  int    `json:"reuniao_id"`
	ReceiverID int    `json:"receiver_id"`
	Status     string `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	SenderID   int 	`json:"sender_id"`
}
