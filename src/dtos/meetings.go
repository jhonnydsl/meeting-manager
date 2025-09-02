package dtos

import "time"

type MeetingInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime 	string `json:"end_time" binding:"required"`
}

type MeetingOutput struct {
	ID 			int `json:"id"`
	Title 		string `json:"title"`
	Description string `json:"description"`
	StartTime 	time.Time `json:"start_time"`
	EndTime 	time.Time `json:"end_time"`
	OwnerID 	int `json:"owner_id"`
	CreatedAt 	time.Time `json:"created_at"`
}

type Meeting struct {
	Title 		string
	Description string
	StartTime 	time.Time
	EndTime 	time.Time
	OwnerID 	int
}