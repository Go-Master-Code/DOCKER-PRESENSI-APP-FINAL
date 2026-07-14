package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	AppPort     string
	FrontendURL string

	JWTSecret string

	BackupPath string

	AppEnv        string
	GinMode       string
	JWTExpireHour string

	UploadPath string
}

var App Config // instance config

func LoadConfig() {
	App = Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		AppPort:       os.Getenv("APP_PORT"),
		FrontendURL:   os.Getenv("FRONTEND_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		BackupPath:    os.Getenv("BACKUP_PATH"),
		AppEnv:        os.Getenv("APP_ENV"),
		GinMode:       os.Getenv("GIN_MODE"),
		JWTExpireHour: os.Getenv("JWT_EXPIRE_HOURS"),
		UploadPath:    os.Getenv("UPLOAD_PATH"),
	}
	// eksekusi method untuk validasi masing-masing field App di atas
	validateConfig()
}

func validateConfig() {
	// buat var map untuk diiterasi dan validasi apakah sudah ada isinya belum tiap var nya
	requiredConfigs := map[string]string{
		"DB_HOST": App.DBHost,
		"DB_PORT": App.DBHost,
		"DB_USER": App.DBHost,
		//"DB_PASSWORD":      App.DBHost, ga perlu validasi, bisa aja password kosong saat local development
		"DB_NAME":          App.DBHost,
		"APP_PORT":         App.DBHost,
		"FRONTEND_URL":     App.DBHost,
		"JWT_SECRET":       App.DBHost,
		"BACKUP_PATH":      App.DBHost,
		"APP_ENV":          App.DBHost,
		"GIN_MODE":         App.DBHost,
		"JWT_EXPIRE_HOURS": App.DBHost,
	}

	for key, value := range requiredConfigs {
		if value == "" {
			log.Fatalf("Konfigurasi %s belum diisi pada file .env", key)
		}
	}
	// validasi app env dan gin mode agar tidak terjadi typo
	if App.AppEnv != "development" &&
		App.AppEnv != "production" {
		log.Fatal("APP_ENV harus development atau production")
	}

	if App.GinMode != "debug" &&
		App.GinMode != "release" {
		log.Fatal("GIN_MODE harus debug atau release")
	}
}
