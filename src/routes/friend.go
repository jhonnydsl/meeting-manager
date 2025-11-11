package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/controllers"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils/middleware"
)

func SetupRoutesFriend(app *gin.RouterGroup) {
	friendService := &services.FriendService{FriendRepo: &repository.FriendRepository{}}
	friendController := &controllers.FriendController{Service: friendService}

	friends := app.Group("/friends", middleware.AuthMiddleware())
	{
		friends.POST("/", friendController.AddFriend)
		friends.GET("/", friendController.GetFriends)
		friends.GET("/pendings", friendController.GetFriendsPending)
		friends.PUT("/accept", friendController.AcceptedFriend)
		friends.PUT("/refused", friendController.RefuseFriend)
	}
}