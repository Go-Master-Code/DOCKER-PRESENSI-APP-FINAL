package database

import (
	"api-presensi/internal/config"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB // global var

func InitDB() {
	// ambil nilai host, port, user, password, dan db name dari config/config
	host := config.App.DBHost
	port := config.App.DBPort
	user := config.App.DBUser
	password := config.App.DBPassword
	dbname := config.App.DBName

	// bangun dsn dinamis
	dsn := fmt.Sprintf( // Asia%%2FJakarta => karena fmt.Sprintf() memakai % sebagai placeholder, jadi % harus di-escape menjadi %%.
		// "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FJakarta",
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)

	fmt.Println("DB_HOST =", host)
	fmt.Println("DB_PORT =", port)
	fmt.Println("DB_NAME =", dbname)
	fmt.Println("DSN =", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "/r/n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second, // tampilkan warning untuk query lambat
				LogLevel:      logger.Info, // tingkat log: Silent | Error | Warn | Info
				Colorful:      true,        // log berwarna
			},
		),
	})

	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	DB = db // var global agar bisa diakses dari repository
	log.Println("Berhasil terhubung ke database MySQL!")
}
