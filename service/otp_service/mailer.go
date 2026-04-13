package otpservice

import (
	"e-commerce/common/logger"
	"net/smtp"
)

type Mailer interface {
	Send(to []string, message []byte) error
}

type SMTPMailer struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewSMTPMailer(from, password, smtpHost, smtpPort string) Mailer {
	return &SMTPMailer{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (s *SMTPMailer) Send(to []string, message []byte) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)

	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.from, to, message)
	if err != nil {
		return err
	}

	logger.Info("Email sent successfully")
	return nil
}