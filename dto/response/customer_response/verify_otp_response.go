package customerresponse

type VerifyOTPResponse struct {
	IsOTPVerified bool `json:"is_otp_verified" binding:"required"`
	Message string `json:"message" binding:"required"`
}