package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/routes"
	"github.com/joho/godotenv"

	_ "github.com/jhonnydsl/gerenciamento-de-reunioes/docs"
)

func main() {
	// @title Gerenciamento de Reuniões API
	// @version 1.0
	// @description API para gerenciar usuários, reuniões e convites
	// @host localhost:8080
	// @BasePath /

	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization

	err := godotenv.Load()
	if err != nil {
		log.Println("erro ao carregar variaveis de ambiente")
		return
	}
	
	err = repository.Connect()
	if err != nil {
		log.Fatalf("erro ao conectar ao banco de dados: %v", err)
	} else {
		log.Println("Conexão estabelecida.")
	}
	defer repository.DB.Close()
		
	repo := &repository.TableRepository{}
		
	err = repo.CreateTableUsers()
	if err != nil {
		log.Fatalf("erro ao criar tabela users: %v", err)
	}
		
	err = repo.CreateTableReunioes()
	if err != nil {
		log.Fatalf("erro ao criar tabela reunioes: %v", err)
	}
		
	err = repo.CreateTableConvites()
	if err != nil {
		log.Fatalf("erro ao criar tabela convites: %v", err)
	}
		
	app := gin.Default()
	routes.SetupRoutes(app)
	routes.SetupRoutesMeeting(app)
	routes.SetupRoutesInvitation(app)
		
	app.Run(":8080")
}