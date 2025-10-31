package realtime

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
	MeetingID int
	UserID int
}

type Message struct {
	MeetingID int
	UserID int
	Data []byte
}

func (c *Client) readPump(rtc *services.RTCService) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg dtos.SignalMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "offer":
			answerSDP, err := rtc.JoinSession(msg.MeetingID, msg.UserID, msg.Data)
			if err != nil {
				continue
			}

			resp := dtos.SignalMessage{
				Type: "answer",
				UserID: msg.UserID,
				MeetingID: msg.MeetingID,
				Data: answerSDP,
			}
			respBytes, _ := json.Marshal(resp)
			c.Send <- respBytes

		case "ice":
			err := rtc.AddIceCandidate(msg.MeetingID, msg.Data, msg.UserID)
			if err != nil {
				continue
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}