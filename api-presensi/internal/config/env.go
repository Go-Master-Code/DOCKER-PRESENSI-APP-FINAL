package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load() // tugasnya hanya membaca file .env ada / tidak
	if err != nil {
		log.Println("Menjalankan aplikasi tanpa file .env")
		return
	}

	log.Println(".env berhasil dimuat")
}
