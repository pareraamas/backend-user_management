package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pareraamas/user_mgmt/service"
	"github.com/golang-jwt/jwt/v5"
	"encoding/json"
)

// Controller untuk endpoint user (register, login, verifikasi email, JWT, recovery)
// Komentar dalam Bahasa Indonesia

type UserController struct {
	Service      *service.UserService
	EmailService *service.EmailService
	JWTService   *service.JWTService
	TOTPService  *service.TOTPService
}

func NewUserController(service *service.UserService, emailService *service.EmailService, jwtService *service.JWTService, totpService *service.TOTPService) *UserController {
	return &UserController{Service: service, EmailService: emailService, JWTService: jwtService, TOTPService: totpService}
}

func (ctl *UserController) Register(c *fiber.Ctx) error {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	resp := service.RegisterFlow(ctl.Service, body.Email, body.Password)
	if resp.State == "failed" {
		return c.Status(400).JSON(resp)
	}
	return c.Status(201).JSON(resp)
}

// Endpoint login user (JWT, modular flow)
func (ctl *UserController) Login(c *fiber.Ctx) error {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	resp := service.LoginFlow(ctl.Service, ctl.JWTService, body.Email, body.Password)
	if resp.State == "failed" {
		return c.Status(401).JSON(resp)
	}
	return c.JSON(resp)
}

// Endpoint protected: profile user (butuh JWT)
func (ctl *UserController) Profile(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
	}
	parsed, err := ctl.JWTService.ValidateToken(token)
	if err != nil || !parsed.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	claims := parsed.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	user, err := ctl.Service.Repo.FindByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	return c.JSON(user)
}

// Endpoint request password recovery (modular flow)
func (ctl *UserController) RequestPasswordRecovery(c *fiber.Ctx) error {
	type req struct {
		Email string `json:"email"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	resp := service.RecoveryFlow(ctl.Service, body.Email)
	if resp.State == "failed" {
		return c.Status(400).JSON(resp)
	}
	return c.JSON(resp)
}

// Endpoint GET custom attributes user (hanya user login/JWT)
func (ctl *UserController) GetCustomAttributes(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")
	if tokenStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
	}
	parsed, err := ctl.JWTService.ValidateToken(tokenStr)
	if err != nil || !parsed.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	userID := claims["user_id"].(string)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	user, err := ctl.Service.Repo.FindByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"custom_attributes": user.CustomAttributes})
}

// Endpoint PUT custom attributes user (hanya user login/JWT)
func (ctl *UserController) UpdateCustomAttributes(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")
	if tokenStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
	}
	parsed, err := ctl.JWTService.ValidateToken(tokenStr)
	if err != nil || !parsed.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	userID := claims["user_id"].(string)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	user, err := ctl.Service.Repo.FindByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	attr, err := json.Marshal(body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal encode atribut custom"})
	}
	user.CustomAttributes = attr
	if err := ctl.Service.Repo.Update(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update custom attributes"})
	}
	return c.JSON(fiber.Map{"custom_attributes": user.CustomAttributes})
}

// Endpoint verifikasi kode TOTP setelah login password
func (ctl *UserController) LoginTOTP(c *fiber.Ctx) error {
	var body struct {
		FlowID string `json:"flow_id"`
		Code  string `json:"code"`
	}
	if err := c.BodyParser(&body); err != nil || body.FlowID == "" || body.Code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}
	userID, err := ctl.Service.GetUserIDByTOTPFlow(body.FlowID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Flow tidak ditemukan atau kadaluarsa"})
	}
	user, err := ctl.Service.GetByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	if !ctl.TOTPService.VerifyCode(user.TOTPSecret, body.Code) {
		return c.Status(401).JSON(fiber.Map{"error": "Kode TOTP salah"})
	}
	// Generate JWT token
	token, err := ctl.JWTService.GenerateToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
	}
	return c.JSON(fiber.Map{"token": token})
}

// Endpoint update profile user (hanya email dan status, harus login/JWT)
func (ctl *UserController) UpdateProfile(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	tokenStr := c.Get("Authorization")
	if tokenStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak ditemukan"})
	}
	parsed, err := ctl.JWTService.ValidateToken(tokenStr)
	if err != nil || !parsed.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	userID := claims["user_id"].(string)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}
	type req struct {
		Email  string `json:"email"`
		Status string `json:"status"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	user, err := ctl.Service.UpdateProfile(userID, body.Email, body.Status)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// Endpoint logout (pseudo, hanya client-side hapus token)
func (ctl *UserController) Logout(c *fiber.Ctx) error {
	// Di Kratos, session dihapus dari DB. Di JWT, cukup client hapus token, atau blacklist jika ingin lebih aman.
	return c.JSON(fiber.Map{"message": "Logout berhasil, silakan hapus token di client."})
}

// Endpoint reset password
func (ctl *UserController) ResetPassword(c *fiber.Ctx) error {
	type req struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	ok, err := ctl.Service.ResetPasswordWithCode(body.Email, body.Code, body.NewPassword)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Kode recovery salah atau kadaluarsa"})
	}
	return c.JSON(fiber.Map{"message": "Password berhasil direset"})
}

// Endpoint request verifikasi email
func (ctl *UserController) RequestVerification(c *fiber.Ctx) error {
	type req struct {
		Email string `json:"email"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	code, err := ctl.Service.GenerateAndSaveVerificationCode(body.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"flow_id": "", "state": "failed", "error": err.Error(), "next_action": "input_email"})
	}
	return c.JSON(fiber.Map{
		"flow_id": "", // Untuk flow sederhana (stateless)
		"state": "challenge",
		"next_action": "input_code",
		"data": fiber.Map{
			"info": "Kode verifikasi dikirim ke email (simulasi)",
			"code": code, // hanya untuk simulasi/dev
		},
	})
}

// Endpoint verifikasi email (modular flow)
func (ctl *UserController) VerifyEmail(c *fiber.Ctx) error {
	type req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	resp := service.VerificationFlow(ctl.Service, body.Email, body.Code)
	if resp.State == "failed" {
		return c.Status(400).JSON(resp)
	}
	return c.JSON(resp)
}
