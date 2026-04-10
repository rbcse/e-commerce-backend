package otpservice_test

import (
	"e-commerce/mocks/repomocks"
	"e-commerce/mocks/servicemocks"
	otp_service "e-commerce/service/otp_service"
	"testing"
	ae "e-commerce/error"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateOTP(t *testing.T) {
	testCases := []struct {
		description	string
		mockSetup func(r *repomocks.OTPRepository , og *servicemocks.OTPGenerator , sf *servicemocks.SenderFactory , ms *servicemocks.OTPSender)
		identifier string
		otpType string
		expectedOTP string
		expectedErr error
	}{
		{
			description: "Should return no error when the otp is generated successfully and the otp should be saved in Redis",
			mockSetup: func(r *repomocks.OTPRepository, og *servicemocks.OTPGenerator, sf *servicemocks.SenderFactory, ms *servicemocks.OTPSender) {
				og.On("Generate").Return("123456",nil);
				sf.On("GetSender",otp_service.OTPType("EMAIL")).Return(ms,nil);
				ms.On("Send","rahul@gmail.com","123456").Return(nil);
				r.On("SaveOTP","rahul@gmail.com","123456").Return(nil)
			},
			identifier: "rahul@gmail.com",
			otpType: "EMAIL",
			expectedOTP: "123456",
		},
		{
			description: "Should return error in generating otp when the otp is not generated successfully",
			mockSetup: func(r *repomocks.OTPRepository, og *servicemocks.OTPGenerator, sf *servicemocks.SenderFactory, ms *servicemocks.OTPSender) {
				og.On("Generate").Return("",ae.ErrOTPGenerationFailed);
			},
			identifier: "rahul@gmail.com",
			otpType: "EMAIL",
			expectedOTP: "",
			expectedErr: ae.ErrOTPGenerationFailed,
		},
		{
			description: "Should return error when an invalid otp type is send",
			mockSetup: func(r *repomocks.OTPRepository, og *servicemocks.OTPGenerator, sf *servicemocks.SenderFactory, ms *servicemocks.OTPSender) {
				og.On("Generate").Return("123456",nil);
				sf.On("GetSender",otp_service.OTPType("NOTHING")).Return(nil,ae.ErrInvalidOTPType);
			},
			identifier: "r12djhv",
			otpType: "NOTHING",
			expectedOTP: "123456",
			expectedErr: ae.ErrInvalidOTPType,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			mockRepo := new(repomocks.OTPRepository)
			mockSenderFactory := new(servicemocks.SenderFactory)
			mockSender := new(servicemocks.OTPSender)
			mockOTPGenerator := new(servicemocks.OTPGenerator)

			tC.mockSetup(mockRepo,mockOTPGenerator,mockSenderFactory,mockSender);
			svc := otp_service.NewOTPService(mockRepo,mockOTPGenerator,mockSenderFactory);

			otp , err := svc.GenerateOTP(tC.identifier,tC.otpType);
			assert.Equal(t,tC.expectedOTP,otp);
			assert.Equal(t,tC.expectedErr,err);
		})
	}
}