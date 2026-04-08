package otpservice

type OTPSender interface{
	Send(identifier , otp string) error
}