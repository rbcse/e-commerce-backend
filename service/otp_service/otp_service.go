package otpservice

import (
	ae "e-commerce/error"
	otprepository "e-commerce/repository/otp_repository"
	"e-commerce/utils"
)

type OTPService interface {
	GenerateOTP(identifier, otp_type string) (string, error)
	VerifyOTP(identifier , otp_type , otp string) error
}

type OTPGenerator interface {
	Generate() (string , error)
}

type DefaultOTPGenerator struct {}

func (d *DefaultOTPGenerator) Generate() (string,error) {
	return utils.GenerateRandom6DigitOTP()
}

type SenderFactory interface {
	GetSender(OTPType) (OTPSender , error)
}

type otpService struct {
	repo otprepository.OTPRepository
	otpGenerator OTPGenerator
	senderFactory SenderFactory
}

func NewOTPService(repo otprepository.OTPRepository , otpGenerator OTPGenerator , senderFactory SenderFactory) OTPService {
	return &otpService{
		repo: repo,
		otpGenerator: otpGenerator,
		senderFactory: senderFactory,
	}
}

func (os *otpService) GenerateOTP(identifier, otp_type string) (string, error) {

	otp, err := os.otpGenerator.Generate();
	if err != nil {
		return "", err
	}

	sender , err := os.senderFactory.GetSender(OTPType(otp_type))

	if err != nil {
		return otp , err
	}

	err = sender.Send(identifier,otp);

	if err != nil {
		return otp , err
	}

	err = os.repo.SaveOTP(identifier, otp)
	return otp, err

}

func (os *otpService) VerifyOTP(identifier , otp_type , otp string) (error) {

	otpData , err := os.repo.GetOTP(identifier)
	if err != nil {
		return ae.ErrOTPExpired
	}

	if otpData.Attempts == 0{
		return ae.ErrOTPAttemptsExhausted
	}

	if otpData.OTP != otp {
		return ae.ErrWrongOTPEntered
	}

	return nil

}