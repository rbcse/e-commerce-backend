package otpservice

import (
	ae "e-commerce/error"
	"e-commerce/templates"
)

type OTPType string

const (
	Email OTPType = "EMAIL"
	Phone OTPType = "PHONE"
)

type DefaultSenderFactory struct{}

func (d *DefaultSenderFactory) GetSender(t OTPType) (OTPSender, error) {
	switch t {
	case Email:
		mailer := NewSMTPMailer("rahuljain10159@gmail.com", "wmsq tufs lznq luui", "smtp.gmail.com", "587")
		return NewEmailSender(mailer, &templates.HTMLRenderer{}), nil
	case Phone:
		return &PhoneNumberSender{}, nil
	default:
		return nil, ae.ErrInvalidOTPType
	}
}
