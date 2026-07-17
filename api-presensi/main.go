package main

import (
	"api-presensi/internal/config"
	"api-presensi/internal/database"
	"api-presensi/internal/handler"
	"api-presensi/internal/middleware"
	"api-presensi/internal/repository"
	"api-presensi/internal/routes"
	"api-presensi/internal/service"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// load file .env
	config.LoadEnv()
	// test sederhana, tambahkan import "os" dan test log di bawah ini
	// println("DB_HOST =", os.Getenv("DB_HOST"))

	// load config dari config/config.go
	config.LoadConfig()
	// saat awal boleh debug setiap struct config.App sbb:
	// fmt.Println("DB_HOST=", config.App.DBHost)           // test apakah config sudah ter load dengan benar
	// fmt.Println("APP_PORT=", config.App.AppPort)         // test apakah config (APP_PORT) sudah ter load
	// fmt.Println("FRONTEND_URL=", config.App.FrontendURL) // test apakah config (FRONTEND_URL) sudah ter load
	// fmt.Println("JWT_SECRET=", config.App.JWTSecret)     // test apakah config (JWT_SECRET) sudah ter load
	// fmt.Println("BACKUP_PATH=", config.App.BackupPath)   // test apakah config (BACKUP_PATH) sudah ter load
	// BACKEND, GIN MODE, JWT
	// fmt.Println("APP_ENV=", config.App.AppEnv)                 // test apakah config (APP_ENV) sudah ter load
	// fmt.Println("GIN_MODE=", config.App.GinMode)               // test apakah config (GIN_MODE) sudah ter load
	// fmt.Println("JWT_EXPIRE_HOURS=", config.App.JWTExpireHour) // test apakah config (JWT_EXPIRE_HOURS) sudah ter load

	// set gin mode
	gin.SetMode(config.App.GinMode) // mode di set sesuai file env

	// setting time.Now() ke loc Asia/Jakarta (WIB) menghindari UTC
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err == nil {
		time.Local = loc
	}

	// set locale Asia/Jakarta agar jam mengikuti WIB
	log.Println("Go Time =", time.Now())

	// inisialisasi database
	database.InitDB()

	// inisiasi engine gin
	r := gin.New()

	// don't trust all proxies
	r.SetTrustedProxies(nil)

	// tambahkan CORS apabila server backend berbeda dengan frontend
	// ===============================
	// 🔥 CORS CONFIG
	// ===============================
	r.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(config.App.FrontendURL, ","),

		// method yang diizinkan
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		// 🔥 HEADER YANG DIIZINKAN
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization", // WAJIB untuk JWT
		},

		// optional
		ExposeHeaders: []string{
			"Content-Length",
		},

		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}))
	// go get github.com/gin-contrib/cors

	// dependency injection karyawan (repository -> service -> handler)
	repoKaryawan := repository.NewRepositoryKaryawan(database.DB)
	serviceKaryawan := service.NewServiceKaryawan(repoKaryawan)
	handlerKaryawan := handler.NewHandlerKaryawan(serviceKaryawan)

	// dependency injection jenis ijin (repository -> service -> handler)
	repoJenisIjin := repository.NewRepositoryJenisIjin(database.DB)
	serviceJenisIjin := service.NewServiceJenisIjin(repoJenisIjin)
	handlerJenisIjin := handler.NewHandlerJenisIjin(serviceJenisIjin)

	// dependency injection hari libur (repository -> service -> handler)
	repoHariLibur := repository.NewRepositoryHariLibur(database.DB)
	serviceHariLibur := service.NewServiceHariLibur(repoHariLibur)
	handlerHariLibur := handler.NewHandlerHariLibur(serviceHariLibur)

	// dependency injection ijin karyawan
	repoIjinKaryawan := repository.NewRepositoryIjinKaryawan(database.DB)
	serviceIjinKaryawan := service.NewServiceIjinKaryawan(repoIjinKaryawan)
	handlerIjinKaryawan := handler.NewHandlerIjinKaryawan(serviceIjinKaryawan)

	// dependency injection presensi
	repoPresensi := repository.NewRepositoryPresensi(database.DB)
	servicePresensi := service.NewServicePresensi(repoPresensi, repoKaryawan)
	handlerPresensi := handler.NewHandlerPresensi(servicePresensi)

	// dependency injection roles
	repoRole := repository.NewRepositoryRole(database.DB)
	serviceRole := service.NewServiceRole(repoRole)
	handlerRole := handler.NewHandlerRole(serviceRole)

	// dependency injection user
	repoUser := repository.NewRepositoryUser(database.DB)
	serviceUser := service.NewServiceUser(repoUser)
	handlerUser := handler.NewHandlerUser(serviceUser)

	// dependency injection logger ke DB
	repoLog := repository.NewRepositoryLog(database.DB)
	serviceLog := service.NewServiceLog(repoLog)
	handlerLog := handler.NewHandlerLog(serviceLog)

	// dependency injection backup db
	serviceBackup := service.NewBackupService()
	handlerBackup := handler.NewHandlerBackup(serviceBackup)

	// ✅ Tambahkan middleware logging di sini
	r.Use(gin.Recovery()) // mencegah server mati, return http 500 jika server mati
	// RequestLogger -> log di console
	r.Use(middleware.RequestLogger())
	// middleware logger ke DB (hanya untuk POST, PUT, dan DELETE)
	r.Use(middleware.RequestLoggerDB(serviceLog))

	// group public untuk endpoint yang tidak perlu protection, misalnya login
	public := r.Group("/api")
	public.POST("/login", middleware.RateLimiter(), handlerUser.Login) // pakai middleware RateLimiter (batas 5x login per 1 menit)
	public.POST("/backup", handlerBackup.BackupDatabase)

	/*
		Fungsinya:

		membuat file .sql
		menjalankan mysqldump

		✔ Ini = CREATE backup
	*/

	// authorization yang akan dipasang pada tiap endpoint yang dilindungi (harus punya token)
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthRequired()) // file auth_required.go -> handler ini dieksekusi dulu sebeleum eksekusi handler endpoint
	{
		// list handler karyawan
		routes.RegisterRoutesKaryawan(authorized, handlerKaryawan)
		// list handler jenis ijin
		routes.RegisterRoutesJenisIjin(authorized, handlerJenisIjin)
		// list handler hari libur
		routes.RegisterRoutesHariLibur(authorized, handlerHariLibur)
		// list handler ijin karyawan
		routes.RegisterRoutesIjinKaryawan(authorized, handlerIjinKaryawan)
		// list handler presensi
		routes.RegisterRoutesPresensi(authorized, handlerPresensi)
		// list handler role
		routes.RegisterRoutesRole(authorized, handlerRole)
		// list handler user
		routes.RegisterRoutesUser(authorized, handlerUser)
		// list handler log (untuk ditampilkan di table untuk super admin)
		routes.RegisterRoutesLog(authorized, handlerLog)
	}

	// run server
	r.Run(":" + config.App.AppPort)
}
