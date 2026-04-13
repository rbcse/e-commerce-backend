package otpservice_test

import (
	ae "e-commerce/error"
	otpservice "e-commerce/service/otp_service"
	"testing"
	"github.com/stretchr/testify/assert"
)	

func Test_GetSender(t *testing.T) {
	testCases := []struct {
		description string
		otpType     otpservice.OTPType
		expectedErr error
	}{
		{
			description: "Should return EmailSender as the otp sender and no error when the otpType is EMAIL",
			otpType: otpservice.Email,
		},
		{
			description: "Should return PhoneNumberSender as the otp sender and no error when the otpType is PHONE",
			otpType: otpservice.Phone,
		},
		{
			description: "Should return Invalid OTP type error when the otp type is other than EMAIL and PHONE",
			otpType: otpservice.OTPType("RANDOMTYPE"),
			expectedErr: ae.ErrInvalidOTPType,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			dsf := otpservice.DefaultSenderFactory{};
			otpSender , err := dsf.GetSender(tC.otpType);
			if tC.expectedErr == nil {
				assert.NotNil(t,otpSender);
			}
			assert.Equal(t,tC.expectedErr,err);
		})
	}
}