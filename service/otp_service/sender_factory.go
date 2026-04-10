package otpservice

import ae "e-commerce/error"

type OTPType string

const (
	Email OTPType = "EMAIL"
	Phone OTPType = "PHONE"
)

type DefaultSenderFactory struct{}

func (d *DefaultSenderFactory) GetSender(t OTPType) (OTPSender, error) {
	switch t {
	case Email:
		return &EmailSender{}, nil
	case Phone:
		return &PhoneNumberSender{}, nil
	default:
		return nil, ae.ErrInvalidOTPType
	}
}