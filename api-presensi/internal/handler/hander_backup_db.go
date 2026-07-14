package handler

import (
	"api-presensi/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// no interface, langsung struct implementasi
type HandlerBackup struct {
	service *service.BackupService
}

// constructor
func NewHandlerBackup(service *service.BackupService) *HandlerBackup {
	return &HandlerBackup{
		service: service,
	}
}

// =====================================================
// BACKUP DATABASE
//
// Flow:
//
// Frontend
//      │
//      ▼
// POST /api/backup
//      │
//      ▼
// BackupService.CreateBackup()
//      │
//      ▼
// File SQL dibuat
//      │
//      ▼
// Browser langsung download
// =====================================================

// =========================
// 🔥 MANUAL BACKUP DATABASE
// =========================
func (h *HandlerBackup) BackupDatabase(c *gin.Context) {
	// jalankan proses backup database
	result, err := h.service.CreateBackup()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(), // tampilkan error asli
		})
		return
	}

	// kirim file sql ke browser, browser otomatis mengunduh file
	c.FileAttachment(
		result.FilePath,
		result.FileName,
	)
}
