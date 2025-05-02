package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config содержит настройки приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
	Email    EmailConfig
	Security SecurityConfig
	Logging  LoggingConfig
}

// ServerConfig содержит настройки сервера
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig содержит настройки базы данных
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// CORSConfig содержит настройки CORS
type CORSConfig struct {
	AllowedOrigins []string
}

// EmailConfig содержит настройки электронной почты
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	MailFrom     string
	MailFromName string
	CompanyEmail string
}

// SecurityConfig содержит настройки безопасности
type SecurityConfig struct {
	APIRateLimit int
	JWTSecret    string
	EnableHTTPS  bool
}

// LoggingConfig содержит настройки логирования
type LoggingConfig struct {
	Level string
}

// LoadConfig загружает настройки из .env файла
func LoadConfig() (Config, error) {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load("../../.env")
		if err != nil {
			fmt.Println("Файл .env не найден, используются переменные окружения")
		}
	}

	config := Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "creality_workshop"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", ""),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			MailFrom:     getEnv("MAIL_FROM", "no-reply@example.com"),
			MailFromName: getEnv("MAIL_FROM_NAME", "Creality Workshop"),
			CompanyEmail: getEnv("COMPANY_EMAIL", "info@example.com"),
		},
		Security: SecurityConfig{
			APIRateLimit: getEnvAsInt("API_RATE_LIMIT", 100),
			JWTSecret:    getEnv("JWT_SECRET", "change_this_to_something_secure"),
			EnableHTTPS:  getEnv("ENABLE_HTTPS", "false") == "true",
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	return config, nil
}

// DSN возвращает строку подключения к базе данных
func (db *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
}

// Вспомогательные функции для получения переменных окружения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
