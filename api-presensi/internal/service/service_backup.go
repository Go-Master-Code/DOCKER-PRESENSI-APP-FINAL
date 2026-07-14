package service

import (
	"api-presensi/internal/config"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// BACKUP SERVICE
// Bertanggung jawab membuat file backup database
type BackupService struct {
	backupPath string
	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPass     string
}

// constructor
// Membuat instance BackupService
func NewBackupService() *BackupService {
	return &BackupService{
		backupPath: config.App.BackupPath,
		dbPort:     config.App.DBPort,
		dbHost:     config.App.DBHost,
		dbName:     config.App.DBName,
		dbUser:     config.App.DBUser,
		dbPass:     config.App.DBPassword,
	}
}

// struct backup result untuk digunakan sebagai nilai balik dari service
type BackupResult struct {
	// nama file backup
	FileName string
	// lokasi lengkap file backup
	FilePath string
}

// Helper validasi service mysql
// memastikan seluruh konfigurasi penting sudah terisi
func (s *BackupService) validateConfig() error {
	// DB Host
	if strings.TrimSpace(s.dbHost) == "" {
		return fmt.Errorf("DB_HOST belum dikonfigurasi")
	}

	// DB Port
	if strings.TrimSpace(s.dbPort) == "" {
		return fmt.Errorf("DB_PORT belum dikonfigurasi")
	}

	// DB Name
	if strings.TrimSpace(s.dbName) == "" {
		return fmt.Errorf("DB_NAME belum dikonfigurasi")
	}

	// DB User
	if strings.TrimSpace(s.dbUser) == "" {
		return fmt.Errorf("DB_USER belum dikonfigurasi")
	}

	// DB Password
	if strings.TrimSpace(s.dbPass) == "" {
		return fmt.Errorf("DB_PASSWORD belum dikonfigurasi")
	}

	return nil
}

// CREATE BACKUP
// Membuat file SQL ke folder backup
func (s *BackupService) CreateBackup() (*BackupResult, error) {
	// Validasi konfigurasi
	err := s.validateConfig()
	if err != nil {
		return nil, err
	}

	// =========================
	// BUAT NAMA FILE
	// =========================

	now := time.Now().Format("2006-01-02_15-04-05")

	fileName := fmt.Sprintf(
		"backup_%s.sql",
		now,
	)

	filePath := fmt.Sprintf(
		"%s/%s",
		s.backupPath,
		fileName,
	)

	// =========================
	// PASTIKAN FOLDER BACKUP ADA
	// =========================

	err = os.MkdirAll(
		s.backupPath,
		os.ModePerm,
	)

	if err != nil {
		return nil, err
	}

	// =========================
	// CREATE FILE
	// =========================
	outFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	defer outFile.Close()

	// =========================
	// RUN MYSQLDUMP
	// =========================

	// Siapkan command mysqldump
	cmd := exec.Command(
		"mysqldump",

		// Host database
		"-h", s.dbHost,

		// Port database
		"-P", s.dbPort,

		// Username database
		"-u", s.dbUser,

		// Password database
		fmt.Sprintf("-p%s", s.dbPass),

		// Nama database yang akan dibackup
		s.dbName,
	)

	// Output SQL ditulis ke file backup
	cmd.Stdout = outFile

	// Tangkap hanya stderr dari mysqldump
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Jalankan mysqldump
	err = cmd.Run()

	if err != nil {
		return nil, fmt.Errorf(
			"mysqldump gagal: %v\n%s",
			err,
			stderr.String(),
		)
	}

	// -------------------------------------------------
	// Kembalikan informasi file backup
	// -------------------------------------------------

	return &BackupResult{
		FileName: fileName,
		FilePath: filePath,
	}, nil
}
