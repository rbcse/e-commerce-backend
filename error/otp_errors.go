package error

import "errors"

var (
	ErrOTPAttemptsExhausted = errors.New("OTP Attempts are exhausted");
	ErrWrongOTPEntered = errors.New("Incorrect OTP entered")
	ErrOTPExpired = errors.New("OTP expired")
	ErrOTPGenerationFailed = errors.New("Failed to generated otp");
	ErrInvalidOTPType = errors.New("Invalid otp type")
)