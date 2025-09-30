package utils

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
)

func SendInvitationEmail(to string, inv dtos.InvitationOutput) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	subject := "Convite para reunião"
	body := fmt.Sprintf("You have received a meeting invitation ID %d.\nStatus: %s", inv.ReuniaoID, inv.Status)

	message := []byte("Subject: " + subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth("", user, pass, host)
	err := smtp.SendMail(host+":"+port, auth, user, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}