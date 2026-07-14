package auth

import (
	"api-presensi/internal/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// var secretKey = []byte("secret") // metode lama sebelum pakai .env

func getSecretKey() []byte {
	return []byte(config.App.JWTSecret) // output harus berupa []byte, lihat metode lama di atas
}

// isi: func GenerateToken dan ValidateToken
func GenerateToken(username string, role string) (string, error) {
	// convert var .env config.app.jwtexpiredHour ke int
	hour, err := strconv.Atoi(config.App.JWTExpireHour)
	if err != nil { // jika terjadi error, misalnya di .env value dari jwtexpiredhour nya=abc, jadikan default 1
		hour = 1
	}

	//create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * time.Duration(hour)).Unix(), // harus exp bukan expired, token expired dalam 1 jam
	})
	return token.SignedString(getSecretKey())
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return getSecretKey(), nil
	})
}
