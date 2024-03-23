package email

import (
	"fmt"
	"net/smtp"
	"server/internal/apperrors"
)

type Sender interface {
	Send(addr, topic, message string) error
}

type sender struct {
	auth        smtp.Auth
	senderEmail string
	senderSMTP  string
	port        string
}

func InitSender(senderEmail, senderKey, senderSMTP, port string) Sender {
	auth := smtp.PlainAuth("", senderEmail, senderKey, senderSMTP)
	return &sender{auth: auth, senderEmail: senderEmail, senderSMTP: senderSMTP, port: port}
}

func (s *sender) Send(addr, topic, message string) error {
	to := []string{addr}
	msgStr := fmt.Sprintf("To: %v\r\nSubject: %v\r\n\r\n %v\r\n", addr, topic, message)

	addrPort := fmt.Sprintf("%v:%v", s.senderSMTP, s.port)
	err := smtp.SendMail(addrPort, s.auth, addr, to, []byte(msgStr))
	if err != nil {
		appErr := apperrors.SendErr.AppendMessage(err)
		return appErr
	}

	return nil
}
