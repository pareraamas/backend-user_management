package service

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// Service untuk TOTP (2FA)
// Komentar Bahasa Indonesia
type TOTPService struct{}

func NewTOTPService() *TOTPService {
	return &TOTPService{}
}

// Generate secret TOTP baru untuk user
func (s *TOTPService) GenerateSecret() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "KratosGoFiber",
		AccountName: "user@example.com",
	})
	if err != nil {
		return "", err
	}
	return key.Secret(), nil
}

// Generate kode QR otpauth URL untuk aplikasi authenticator
func (s *TOTPService) GenerateQRCodeURL(email, secret string) (string, error) {
	key, err := otp.NewKeyFromURL("otpauth://totp/KratosGoFiber:" + email + "?secret=" + secret + "&issuer=KratosGoFiber")
	if err != nil {
		return "", err
	}
	return key.URL(), nil
}

// Verifikasi kode TOTP dari user
func (s *TOTPService) VerifyCode(secret, code string) bool {
	return totp.Validate(code, secret)
}
