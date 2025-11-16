package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/services"
)

type MeetingController struct {
	Service *services.MeetingService
}

// CreateMeeting godoc
// @Summary Cria uma nova reunião
// @Description Cria uma reunião com título, descrição, horário de início e fim. O usuário logado será o proprietário (owner).
// @Tags meetings
// @Accept json
// @Produce json
// @Param meeting body dtos.MeetingInput true "Dados da reunião"
// @Success 201 {object} dtos.MeetingOutput
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /meetings [post]
// @Security BearerAuth
func (controller *MeetingController) CreateMeeting(c *gin.Context) {
	var meetingInput dtos.MeetingInput

	err := c.ShouldBindJSON(&meetingInput)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	layoutBR := "02/01/2006 15:04"  // dd/MM/yyyy HH:mm

	startTime, err := time.Parse(layoutBR, meetingInput.StartTime)
	if err != nil {
		c.JSON(400, gin.H{"error": "start_time inválido"})
		return
	}

	endTime, err := time.Parse(layoutBR, meetingInput.EndTime)
	if err != nil {
		c.JSON(400, gin.H{"error": "end_time inválido"})
		return
	}

	ownerID := c.GetInt("userID")

	meeting := dtos.Meeting{
		Title: meetingInput.Title,
		Description: meetingInput.Description,
		StartTime: startTime,
		EndTime: endTime,
		OwnerID: ownerID,
	}

	createdMeeting, err := controller.Service.CreateMeeting(meeting, ownerID)
	if err != nil {
		if strings.Contains(err.Error(), "schedule conflict") {
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, createdMeeting)
}

// GetAllMeetings godoc
// @Summary Lista todas as reuniões do usuário autenticado
// @Description Retorna apenas as reuniões criadas pelo usuário logado
// @Tags meetings
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} dtos.MeetingOutput
// @Failure 400 {object} map[string]string "usuário não autenticado ou id inválido"
// @Failure 500 {object} map[string]string "erro ao buscar reuniões"
// @Router /meetings [get]
func (controller *MeetingController) GetAllMeetings(c *gin.Context) {
	ownerID := c.GetInt("userID")

	meetings, err := controller.Service.GetAllMeetings(ownerID)
	if err != nil {
		c.JSON(500, gin.H{"error": "error fetching meetings"})
		return
	}

	c.JSON(200, meetings)
}

// UpdateController godoc
// @Summary Atualiza uma reunião
// @Description Atualiza os dados de uma reunião (título, descrição, horários). O usuário logado deve ser o dono da reunião.
// @Tags meetings
// @Accept json
// @Produce json
// @Param meeting body dtos.UpdateMeeting true "Dados da reunião com ID"
// @Success 200 {object} dtos.MeetingOutput
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /meetings/update [put]
// @Security BearerAuth
func (controller *MeetingController) UpdateController(c *gin.Context) {
	var meetingInput dtos.UpdateMeeting

	err := c.ShouldBindJSON(&meetingInput)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	layoutBR := "02/01/2006 15:04"

	if _, err := time.Parse(layoutBR, meetingInput.StartTime); err != nil {
		c.JSON(400, gin.H{"error": "invalid start_time"})
		return
	}

	if _, err := time.Parse(layoutBR, meetingInput.EndTime); err != nil {
		c.JSON(400, gin.H{"error": "invalid end_time"})
		return
	}

	ownerID := c.GetInt("userID")

	meetingUpdate, err := controller.Service.UpdateMeeting(meetingInput, ownerID)
	if err != nil {
		if strings.Contains(err.Error(), "schedule conflict") {
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, meetingUpdate)
}

// @Summary Deletar reunião
// @Description Remove uma reunião existente
// @Tags meetings
// @Accept  json
// @Produce  json
// @Param   id   path   int  true  "ID da reunião"
// @Success 200 {object} map[string]string "reunião excluída com sucesso"
// @Failure 400 {object} map[string]string "id inválido ou erro na exclusão"
// @Failure 401 {object} map[string]string "usuário não autenticado"
// @Failure 404 {object} map[string]string "reunião não encontrada"
// @Router /meetings/delete/{id} [delete]
// @Security BearerAuth
func (controller *MeetingController) DeleteMeeting(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	ownerID := c.GetInt("userID")

	err = controller.Service.DeleteMeeting(id, ownerID)
	if err != nil {
		c.JSON(400, gin.H{"error": "error deleting meeting"})
		return
	}

	c.JSON(200, gin.H{"message": "meeting deleted successfully"})
}