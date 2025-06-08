package ruangan

import (
	"fmt"
	"net/http"
	"simpadu/helper"
	"simpadu/structs"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InDB struct {
	DB *gorm.DB
}

// CreateRuangan creates a new room
func (Idb *InDB) CreateRuangan(c *gin.Context) {
	var input struct {
		Gedung       string `json:"gedung" binding:"required"`       // A, B, C, dll
		Tipe         string `json:"tipe" binding:"required"`         // LAB, KLS, STU, dll
		Nomor        int    `json:"nomor" binding:"required"`        // 1, 2, 3, dll
		Nama_ruangan string `json:"nama_ruangan" binding:"required"` // Nama lengkap ruangan
		Kode_prodi   string `json:"kode_prodi" binding:"required"`
		Kapasitas    int    `json:"kapasitas" binding:"required"`
		Lokasi       string `json:"lokasi" binding:"required"` // Lantai 1, Lantai 2, dll
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Validasi token
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Validasi role (admin prodi atau admin akademik)
	if claims["role"] != "admin_prodi" && claims["role"] != "admin_akademik" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	// Generate kode ruangan
	kodeRuangan := fmt.Sprintf("%s-%s-%d", input.Gedung, input.Tipe, input.Nomor)

	// Validasi: cek apakah kode ruangan sudah ada
	var existingRuangan structs.Kode_ruangan
	if err := Idb.DB.Where("kode_ruangan = ?", kodeRuangan).First(&existingRuangan).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode ruangan sudah digunakan"})
		return
	}

	// Simpan ruangan
	ruangan := structs.Kode_ruangan{
		Kode_ruangan: kodeRuangan,
		Nama_ruangan: input.Nama_ruangan,
		Kode_prodi:   input.Kode_prodi,
		Gedung:       input.Gedung,
		Tipe:         input.Tipe,
		Nomor:        input.Nomor,
		Kapasitas:    input.Kapasitas,
		Lokasi:       input.Lokasi,
		Status:       "Aktif",
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := Idb.DB.Create(&ruangan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ruangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ruangan berhasil dibuat",
		"data":    ruangan,
	})
}

// GetRuanganByGedung gets all rooms in a specific building
func (Idb *InDB) GetRuanganByGedung(c *gin.Context) {
	gedung := c.Param("gedung")

	var ruangans []structs.Kode_ruangan
	if err := Idb.DB.Where("gedung = ?", gedung).
		Order("tipe, nomor").
		Find(&ruangans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ruangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ruangans,
	})
}

// GetRuanganByProdi gets all rooms for a specific program
func (Idb *InDB) GetRuanganByProdi(c *gin.Context) {
	kodeProdi := c.Param("kode_prodi")

	var ruangans []structs.Kode_ruangan
	if err := Idb.DB.Where("kode_prodi = ?", kodeProdi).
		Order("gedung, tipe, nomor").
		Find(&ruangans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ruangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ruangans,
	})
}

// GetRuanganByTipe gets rooms by type
func (Idb *InDB) GetRuanganByTipe(c *gin.Context) {
	tipe := c.Param("tipe")

	var ruangans []structs.Kode_ruangan
	if err := Idb.DB.Where("tipe = ?", tipe).
		Order("gedung, nomor").
		Find(&ruangans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ruangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ruangans,
	})
}

// UpdateRuangan updates room information
func (Idb *InDB) UpdateRuangan(c *gin.Context) {
	kodeRuangan := c.Param("kode_ruangan")
	var input struct {
		Nama_ruangan string `json:"nama_ruangan"`
		Kapasitas    int    `json:"kapasitas"`
		Lokasi       string `json:"lokasi"`
		Status       string `json:"status"` // "Aktif", "Maintenance", "Tidak Aktif"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Validasi token
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Validasi role
	if claims["role"] != "admin_prodi" && claims["role"] != "admin_akademik" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	// Update ruangan
	updates := map[string]interface{}{
		"nama_ruangan": input.Nama_ruangan,
		"kapasitas":    input.Kapasitas,
		"lokasi":       input.Lokasi,
		"status":       input.Status,
		"updated_at":   time.Now(),
	}

	if err := Idb.DB.Model(&structs.Kode_ruangan{}).
		Where("kode_ruangan = ?", kodeRuangan).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate ruangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ruangan berhasil diupdate",
	})
}

// DeleteRuangan deletes a room
func (Idb *InDB) DeleteRuangan(c *gin.Context) {
	kodeRuangan := c.Param("kode_ruangan")

	// Validasi token
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Validasi role
	if claims["role"] != "admin_prodi" && claims["role"] != "admin_akademik" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	// Cek apakah ruangan sedang digunakan dalam jadwal
	var jadwalCount int64
	Idb.DB.Model(&structs.Jadwal_matakuliah{}).
		Where("kode_ruangan = ?", kodeRuangan).
		Count(&jadwalCount)

	if jadwalCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ruangan tidak dapat dihapus karena masih digunakan dalam jadwal"})
		return
	}

	// Hapus ruangan
	if err := Idb.DB.Where("kode_ruangan = ?", kodeRuangan).
		Delete(&structs.Kode_ruangan{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus ruangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ruangan berhasil dihapus",
	})
}

// GetRuanganStatus gets room status and schedule
func (Idb *InDB) GetRuanganStatus(c *gin.Context) {
	kodeRuangan := c.Param("kode_ruangan")
	hari := c.Query("hari")

	var ruangan structs.Kode_ruangan
	if err := Idb.DB.Where("kode_ruangan = ?", kodeRuangan).First(&ruangan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruangan tidak ditemukan"})
		return
	}

	// Ambil jadwal ruangan untuk hari ini
	var jadwals []structs.Jadwal_matakuliah
	query := Idb.DB.Where("kode_ruangan = ?", kodeRuangan)
	if hari != "" {
		query = query.Where("hari = ?", hari)
	}
	if err := query.Order("jam_mulai").Find(&jadwals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil jadwal ruangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ruangan": ruangan,
		"jadwal":  jadwals,
	})
}
