package config

import (
	"github.com/joho/godotenv"
)

func LoadEnv() {
	// OLD METHOD sebelum implementasi ke docker
	// err := godotenv.Load() // tugasnya hanya membaca file .env ada / tidak
	// if err != nil {
	// 	log.Println("Menjalankan aplikasi tanpa file .env")
	// 	return
	// }

	// log.Println(".env berhasil dimuat")

	// Abaikan jika file .env tidak ada.
	// Environment variable dari Docker tetap akan digunakan.
	_ = godotenv.Load()
}
