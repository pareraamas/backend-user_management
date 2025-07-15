package config

// Konfigurasi SMTP dummy untuk simulasi pengiriman email

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

func GetSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:     "smtp.dummy.local",
		Port:     587,
		Username: "dummy@dummy.local",
		Password: "dummy123",
		From:     "noreply@dummy.local",
		FromName: "Dummy Sender",
	}
}
