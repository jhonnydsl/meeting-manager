package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils"
)

type InvitationController struct {
	Service *services.InvitationService
}

// CreateInvitation godoc
// @Summary Cria um novo convite
// @Description Cria um convite para outro usuário participar de uma reunião
// @Tags invitations
// @Accept json
// @Produce json
// @Param invitation body dtos.InvitationInput true "Dados do convite"
// @Success 201 {object} dtos.InvitationOutput
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /invitations [post]
// @Security BearerAuth
func (controller *InvitationController) CreateInvitation(c *gin.Context) {
	var invitationInput dtos.InvitationInput

	err := c.ShouldBindJSON(&invitationInput)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	senderID := c.GetInt("userID")
	
	createdInvitation, err := controller.Service.CreateInvitation(invitationInput, senderID)
	if err != nil {
		c.JSON(400, gin.H{"error": "erro ao criar convite"})
		return
	}

	receiverEmail, err := controller.Service.ReturnUserByEmail(createdInvitation.ReceiverID)
	if err != nil {
		fmt.Println("Erro ao pegar email:", err)
	} else {
		go func(invID int, email string) {
			if err := utils.SendInvitationEmail(email, createdInvitation); err != nil {
				fmt.Println("erro ao enviar email:", err)
				return
			} 
			
			if err := controller.Service.UpdateInvitationStatus(invID, "sent"); err != nil {
				fmt.Println("erro ao atualizar status", err)
			}
		}(createdInvitation.ID, receiverEmail)
	}

	c.JSON(201, gin.H{"message": "convite enviado com sucesso"})
}

// GetAllInvitations godoc
// @Summary Lista todos os convites enviados pelo usuário logado
// @Description Retorna uma lista de convites em que o usuário autenticado é o remetente (sender).
// @Tags invitations
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dtos.InvitationOutput
// @Failure 400 {object} map[string]string "id do usuário inválido ou erro de query"
// @Failure 401 {object} map[string]string "usuário não autenticado"
// @Router /invitations/sent [get]
func (controller *InvitationController) GetAllInvitations(c *gin.Context) {
	senderID := c.GetInt("userID")

	invitations, err := controller.Service.GetAllInvitations(senderID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, invitations)
}

// GetAllInvitations godoc
// @Summary Lista todos os convites recebidos pelo usuário logado
// @Description Retorna uma lista de convites em que o usuário autenticado recebeu (receiver).
// @Tags invitations
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dtos.InvitationOutput
// @Failure 400 {object} map[string]string "id do usuário inválido ou erro de query"
// @Failure 401 {object} map[string]string "usuário não autenticado"
// @Router /invitations/received [get]
func (controller *InvitationController) GetReceiver(c *gin.Context) {
	receiverID := c.GetInt("userID")

	invitations, err := controller.Service.GetReceiver(receiverID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, invitations)
}

// DeleteInvitation godoc
// @Summary Excluir convite
// @Description Permite que apenas o dono da reunião exclua um convite
// @Tags invitations
// @Param id path int true "ID do convite"
// @Success 200 {object} map[string]string "convite deletado com sucesso"
// @Failure 400 {object} map[string]string "erro ao excluir convite"
// @Failure 401 {object} map[string]string "usuário não autenticado"
// @Security BearerAuth
// @Router /invitations/{id} [delete]
func (controller *InvitationController) DeleteInvitation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ownerID := c.GetInt("userID")

	err = controller.Service.DeleteInvitation(id, ownerID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "convite deletado com sucesso"})
}