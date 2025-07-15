package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pareraamas/user_mgmt/controller"
)

// Fungsi untuk setup routing user
func SetupRoutes(app *fiber.App, userController *controller.UserController, totpController *controller.TOTPController, jwtMiddleware fiber.Handler) {
	app.Post("/register", userController.Register)
	app.Post("/login", userController.Login)
	app.Post("/request-verification", userController.RequestVerification)
	app.Post("/verify-email", userController.VerifyEmail)
	app.Get("/profile", userController.Profile)
	app.Post("/request-password-recovery", userController.RequestPasswordRecovery)
	app.Post("/reset-password", userController.ResetPassword)
	app.Post("/logout", userController.Logout)
	app.Put("/profile", userController.UpdateProfile)
	app.Get("/profile/custom", userController.GetCustomAttributes)
	app.Put("/profile/custom", userController.UpdateCustomAttributes)

	// Endpoint login TOTP (2FA step kedua)
	app.Post("/login/totp", userController.LoginTOTP)

	// Endpoint TOTP (2FA), hanya untuk user login
	totp := app.Group("/totp", jwtMiddleware)
	totp.Post("/setup", totpController.SetupTOTP)
	totp.Post("/enable", totpController.EnableTOTP)
}
