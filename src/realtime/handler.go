package realtime

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, rtc *services.RTCService, w http.ResponseWriter, r *http.Request, meetingID, userID int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Hub: hub,
		Conn: conn,
		Send: make(chan []byte, 256),
		MeetingID: meetingID,
		UserID: userID,
	}

	hub.Register <- client

	go client.writePump()
	go client.readPump(rtc)
}