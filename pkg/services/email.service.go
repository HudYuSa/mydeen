package services

import (
	"bytes"
	_ "embed"
	"net/smtp"
	"text/template"

	"github.com/HudYuSa/mydeen/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailVerificationData struct {
	VerificationCode string
}

type OtpData struct {
	Code       string
	ExpireTime int // minutes
}

//go:embed templates/email_verification.html
var verificationEmailTemplate string

//go:embed templates/otp.html
var otpTemplate string

func SendVerificationCodeSMTP(verificationCode string, to []string) error {
	// Create an email template from the embedded content.
	tmpl, err := template.New("emailVerificationTemplate").Parse(verificationEmailTemplate)
	if err != nil {
		return err
	}

	// Replace the placeholder with the actual verification link.
	data := EmailVerificationData{
		VerificationCode: "http://localhost:8000/api/master/verify?token=" + verificationCode,
	}

	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, data); err != nil {
		return err
	}
	emailSubject := "Subject: Account Verification Code\r\n"
	emailMessage := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		emailSubject + emailContent.String()

	auth := smtp.PlainAuth("", "hudyusufatsigah@gmail.com", config.GlobalConfig.GoogleAppPassword, "smtp.gmail.com")

	err = smtp.SendMail("smtp.gmail.com:587", auth, "hudyusufatsigah@gmail.com", to, []byte(emailMessage))
	if err != nil {
		return err
	}

	return nil
}

func SendVerificationCodeGomail(verificationCode string, to []string) error {
	// Create an email template from the embedded content.
	tmpl, err := template.New("emailTemplate").Parse(verificationEmailTemplate)
	if err != nil {
		return err
	}

	// Replace the placeholder with the actual verification link.
	data := EmailVerificationData{
		VerificationCode: "http://" + config.GlobalConfig.ServerUrl + "/api/master/verify?verification_code=" + verificationCode,
	}

	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, data); err != nil {
		return err
	}

	// create the message
	m := gomail.NewMessage()
	m.SetHeader("From", "hudyusufatsigah@gmail.com")
	m.SetHeader("To", to...)
	m.SetHeader("Content-Type", "text/html; charset=UTF-8")
	m.SetHeader("Subject", "Account Verification Code")
	m.SetBody("text/html", emailContent.String())

	// create dialer to send message
	d := gomail.NewDialer("smtp.gmail.com", 587, "hudyusufatsigah@gmail.com", config.GlobalConfig.GoogleAppPassword)

	// send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendOtpCode(otpCode string, to []string) error {
	// Create an email template from the embedded content.
	tmpl, err := template.New("otpTemplate").Parse(otpTemplate)
	if err != nil {
		return err
	}

	// Replace the placeholder with the actual verification link.
	data := OtpData{
		Code:       otpCode,
		ExpireTime: 10,
	}

	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, data); err != nil {
		return err
	}

	// create the message
	m := gomail.NewMessage()
	m.SetHeader("From", "hudyusufatsigah@gmail.com")
	m.SetHeader("To", to...)
	m.SetHeader("Content-Type", "text/html; charset=UTF-8")
	m.SetHeader("Subject", "OTP code")
	m.SetBody("text/html", emailContent.String())

	// create dialer to send message
	d := gomail.NewDialer("smtp.gmail.com", 587, "hudyusufatsigah@gmail.com", config.GlobalConfig.GoogleAppPassword)

	// send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
