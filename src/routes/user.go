package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/controllers"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(app *gin.Engine) {
	userService := &services.UserService{UserRepo: &repository.UserRepository{}}
	userController := &controllers.UserController{Service: userService}

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server online"})
	})

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	users := app.Group("/users")
	{
		users.POST("", userController.CreateUser)
		users.GET("/", middleware.AuthMiddleware(), userController.GetAllUsers)
		users.GET("/profile", middleware.AuthMiddleware(), userController.GetProfile)
		users.DELETE("/", middleware.AuthMiddleware(), userController.DeleteUser)
	}

	app.POST("/login", userController.LoginUser)
}