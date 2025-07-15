package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/pareraamas/user_mgmt/service"
)

// Controller untuk TOTP (2FA)
// Komentar Bahasa Indonesia
type TOTPController struct {
	UserService *service.UserService
	TOTPService *service.TOTPService
}

func NewTOTPController(userService *service.UserService, totpService *service.TOTPService) *TOTPController {
	return &TOTPController{
		UserService: userService,
		TOTPService: totpService,
	}
}

// Endpoint untuk setup TOTP (generate secret & QR code)
func (ctl *TOTPController) SetupTOTP(c *fiber.Ctx) error {
	// Ambil user dari JWT (implementasi sederhana, asumsikan user_id di context)
	userID := c.Locals("user_id").(string)
	user, err := ctl.UserService.GetByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	if user.TOTPEnabled {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP sudah aktif"})
	}
	secret, err := ctl.TOTPService.GenerateSecret()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate secret"})
	}
	qrURL, err := ctl.TOTPService.GenerateQRCodeURL(user.Email, secret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate QR code"})
	}
	// Simpan secret sementara (belum enable sampai diverifikasi)
	user.TOTPSecret = secret
	ctl.UserService.Update(user)
	return c.JSON(fiber.Map{"secret": secret, "otpauth_url": qrURL})
}

// Endpoint untuk verifikasi kode TOTP pertama kali (enable 2FA)
func (ctl *TOTPController) EnableTOTP(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	user, err := ctl.UserService.GetByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	if user.TOTPEnabled {
		return c.Status(400).JSON(fiber.Map{"error": "TOTP sudah aktif"})
	}
	var body struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(c.Body(), &body); err != nil || body.Code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Kode TOTP wajib diisi"})
	}
	if ctl.TOTPService.VerifyCode(user.TOTPSecret, body.Code) {
		user.TOTPEnabled = true
		ctl.UserService.Update(user)
		return c.JSON(fiber.Map{"message": "TOTP berhasil diaktifkan"})
	}
	return c.Status(400).JSON(fiber.Map{"error": "Kode TOTP salah"})
}
