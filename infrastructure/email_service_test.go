package infrastructure_test

import (
	. "EthioGuide/infrastructure"
	"errors"
	"fmt"
	"io"
	"mime/quotedprintable"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
)

// MockDialer is a mock for the gomail.Dialer interface
type MockDialer struct {
	mock.Mock
}

func (m *MockDialer) DialAndSend(msgs ...*gomail.Message) error {
	args := m.Called(msgs)
	return args.Error(0)
}

// EmailServiceTestSuite is the test suite for SmtpEmailService
type EmailServiceTestSuite struct {
	suite.Suite
	mockDialer       *MockDialer
	service          *SmtpEmailService
	verificationUrl  string
	resetPasswordUrl string
}

// SetupSuite runs once before the entire suite, creating temp directories/files.
func (s *EmailServiceTestSuite) SetupSuite() {
	// Create a temporary templates directory for our tests
	err := os.Mkdir("templates", 0755)
	s.Require().NoError(err)

	// Create dummy template files
	verificationTpl := `Hello {{.Username}}, please verify here: {{.ActionURL}}`
	resetTpl := `Hello {{.Username}}, please reset here: {{.ActionURL}}`
	err = os.WriteFile("templates/verification.html", []byte(verificationTpl), 0644)
	s.Require().NoError(err)
	err = os.WriteFile("templates/password_reset.html", []byte(resetTpl), 0644)
	s.Require().NoError(err)
}

// TearDownSuite runs once after the entire suite, cleaning up temp files.
func (s *EmailServiceTestSuite) TearDownSuite() {
	err := os.RemoveAll("templates")
	s.Require().NoError(err)
}

// SetupTest is run before each test in the suite
func (s *EmailServiceTestSuite) SetupTest() {
	s.mockDialer = new(MockDialer)
	s.verificationUrl = "http://localhost/verify"
	s.resetPasswordUrl = "http://localhost/reset"

	// We cast the result of NewSMTPEmailService to our concrete type to access the unexported dialer field
	// Note: The constructor returns an interface, so we need to cast it.
	service := NewSMTPEmailService(
		"smtp.example.com",
		587,
		"user",
		"pass",
		"from@example.com",
		s.verificationUrl,
		s.resetPasswordUrl,
	).(*SmtpEmailService)

	// Replace the real dialer with our mock
	service.SetDialer(s.mockDialer)
	s.service = service
}

// We need to add an exported method to SmtpEmailService to inject the mock dialer for testing.
// Add this method to your `email_service.go` file:
/*
func (s *SmtpEmailService) SetDialer(d dialer) {
    s.dialer = d
}
*/

func (s *EmailServiceTestSuite) TestSendPasswordResetEmail_Success() {
	// Arrange
	toEmail := "test@example.com"
	username := "TestUser"
	resetToken := "reset123"
	expectedURL := fmt.Sprintf("%s?token=%s", s.resetPasswordUrl, resetToken)

	s.mockDialer.On("DialAndSend", mock.MatchedBy(func(msgs []*gomail.Message) bool {
		s.Len(msgs, 1)
		msg := msgs[0]
		s.Equal(toEmail, msg.GetHeader("To")[0])
		s.Equal("Reset Your EthioGuide Password", msg.GetHeader("Subject")[0])

		// --- FINAL CORRECTED ASSERTION ---
		// Render the entire multipart message to a buffer.
		bodyBuffer := new(strings.Builder)
		_, err := msg.WriteTo(bodyBuffer)
		s.NoError(err)
		fullEmailBody := bodyBuffer.String()

		// Create a Quoted-Printable reader and decode the body.
		qpReader := quotedprintable.NewReader(strings.NewReader(fullEmailBody))
		decodedBody, err := io.ReadAll(qpReader)
		s.NoError(err, "Failed to decode quoted-printable body")

		// Now assert against the clean, decoded body.
		s.Contains(string(decodedBody), fmt.Sprintf("Hi %s", username))
		s.Contains(string(decodedBody), expectedURL)

		return true
	})).Return(nil).Once()

	// Act
	err := s.service.SendPasswordResetEmail(toEmail, username, resetToken)

	// Assert
	s.NoError(err)
	s.mockDialer.AssertExpectations(s.T())
}

func (s *EmailServiceTestSuite) TestSendVerificationEmail_Success() {
	// Arrange
	toEmail := "verify@example.com"
	username := "NewUser"
	activationToken := "activate123"
	expectedURL := fmt.Sprintf("%s?token=%s", s.verificationUrl, activationToken)

	s.mockDialer.On("DialAndSend", mock.MatchedBy(func(msgs []*gomail.Message) bool {
		s.Len(msgs, 1)
		msg := msgs[0]
		s.Equal(toEmail, msg.GetHeader("To")[0])
		s.Equal("Welcome to EthioGuide! Please Verify Your Account", msg.GetHeader("Subject")[0])

		// --- FINAL CORRECTED ASSERTION ---
		// Render and decode the body to get the clean content.
		bodyBuffer := new(strings.Builder)
		_, err := msg.WriteTo(bodyBuffer)
		s.NoError(err)
		fullEmailBody := bodyBuffer.String()

		qpReader := quotedprintable.NewReader(strings.NewReader(fullEmailBody))
		decodedBody, err := io.ReadAll(qpReader)
		s.NoError(err, "Failed to decode quoted-printable body")

		s.Contains(string(decodedBody), fmt.Sprintf("Hi %s", username))
		s.Contains(string(decodedBody), expectedURL)

		return true
	})).Return(nil).Once()

	// Act
	err := s.service.SendVerificationEmail(toEmail, username, activationToken)

	// Assert
	s.NoError(err)
	s.mockDialer.AssertExpectations(s.T())
}

func (s *EmailServiceTestSuite) TestSend_DialerError() {
	// Arrange
	expectedErr := errors.New("smtp connection failed")
	s.mockDialer.On("DialAndSend", mock.Anything).Return(expectedErr).Once()

	// Act
	err := s.service.SendVerificationEmail("test@example.com", "User", "token")

	// Assert
	s.Error(err)
	s.Equal(expectedErr, err)
	s.mockDialer.AssertExpectations(s.T())
}

func (s *EmailServiceTestSuite) TestSend_TemplateParsingError() {
	// Arrange: To simulate a parsing error, we can temporarily remove the template file.
	s.Require().NoError(os.Rename("templates/verification.html", "templates/verification.html.tmp"))
	defer os.Rename("templates/verification.html.tmp", "templates/verification.html") // Ensure cleanup

	// Act
	err := s.service.SendVerificationEmail("test@example.com", "User", "token")

	// Assert
	s.Error(err)
	s.Contains(err.Error(), "could not parse email template")
	// The dialer should not be called if parsing fails
	s.mockDialer.AssertNotCalled(s.T(), "DialAndSend")
}

// TestEmailServiceTestSuite runs the entire suite
func TestEmailServiceTestSuite(t *testing.T) {
	suite.Run(t, new(EmailServiceTestSuite))
}
