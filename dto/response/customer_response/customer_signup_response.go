package customerresponse

type CustomerSignupResponse struct {
	IsSignUpSuccessful bool   `json:"is_signup_successful" binding:"required"`
	Message            string `json:"message" binding:"required"`
}
