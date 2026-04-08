package customerrequest

type GenerateOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Type       string `json:"type" binding:"required"`
}
