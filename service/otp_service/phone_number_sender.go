package otpservice

import "fmt"

type PhoneNumberSender struct{}

func (p *PhoneNumberSender) Send(identifier string, otp string) error {
	fmt.Println("Sending SMS OTP to", identifier, "OTP:", otp)
	return nil
}