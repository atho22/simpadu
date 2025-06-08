package jadwal

import (
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

// CreateJadwal creates a new course schedule
func (Idb *InDB) CreateJadwal(c *gin.Context) {
	var input struct {
		Kode_matakuliah     string    `json:"kode_matakuliah" binding:"required"`
		Kelas_id            uint      `json:"kelas_id" binding:"required"`
		Kode_dosen          string    `json:"kode_dosen" binding:"required"`
		Kode_ruangan        string    `json:"kode_ruangan" binding:"required"`
		Kode_prodi          string    `json:"kode_prodi" binding:"required"`
		Jam_mulai           time.Time `json:"jam_mulai" binding:"required"`
		Jam_selesai         time.Time `json:"jam_selesai" binding:"required"`
		Hari                string    `json:"hari" binding:"required"`
		Kode_tahun_akademik string    `json:"kode_tahun_akademik" binding:"required"`
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

	// Validasi jam mulai harus sebelum jam selesai
	if input.Jam_mulai.After(input.Jam_selesai) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jam mulai harus sebelum jam selesai"})
		return
	}

	// Validasi ruangan tersedia
	var existingJadwal structs.Jadwal_matakuliah
	if err := Idb.DB.Where(
		"kode_ruangan = ? AND hari = ? AND ((jam_mulai <= ? AND jam_selesai > ?) OR (jam_mulai < ? AND jam_selesai >= ?) OR (jam_mulai >= ? AND jam_selesai <= ?))",
		input.Kode_ruangan,
		input.Hari,
		input.Jam_mulai,
		input.Jam_mulai,
		input.Jam_selesai,
		input.Jam_selesai,
		input.Jam_mulai,
		input.Jam_selesai,
	).First(&existingJadwal).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ruangan sudah digunakan pada jadwal tersebut"})
		return
	}

	// Validasi dosen tidak bentrok
	if err := Idb.DB.Where(
		"kode_dosen = ? AND hari = ? AND ((jam_mulai <= ? AND jam_selesai > ?) OR (jam_mulai < ? AND jam_selesai >= ?) OR (jam_mulai >= ? AND jam_selesai <= ?))",
		input.Kode_dosen,
		input.Hari,
		input.Jam_mulai,
		input.Jam_mulai,
		input.Jam_selesai,
		input.Jam_selesai,
		input.Jam_mulai,
		input.Jam_selesai,
	).First(&existingJadwal).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosen sudah memiliki jadwal pada waktu tersebut"})
		return
	}

	// Validasi kelas tidak bentrok
	if err := Idb.DB.Where(
		"kelas_id = ? AND hari = ? AND ((jam_mulai <= ? AND jam_selesai > ?) OR (jam_mulai < ? AND jam_selesai >= ?) OR (jam_mulai >= ? AND jam_selesai <= ?))",
		input.Kelas_id,
		input.Hari,
		input.Jam_mulai,
		input.Jam_mulai,
		input.Jam_selesai,
		input.Jam_selesai,
		input.Jam_mulai,
		input.Jam_selesai,
	).First(&existingJadwal).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kelas sudah memiliki jadwal pada waktu tersebut"})
		return
	}

	// Simpan jadwal
	jadwal := structs.Jadwal_matakuliah{
		Kode_matakuliah:     input.Kode_matakuliah,
		Kelas_id:            input.Kelas_id,
		Kode_dosen:          input.Kode_dosen,
		Kode_ruangan:        input.Kode_ruangan,
		Kode_prodi:          input.Kode_prodi,
		Jam_mulai:           input.Jam_mulai,
		Jam_selesai:         input.Jam_selesai,
		Hari:                input.Hari,
		Kode_tahun_akademik: input.Kode_tahun_akademik,
		Created_at:          time.Now(),
		Created_by:          claims["user_id"].(string),
	}

	if err := Idb.DB.Create(&jadwal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jadwal berhasil dibuat",
		"data":    jadwal,
	})
}

// GetJadwalByKelas gets schedules for a specific class
func (Idb *InDB) GetJadwalByKelas(c *gin.Context) {
	kelasID := c.Param("kelas_id")

	var jadwals []structs.Jadwal_matakuliah
	if err := Idb.DB.Where("kelas_id = ?", kelasID).
		Order("hari, jam_mulai").
		Find(&jadwals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": jadwals,
	})
}

// GetJadwalByDosen gets schedules for a specific lecturer
func (Idb *InDB) GetJadwalByDosen(c *gin.Context) {
	kodeDosen := c.Param("kode_dosen")

	var jadwals []structs.Jadwal_matakuliah
	if err := Idb.DB.Where("kode_dosen = ?", kodeDosen).
		Order("hari, jam_mulai").
		Find(&jadwals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": jadwals,
	})
}

