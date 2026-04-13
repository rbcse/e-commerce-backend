package otpservice_test

import (
    "e-commerce/mocks/servicemocks"
    otpservice "e-commerce/service/otp_service"
    "fmt"
    "testing"
	ae "e-commerce/error"
    "github.com/stretchr/testify/assert"
)

func buildEmailMessage(subject, body string) []byte {
    msg := fmt.Sprintf("Subject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s", subject, body)
    return []byte(msg)
}

func Test_EmailSender(t *testing.T) {
    testCases := []struct {
        description string
        mockSetup   func(m *servicemocks.Mailer, r *servicemocks.TemplateRenderer)
        identifier  string
        otp         string
        expectedErr error
    }{
        {
            description: "Should return no error when the template for the email is rendered properly and the message is sent successfully",
            mockSetup: func(m *servicemocks.Mailer, r *servicemocks.TemplateRenderer) {
                renderedBody := "OTP : 123456"

                r.On("Render", "service/otp_service/templates/email_otp.html", map[string]interface{}{
                    "OTP": "123456",
                }).Return(renderedBody, nil)

                expectedMsg := buildEmailMessage("Your OTP Code", renderedBody)
                m.On("Send", []string{"rahul@gmail.com"}, expectedMsg).Return(nil)
            },
            identifier:  "rahul@gmail.com",
            otp:         "123456",
            expectedErr: nil,
        },
        {
            description: "Should return error in rendering the template when the template for the email is not rendered",
            mockSetup: func(m *servicemocks.Mailer, r *servicemocks.TemplateRenderer) {

                r.On("Render", "service/otp_service/templates/email_otp.html", map[string]interface{}{
                    "OTP": "123456",
                }).Return("",ae.ErrRenderingTemplate)

            },
            identifier:  "rahul@gmail.com",
            otp:         "123456",
            expectedErr: ae.ErrRenderingTemplate,
        },
    }

    for _, tC := range testCases {
        t.Run(tC.description, func(t *testing.T) {
            mockMailer := new(servicemocks.Mailer)
            mockRenderer := new(servicemocks.TemplateRenderer)

            tC.mockSetup(mockMailer, mockRenderer)

            sender := otpservice.NewEmailSender(mockMailer, mockRenderer)
            err := sender.Send(tC.identifier, tC.otp)

            assert.Equal(t, tC.expectedErr, err)
            mockMailer.AssertExpectations(t)
            mockRenderer.AssertExpectations(t)
        })
    }
}