package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
)

type FriendController struct {
	Service *services.FriendService
}

// AddFriend godoc
// @Summary Envia solicitação de amizade
// @Description Envia uma solicitação de amizade ao id selecionado
// @Tags friends
// @Accept json
// @Produce json
// @Param invitation body dtos.AddFriendInput true "Dados da solicitação"
// @Success 201 {object} map[string]string "request sent successfully"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /friends [post]
// @Security BearerAuth
func (controller *FriendController) AddFriend(c *gin.Context) {
	var input dtos.AddFriendInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("userID")

	err = controller.Service.AddFriend(userID, input.FriendID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "request sent successfully"})
}

// GetFriends godoc
// @Summary Lista todos os amigos do usuário logado
// @Description Retorna uma lista de amigos do usuário autenticado.
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dtos.FriendOutput
// @Failure 400 {object} map[string]string "user id invalid or query error"
// @Failure 401 {object} map[string]string "user unauthorizad"
// @Router /friends [get]
func (controller *FriendController) GetFriends(c *gin.Context) {
	userID := c.GetInt("userID")

	friends, err := controller.Service.GetFriends(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, friends)
}