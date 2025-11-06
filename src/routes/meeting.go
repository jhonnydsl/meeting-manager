package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/controllers"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils/middleware"
)

func SetupRoutesMeeting(app *gin.Engine) {
	meetingService := &services.MeetingService{MeetingRepo: &repository.MeetingRepository{}}
	meetingController := &controllers.MeetingController{Service: meetingService}

	meetings := app.Group("/meetings", middleware.AuthMiddleware())
	{
		meetings.POST("/", meetingController.CreateMeeting)
		meetings.GET("", meetingController.GetAllMeetings)
		meetings.PUT("/update", meetingController.UpdateController)
		meetings.DELETE("/delete/:id", meetingController.DeleteMeeting)
	}
}