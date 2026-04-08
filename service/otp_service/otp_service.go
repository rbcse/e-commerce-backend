package otpservice

import (
	"crypto/rand"
	otprepository "e-commerce/repository/otp_repository"
	"fmt"
	"math/big"
)

type OTPService interface {
	GenerateOTP(identifier, otp_type string) (string, error)
}

type otpService struct {
	repo otprepository.OTPRepository
}

func NewOTPService(repo otprepository.OTPRepository) OTPService {
	return &otpService{
		repo: repo,
	}
}

func generateRandom6DigitOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64()), err
}

func (os *otpService) GenerateOTP(identifier, otp_type string) (string, error) {

	otp, err := generateRandom6DigitOTP()
	if err != nil {
		return "", err
	}

	sender , err := GetSender(OTPType(otp_type))

	if err != nil {
		return "" , err
	}

	err = sender.Send(identifier,otp);

	if err != nil {
		return "" , err
	}

	err = os.repo.SaveOTP(identifier, otp)
	return otp, err

}
