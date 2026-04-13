package otpservice

import (
	"e-commerce/templates"
)

type EmailSender struct {
	mailer           Mailer
	templateRenderer templates.TemplateRenderer
}

func NewEmailSender(mailer Mailer, templateRenderer templates.TemplateRenderer) *EmailSender {
	return &EmailSender{
		mailer:           mailer,
		templateRenderer: templateRenderer,
	}
}

func (e *EmailSender) Send(identifier, otp string) error {

	body, err := e.templateRenderer.Render("service/otp_service/templates/email_otp.html", map[string]interface{}{
		"OTP": otp,
	})

	if err != nil {
		return err
	}

	to := []string{identifier}

	msg := []byte("Subject: Your OTP Code\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body)

	err = e.mailer.Send(to, msg)
	return err

}
