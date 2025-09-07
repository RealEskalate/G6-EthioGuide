package infrastructure

import (
	"EthioGuide/domain"
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

// dialer interface allows mocking the gomail.Dialer
type dialer interface {
	DialAndSend(...*gomail.Message) error
}

type SmtpEmailService struct {
	host             string
	port             int
	username         string
	password         string
	from             string
	dialer           dialer
	verificationUrl  string
	resetPasswordUrl string
}

// emailData is a struct to hold the dynamic data for our HTML templates.
type emailData struct {
	Username  string
	ActionURL string
}

func NewSMTPEmailService(host string, port int, username, password, from, verificationUrl, resetPasswordUrl string) domain.IEmailService {
	d := gomail.NewDialer(host, port, username, password)

	return &SmtpEmailService{
		host:             host,
		port:             port,
		username:         username,
		password:         password,
		from:             from,
		dialer:           d,
		verificationUrl:  verificationUrl,
		resetPasswordUrl: resetPasswordUrl,
	}
}

func (s *SmtpEmailService) SendPasswordResetEmail(toEmail, username, resetToken string) error {
	subject := "Reset Your EthioGuide Password"
	actionURL := fmt.Sprintf("%s?token=%s", s.resetPasswordUrl, resetToken)

	data := emailData{
		Username:  username,
		ActionURL: actionURL,
	}

	// Specify the template file to use
	templateFile := "templates/password_reset.html"

	return s.send(toEmail, subject, templateFile, data)
}

func (s *SmtpEmailService) SendVerificationEmail(toEmail, username, activationToken string) error {
	subject := "Welcome to EthioGuide! Please Verify Your Account"
	actionURL := fmt.Sprintf("%s?token=%s", s.verificationUrl, activationToken)

	data := emailData{
		Username:  username,
		ActionURL: actionURL,
	}

	// Specify the template file to use
	templateFile := "templates/verification.html"

	return s.send(toEmail, subject, templateFile, data)
}

// send now parses an HTML template and executes it with the provided data.
func (s *SmtpEmailService) send(to, subject, templateFile string, data interface{}) error {
	m := gomail.NewMessage()

	// Parse the HTML template
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("could not parse email template %s: %w", templateFile, err)
	}

	// Execute the template with the data and write to a buffer
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("could not execute email template %s: %w", templateFile, err)
	}

	// Set headers and HTML body
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String(), gomail.SetPartEncoding(gomail.Unencoded))

	// Best Practice: Add a plain text alternative for clients that don't render HTML
	if d, ok := data.(emailData); ok {
		plainTextBody := fmt.Sprintf(
			"Hi %s,\n\nPlease use the following link to complete the action: %s\n\nThank you,\nThe EthioGuide Team",
			d.Username,
			d.ActionURL,
		)
		m.AddAlternative("text/plain", plainTextBody, gomail.SetPartEncoding(gomail.Unencoded))
	}

	return s.dialer.DialAndSend(m)
}
