package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
)

type UserController struct {
	Service *services.UserService
}

// CreateUser godoc
// @Summary Cria um novo usuário
// @Description Cria um usuário com nome, email e senha
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.UserInput true "Dados do usuário"
// @Success 201 {object} dtos.UserOutput
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (controller * UserController) CreateUser(c *gin.Context) {
	var userInput dtos.UserInput
	err := c.ShouldBindJSON(&userInput)		// <= Lendo os dados em formato JSON para a criação de usuario.
	if err != nil {	
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := controller.Service.CreateUser(userInput)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, createdUser)
}

// GetAllUsers godoc
// @Summary Lista todos os usuários
// @Security BearerAuth
// @Description Retorna todos os usuários cadastrados
// @Tags users
// @Produce json
// @Success 200 {array} dtos.UserOutput
// @Failure 400 {object} map[string]string
// @Router /users [get]
func (controller *UserController) GetAllUsers(c *gin.Context) {
	lista, err := controller.Service.GetAllUsers()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, lista)
}

// LoginUser godoc
// @Summary Autentica um usuário
// @Description Faz login com email e senha e retorna um token JWT
// @Tags users
// @Accept json
// @Produce json
// @Param login body dtos.UserLoginInput true "Dados para login"
// @Success 200 {object} map[string]string "Token JWT"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (controller *UserController) LoginUser(c *gin.Context) {
	var login dtos.UserLoginInput

	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := controller.Service.LoginUser(login)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

// GetProfile godoc
// @Summary Retorna o perfil do usuário autenticado
// @Description Obtém os dados do usuário logado a partir do token JWT
// @Tags users
// @Produce json
// @Success 200 {object} dtos.UserOutput
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/profile [get]
// @Security BearerAuth
func (controller *UserController) GetProfile(c *gin.Context) {
	userID := c.GetInt("userID")

	user, err := controller.Service.GetUserByID(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "usuário não encontrado"})
		return
	}

	c.JSON(200, user)
}

// DeleteUser godoc
// @Summary Deleta o usuário autenticado
// @Description Deleta o usuário logado a partir do token JWT
// @Tags users
// @Produce json
// @Success 200 {object} map[string]string "Mensagem de sucesso"
// @Failure 401 {object} map[string]string "Usuário não autenticado"
// @Failure 400 {object} map[string]string "Erro ao deletar usuário"
// @Router /users [delete]
// @Security BearerAuth
func (controller *UserController) DeleteUser(c *gin.Context) {
	userID := c.GetInt("userID")

	if err := controller.Service.DeleteUser(userID); err != nil {
		c.JSON(400, gin.H{"error": "erro ao deletar usuario"})
		return
	}

	c.JSON(200, gin.H{"message": "usuario deletado com sucesso"})
}