// UpdateJadwal updates a schedule
func (Idb *InDB) UpdateJadwal(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Kode_ruangan string    `json:"kode_ruangan"`
		Jam_mulai    time.Time `json:"jam_mulai"`
		Jam_selesai  time.Time `json:"jam_selesai"`
		Hari         string    `json:"hari"`
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

	// Ambil jadwal yang akan diupdate
	var jadwal structs.Jadwal_matakuliah
	if err := Idb.DB.First(&jadwal, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan"})
		return
	}

	// Validasi jam mulai harus sebelum jam selesai
	if input.Jam_mulai.After(input.Jam_selesai) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jam mulai harus sebelum jam selesai"})
		return
	}

	// Validasi ruangan tersedia (kecuali jika ruangan tidak berubah)
	if input.Kode_ruangan != jadwal.Kode_ruangan {
		var existingJadwal structs.Jadwal_matakuliah
		if err := Idb.DB.Where(
			"kode_ruangan = ? AND hari = ? AND id != ? AND ((jam_mulai <= ? AND jam_selesai > ?) OR (jam_mulai < ? AND jam_selesai >= ?) OR (jam_mulai >= ? AND jam_selesai <= ?))",
			input.Kode_ruangan,
			input.Hari,
			id,
			input.Jam_mulai,
			input.Jam_mulai,
			input.Jam_selesai,
			input.Jam_selesai,
			input.Jam_mulai,
			input.Jam_selesai,
		).First(&existingJadwal).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ruangan sudah digunakan pada jadwal tersebut"})
			return
		}
	}

	// Validasi dosen tidak bentrok
	var existingJadwal structs.Jadwal_matakuliah
	if err := Idb.DB.Where(
		"kode_dosen = ? AND hari = ? AND id != ? AND ((jam_mulai <= ? AND jam_selesai > ?) OR (jam_mulai < ? AND jam_selesai >= ?) OR (jam_mulai >= ? AND jam_selesai <= ?))",
		jadwal.Kode_dosen,
		input.Hari,
		id,
		input.Jam_mulai,
		input.Jam_mulai,
		input.Jam_selesai,
		input.Jam_selesai,
		input.Jam_mulai,
		input.Jam_selesai,
	).First(&existingJadwal).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosen sudah memiliki jadwal pada waktu tersebut"})
		return
	}

	// Validasi kelas tidak bentrok
	if err := Idb.DB.Where(
		"kelas_id = ? AND hari = ? AND id != ? AND ((jam_mulai <= ? AND jam_selesai > ?) OR (jam_mulai < ? AND jam_selesai >= ?) OR (jam_mulai >= ? AND jam_selesai <= ?))",
		jadwal.Kelas_id,
		input.Hari,
		id,
		input.Jam_mulai,
		input.Jam_mulai,
		input.Jam_selesai,
		input.Jam_selesai,
		input.Jam_mulai,
		input.Jam_selesai,
	).First(&existingJadwal).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kelas sudah memiliki jadwal pada waktu tersebut"})
		return
	}

	// Update jadwal
	updates := map[string]interface{}{
		"kode_ruangan": input.Kode_ruangan,
		"jam_mulai":    input.Jam_mulai,
		"jam_selesai":  input.Jam_selesai,
		"hari":         input.Hari,
	}

	if err := Idb.DB.Model(&structs.Jadwal_matakuliah{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jadwal berhasil diupdate",
	})
}

// DeleteJadwal deletes a schedule
func (Idb *InDB) DeleteJadwal(c *gin.Context) {
	id := c.Param("id")

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

	// Hapus jadwal
	if err := Idb.DB.Delete(&structs.Jadwal_matakuliah{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jadwal berhasil dihapus",
	})
}

// GetRuanganTersedia gets available rooms for a specific time slot
func (Idb *InDB) GetRuanganTersedia(c *gin.Context) {
	var input struct {
		Hari        string    `json:"hari" binding:"required"`
		Jam_mulai   time.Time `json:"jam_mulai" binding:"required"`
		Jam_selesai time.Time `json:"jam_selesai" binding:"required"`
		Kode_prodi  string    `json:"kode_prodi" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Ambil semua ruangan yang tersedia
	var ruangans []structs.Kode_ruangan
	if err := Idb.DB.Where("kode_prodi = ? AND status = 'Aktif'", input.Kode_prodi).
		Find(&ruangans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ruangan"})
		return
	}

	// Filter ruangan yang sudah digunakan
	var ruanganTersedia []structs.Kode_ruangan
	for _, ruangan := range ruangans {
		var existingJadwal structs.Jadwal_matakuliah
		if err := Idb.DB.Where(
			"kode_ruangan = ? AND hari = ? AND ((jam_mulai <= ? AND jam_selesai > ?) OR (jam_mulai < ? AND jam_selesai >= ?) OR (jam_mulai >= ? AND jam_selesai <= ?))",
			ruangan.Kode_ruangan,
			input.Hari,
			input.Jam_mulai,
			input.Jam_mulai,
			input.Jam_selesai,
			input.Jam_selesai,
			input.Jam_mulai,
			input.Jam_selesai,
		).First(&existingJadwal).Error; err != nil {
			// Jika tidak ada jadwal yang bentrok, ruangan tersedia
			ruanganTersedia = append(ruanganTersedia, ruangan)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ruanganTersedia,
	})
}
