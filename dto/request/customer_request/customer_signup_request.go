package customerrequest

type CustomerSignupRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,min=12,max=14"`
	Password string `json:"password" binding:"required,min=8,max=15"`
}
