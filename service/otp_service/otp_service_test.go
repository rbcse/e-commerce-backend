package otpservice_test

import (
	ae "e-commerce/error"
	"e-commerce/mocks/repomocks"
	"e-commerce/mocks/servicemocks"
	otprepository "e-commerce/repository/otp_repository"
	otp_service "e-commerce/service/otp_service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateOTP(t *testing.T) {
	testCases := []struct {
		description string
		mockSetup   func(r *repomocks.OTPRepository, og *servicemocks.OTPGenerator, sf *servicemocks.SenderFactory, ms *servicemocks.OTPSender)
		identifier  string
		otpType     string
		expectedOTP string
		expectedErr error
	}{
		{
			description: "Should return no error when the otp is generated successfully and the otp should be saved in Redis",
			mockSetup: func(r *repomocks.OTPRepository, og *servicemocks.OTPGenerator, sf *servicemocks.SenderFactory, ms *servicemocks.OTPSender) {
				og.On("Generate").Return("123456", nil)
				sf.On("GetSender", otp_service.OTPType("EMAIL")).Return(ms, nil)
				ms.On("Send", "rahul@gmail.com", "123456").Return(nil)
				r.On("SaveOTP", "rahul@gmail.com", "123456").Return(nil)
			},
			identifier:  "rahul@gmail.com",
			otpType:     "EMAIL",
			expectedOTP: "123456",
		},
		{
			description: "Should return error in generating otp when the otp is not generated successfully",
			mockSetup: func(r *repomocks.OTPRepository, og *servicemocks.OTPGenerator, sf *servicemocks.SenderFactory, ms *servicemocks.OTPSender) {
				og.On("Generate").Return("", ae.ErrOTPGenerationFailed)
			},
			identifier:  "rahul@gmail.com",
			otpType:     "EMAIL",
			expectedOTP: "",
			expectedErr: ae.ErrOTPGenerationFailed,
		},
		{
			description: "Should return error when an invalid otp type is send",
			mockSetup: func(r *repomocks.OTPRepository, og *servicemocks.OTPGenerator, sf *servicemocks.SenderFactory, ms *servicemocks.OTPSender) {
				og.On("Generate").Return("123456", nil)
				sf.On("GetSender", otp_service.OTPType("NOTHING")).Return(nil, ae.ErrInvalidOTPType)
			},
			identifier:  "r12djhv",
			otpType:     "NOTHING",
			expectedOTP: "123456",
			expectedErr: ae.ErrInvalidOTPType,
		},
		{
			description: "Should return error in sending otp when there is an error in sending otp to the identifier",
			mockSetup: func(r *repomocks.OTPRepository, og *servicemocks.OTPGenerator, sf *servicemocks.SenderFactory, ms *servicemocks.OTPSender) {
				og.On("Generate").Return("123456", nil)
				sf.On("GetSender", otp_service.OTPType("EMAIL")).Return(ms, nil)
				ms.On("Send", "rahul@gmail.com", "123456").Return(ae.ErrSendingOTP)
			},
			identifier:  "rahul@gmail.com",
			otpType:     "EMAIL",
			expectedOTP: "123456",
			expectedErr: ae.ErrSendingOTP,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			mockRepo := new(repomocks.OTPRepository)
			mockSenderFactory := new(servicemocks.SenderFactory)
			mockSender := new(servicemocks.OTPSender)
			mockOTPGenerator := new(servicemocks.OTPGenerator)

			tC.mockSetup(mockRepo, mockOTPGenerator, mockSenderFactory, mockSender)
			svc := otp_service.NewOTPService(mockRepo, mockOTPGenerator, mockSenderFactory)

			otp, err := svc.GenerateOTP(tC.identifier, tC.otpType)
			assert.Equal(t, tC.expectedOTP, otp)
			assert.Equal(t, tC.expectedErr, err)
		})
	}
}

func Test_VerifyOTP(t *testing.T) {
	testCases := []struct {
		description string
		mockSetup   func(r *repomocks.OTPRepository)
		identifier  string
		otpType     string
		otp         string
		expectedErr error
	}{
		{
			description: "Should return no error when the otp entered by customer matches the correct otp and otp is not expired",
			mockSetup: func(r *repomocks.OTPRepository) {
				r.On("GetOTP", "rahul@gmail.com").Return(&otprepository.OTPData{
					Identifier: "rahul@gmail.com",
					OTP:        "123456",
					Attempts:   3,
				})
			},
			identifier: "rahul@gmail.com",
			otpType:    "EMAIL",
			otp:        "123456",
		},
		{
			description: "Should return otp expired error when the customer enters the otp after 5 minutes the otp is being sent",
			mockSetup: func(r *repomocks.OTPRepository) {
				r.On("GetOTP", "rahul@gmail.com").Return(nil)
			},
			identifier: "rahul@gmail.com",
			otpType:    "EMAIL",
			otp:        "123456",
			expectedErr: ae.ErrOTPExpired,
		},
		{
			description: "Should return error attempts are exhausted when the customer enters the incorrect otp for 3 times and then entering otp for the 4th time",
			mockSetup: func(r *repomocks.OTPRepository) {
				r.On("GetOTP", "rahul@gmail.com").Return(&otprepository.OTPData{
					Identifier: "rahul@gmail.com",
					OTP:        "123456",
					Attempts:   0,
				})
			},
			identifier: "rahul@gmail.com",
			otpType:    "EMAIL",
			otp:        "123456",
			expectedErr: ae.ErrOTPAttemptsExhausted,
		},
		{
			description: "Should return wrong otp error when the otp entered by customer does not matches the correct otp",
			mockSetup: func(r *repomocks.OTPRepository) {
				r.On("GetOTP", "rahul@gmail.com").Return(&otprepository.OTPData{
					Identifier: "rahul@gmail.com",
					OTP:        "133456",
					Attempts:   3,
				})
			},
			identifier: "rahul@gmail.com",
			otpType:    "EMAIL",
			otp:        "123456",
			expectedErr: ae.ErrWrongOTPEntered,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			mockRepo := new(repomocks.OTPRepository)
			mockSenderFactory := new(servicemocks.SenderFactory)
			mockOTPGenerator := new(servicemocks.OTPGenerator)
			tC.mockSetup(mockRepo);
			svc := otp_service.NewOTPService(mockRepo, mockOTPGenerator, mockSenderFactory)
			err := svc.VerifyOTP(tC.identifier,tC.otpType,tC.otp);
			assert.Equal(t,tC.expectedErr,err);
		})
	}
}

func Test_DefaultOTPGenerator(t *testing.T) {
	testCases := []struct {
		description	string
		expectedErr error
		expectedOTPLength int
	}{
		{
			description: "Should return no error when GenerateRandom6DigitOTP is called and should generate a 6 digit otp",
			expectedOTPLength: 6,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			dg := otp_service.DefaultOTPGenerator{};
			otp , err := dg.Generate()
			assert.Equal(t,tC.expectedErr,err);
			assert.Equal(t,tC.expectedOTPLength,len(otp))
		})
	}
}