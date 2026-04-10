package customerrequest

type VerifyOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	OTPType string `json:"otp_type" binding:required"`
	OTP string `json:"otp" binding:"required"`
}