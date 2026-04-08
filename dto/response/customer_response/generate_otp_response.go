package customerresponse

type GenerateOTPResponse struct {
	Otp     string `json:"otp" binding:"required"`
	Message string `json:"error_message" binding:"required"`
}
