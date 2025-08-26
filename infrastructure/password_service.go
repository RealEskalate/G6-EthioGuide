package infrastructure

import (
	"EthioGuide/domain"

	"golang.org/x/crypto/bcrypt"
)

// bcryptService is the concrete implementation of PasswordService using the bcrypt algorithm.
// The struct is empty because the service is stateless.

type bcryptService struct{}

// NewPasswordService creates a new instance of our password service.
func NewPasswordService() domain.IPasswordService {
	return &bcryptService{}
}

// HashPassword generates a secure bcrypt hash of the password.
func (s *bcryptService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ComparePassword securely compares a hash with a plain-text password.
// It returns nil on success or an error if they don't match.
func (s *bcryptService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
