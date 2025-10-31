package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/realtime"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
)

func SetupWebSocketsRoutes(r *gin.Engine, hub *realtime.Hub, rtc *services.RTCService) {
	r.GET("/ws/:meetingID/:userID", func(c *gin.Context) {
		meetingID, _ := strconv.Atoi(c.Param("meetingID"))
		userID, _ := strconv.Atoi(c.Param("userID"))
		realtime.ServeWS(hub, rtc, c.Writer, c.Request, meetingID, userID)
	})
}