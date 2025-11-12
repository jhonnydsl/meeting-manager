package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/realtime"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/routes"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils/middleware"

	_ "github.com/jhonnydsl/gerenciamento-de-reunioes/docs"
)

func main() {
	// @title Gerenciamento de Reuniões API
	// @version 1.0
	// @description API para gerenciar usuários, reuniões e convites
	// @host meeting-manager-xzk7.onrender.com
	// @BasePath /api/v1

	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization

	/*err := godotenv.Load()
	if err != nil {
		log.Println("error loading environment variables")
		return
	}*/
	
	err := repository.Connect()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	} else {
		log.Println("connection established")
	}
	defer repository.DB.Close()
		
	repo := &repository.TableRepository{}

	err = repo.CreateTableUsers()
	if err != nil {
		log.Fatalf("error creating users table: %v", err)
	}	

	err = repo.CreateTableReunioes()
	if err != nil {
		log.Fatalf("error creating reunioes table: %v", err)
	}		

	err = repo.CreateTableConvites()
	if err != nil {		log.Fatalf("error creating convites table: %v", err)
	}
	
	err = repo.CreateTableFriends()
	if err != nil {
		log.Fatalf("error creating friends table: %v", err)
	}
		
	app := gin.Default()
	app.Use(middleware.Cors())
	app.Use(middleware.ErrorMiddlewareHandle())

	v1 := app.Group("/api/v1")
	{
		routes.SetupUserRoutes(v1)
		routes.SetupRoutesMeeting(v1)
		routes.SetupRoutesInvitation(v1)
		routes.SetupRoutesFriend(v1)
	}


	hub := realtime.NewHub()
	go hub.Run()
	rtcService := services.NewRTCService()

	routes.SetupWebSocketsRoutes(app, hub, rtcService)

	rtcService.OnICECandidate = func(meetingID, userID int, candidate string) {
		hub.Mutex.Lock()
		defer hub.Mutex.Unlock()

		clients, ok := hub.Clients[meetingID]
		if !ok {
			return
		}

		client, ok := clients[userID]
		if !ok {
			return
		}

		msg := dtos.SignalMessage{
			Type: "ice",
			MeetingID: meetingID,
			UserID: userID,
			Data: candidate,
		}

		data, _ := json.Marshal(msg)
		client.Send <- data
	}
	
	app.Run(":8080")
}