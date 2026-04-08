package otpservice

import (
	"e-commerce/utils"
	"fmt"
	"net/smtp"
)

type EmailSender struct{}

func (e *EmailSender) Send(identifier, otp string) error {

	body, err := utils.RenderTemplate("service/otp_service/templates/email_otp.html",map[string]interface{}{
		"OTP" : otp,
	})

	if err != nil {
		return err
	}

	from := "rahuljain10159@gmail.com"
	password := "wmsq tufs lznq luui"
	to := []string{identifier}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("Subject: Your OTP Code\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from,to, msg)
	if err != nil {
		return err
	}

	fmt.Println("Email sent successfully to", identifier)
	return nil

}