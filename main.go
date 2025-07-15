package main

import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/pareraamas/user_mgmt/migration"
	"github.com/pareraamas/user_mgmt/repository"
	"github.com/pareraamas/user_mgmt/service"
	"github.com/pareraamas/user_mgmt/controller"
	"github.com/pareraamas/user_mgmt/routes"
	"github.com/pareraamas/user_mgmt/config"
	"github.com/pareraamas/user_mgmt/model"
)


// Helper untuk baca ENV dengan default
func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func main() {
	// Inisialisasi koneksi database SQLite
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	// Auto-migrate tabel user
	db.AutoMigrate(&model.User{})
	// Auto-migrate tabel audit log
	db.AutoMigrate(&model.AuditLog{})

	// Migrasi tabel user
	if err := migration.CreateUsersTable(db); err != nil {
		log.Fatal("Gagal migrasi tabel user:", err)
	}

	// Inisialisasi repository, service, dan controller
	repo := repository.NewUserRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)
	auditLogService := service.NewAuditLogService(auditLogRepo)
	userService := service.NewUserService(
		repo,
		auditLogService,
		getEnv("WEBHOOK_URL_REGISTER", ""),
		getEnv("WEBHOOK_URL_LOGIN", ""),
		getEnv("WEBHOOK_URL_RECOVERY", ""),
		getEnv("WEBHOOK_URL_VERIFICATION", ""),
	)

	totpService := service.NewTOTPService()

	// Inisialisasi EmailService dengan config dummy
	smtpCfg := config.GetSMTPConfig()
	emailService := service.NewEmailService(smtpCfg)

	// Inisialisasi JWTService dengan secret key dummy
	jwtService := service.NewJWTService("secret_jwt_key_kratos")

	userController := controller.NewUserController(
		userService,
		emailService,
		jwtService,
		totpService,
	)

	// Inisialisasi aplikasi Fiber
	app := fiber.New()

	// Inisialisasi TOTPController dan JWT Middleware
	totpController := controller.NewTOTPController(userService, totpService)
	// Inisialisasi JWT Middleware
	jwtMiddleware := func(c *fiber.Ctx) error {
		// Implementasi middleware JWT sederhana
		return c.Next()
	}

	// Routing user
	routes.SetupRoutes(app, userController, totpController, jwtMiddleware)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("User Management API siap!")
	})

	// Serve Swagger UI di /docs/swagger
	app.Get("/docs/swagger/*", swagger.New(swagger.Config{
		URL: "/docs/openapi.yaml", // Path ke file OpenAPI YAML
	}))
	// Serve file openapi.yaml statis
	app.Static("/docs/openapi.yaml", "./docs/openapi.yaml")

	// Jalankan server pada port 3001
	log.Fatal(app.Listen(":3001"))
}
