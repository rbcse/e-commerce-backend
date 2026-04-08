package otpservice

import "fmt"

type OTPType string

const (
	Email OTPType = "EMAIL"
	Phone OTPType = "PHONE"
)

func GetSender(t OTPType) (OTPSender, error) {
	switch t {
	case Email:
		return &EmailSender{}, nil
	case Phone:
		return &PhoneNumberSender{}, nil
	default:
		return nil, fmt.Errorf("invalid otp type")
	}
}