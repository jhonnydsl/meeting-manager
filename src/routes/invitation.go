package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/controllers"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils/middleware"
)

func SetupRoutesInvitation(app *gin.RouterGroup) {
	invitationService := &services.InvitationService{InvitRepo: &repository.InvitationRepository{}}
	invitationController := &controllers.InvitationController{Service: invitationService}

	invitations := app.Group("/invitations", middleware.AuthMiddleware())
	{
		invitations.POST("/", invitationController.CreateInvitation)
		invitations.GET("/sent", invitationController.GetAllInvitations)
		invitations.GET("/received", invitationController.GetReceiver)
		invitations.DELETE("/:id", invitationController.DeleteInvitation)
	}
}