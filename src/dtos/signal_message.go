package dtos

type SignalMessage struct {
	Type      string `json:"type"`
	UserID    int    `json:"userID"`
	MeetingID int    `json:"meetingID"`
	Data      string `json:"data"`
}