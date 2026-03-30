package customerrequest

type CustomerSignupRequest struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}