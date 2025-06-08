package absensi

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

// BukaPertemuan opens a new meeting for attendance
func (Idb *InDB) BukaPertemuan(c *gin.Context) {
	var input struct {
		Kode_matakuliah string `json:"kode_matakuliah" binding:"required"`
		Pertemuan_ke    int    `json:"pertemuan_ke" binding:"required"`
		Durasi          int    `json:"durasi" binding:"required"` // durasi dalam menit
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

	// Validasi role dosen
	if claims["role"] != "dosen" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya dosen yang dapat membuka pertemuan"})
		return
	}

	// Validasi: apakah sudah dibuka sebelumnya
	var existing structs.Pertemuan
	err := Idb.DB.Where("kode_matakuliah = ? AND pertemuan_ke = ?", input.Kode_matakuliah, input.Pertemuan_ke).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pertemuan ini sudah dibuka"})
		return
	}

	// Hitung waktu tutup
	waktuDibuka := time.Now()
	waktuDitutup := waktuDibuka.Add(time.Duration(input.Durasi) * time.Minute)

	// Generate kode pertemuan
	kodePertemuan := helper.GenerateKodePertemuan(input.Kode_matakuliah, input.Pertemuan_ke)

	// Simpan ke DB
	pertemuan := structs.Pertemuan{
		Kode_pertemuan:  kodePertemuan,
		Kode_matakuliah: input.Kode_matakuliah,
		Pertemuan_ke:    input.Pertemuan_ke,
		Dibuka:          true,
		Waktu_dibuka:    waktuDibuka,
		Waktu_ditutup:   waktuDitutup,
		Dibuka_by:       claims["user_id"].(string),
		Created_at:      time.Now(),
		Updated_at:      time.Now(),
	}

	if err := Idb.DB.Create(&pertemuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka pertemuan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pertemuan berhasil dibuka",
		"data": gin.H{
			"kode_pertemuan": kodePertemuan,
			"waktu_dibuka":   waktuDibuka,
			"waktu_ditutup":  waktuDitutup,
		},
	})
}

// TutupPertemuan closes an active meeting
func (Idb *InDB) TutupPertemuan(c *gin.Context) {
	kodePertemuan := c.Param("kode_pertemuan")

	// Validasi token
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Validasi role dosen
	if claims["role"] != "dosen" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya dosen yang dapat menutup pertemuan"})
		return
	}

	// Update status pertemuan
	result := Idb.DB.Model(&structs.Pertemuan{}).
		Where("kode_pertemuan = ? AND dibuka = ?", kodePertemuan, true).
		Updates(map[string]interface{}{
			"dibuka":     false,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menutup pertemuan"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan atau sudah ditutup"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pertemuan berhasil ditutup"})
}

// Absen marks student attendance
func (Idb *InDB) Absen(c *gin.Context) {
	var input struct {
		Kode_pertemuan string `json:"kode_pertemuan" binding:"required"`
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

	// Validasi role mahasiswa
	if claims["role"] != "mahasiswa" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya mahasiswa yang dapat melakukan absensi"})
		return
	}

	// Cek status pertemuan
	var pertemuan structs.Pertemuan
	if err := Idb.DB.Where("kode_pertemuan = ?", input.Kode_pertemuan).First(&pertemuan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	if !pertemuan.Dibuka {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pertemuan sudah ditutup"})
		return
	}

	if time.Now().After(pertemuan.Waktu_ditutup) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Waktu absensi sudah berakhir"})
		return
	}

	// Cek apakah sudah absen
	var existingAbsensi structs.Absensi
	if err := Idb.DB.Where("kode_pertemuan = ? AND nim = ?", input.Kode_pertemuan, claims["nim"]).First(&existingAbsensi).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan absensi"})
		return
	}

	// Simpan absensi
	absensi := structs.Absensi{
		Kode_pertemuan:  input.Kode_pertemuan,
		Kode_matakuliah: pertemuan.Kode_matakuliah,
		Nim:             claims["nim"].(string),
		Created_at:      time.Now(),
	}

	if err := Idb.DB.Create(&absensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Absensi berhasil"})
}

// GetStatusPertemuan gets the current status of a meeting
func (Idb *InDB) GetStatusPertemuan(c *gin.Context) {
	kodeMatakuliah := c.Param("kode_matakuliah")
	var pertemuan structs.Pertemuan

	// Ambil pertemuan terakhir yang aktif
	err := Idb.DB.Where("kode_matakuliah = ? AND dibuka = ?", kodeMatakuliah, true).
		Order("pertemuan_ke desc").
		First(&pertemuan).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":       "belum_dibuka",
			"pertemuan_ke": 0,
			"message":      "Belum ada pertemuan yang dibuka",
		})
		return
	}

	// Hitung sisa waktu menggunakan time.Until
	sisaWaktu := time.Until(pertemuan.Waktu_ditutup)
	if sisaWaktu < 0 {
		sisaWaktu = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "dibuka",
		"pertemuan_ke":   pertemuan.Pertemuan_ke,
		"kode_pertemuan": pertemuan.Kode_pertemuan,
		"sisa_waktu":     int(sisaWaktu.Seconds()),
		"waktu_dibuka":   pertemuan.Waktu_dibuka,
		"waktu_ditutup":  pertemuan.Waktu_ditutup,
	})
}

// GetRekapAbsensi gets attendance summary for a course
func (Idb *InDB) GetRekapAbsensi(c *gin.Context) {
	kodeMatakuliah := c.Param("kode_matakuliah")

	// Validasi token
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Hanya dosen yang dapat melihat rekap
	if claims["role"] != "dosen" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	// Ambil semua pertemuan untuk matakuliah ini
	var pertemuans []structs.Pertemuan
	if err := Idb.DB.Where("kode_matakuliah = ?", kodeMatakuliah).
		Order("pertemuan_ke asc").
		Find(&pertemuans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pertemuan"})
		return
	}

	// Ambil semua absensi
	var absensis []structs.Absensi
	if err := Idb.DB.Where("kode_matakuliah = ?", kodeMatakuliah).Find(&absensis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data absensi"})
		return
	}

	// Format response
	var rekap []gin.H
	for _, p := range pertemuans {
		// Hitung jumlah absen untuk pertemuan ini
		var jumlahAbsen int64
		Idb.DB.Model(&structs.Absensi{}).
			Where("kode_pertemuan = ?", p.Kode_pertemuan).
			Count(&jumlahAbsen)

		rekap = append(rekap, gin.H{
			"pertemuan_ke":   p.Pertemuan_ke,
			"kode_pertemuan": p.Kode_pertemuan,
			"status":         p.Dibuka,
			"waktu_dibuka":   p.Waktu_dibuka,
			"waktu_ditutup":  p.Waktu_ditutup,
			"jumlah_absen":   jumlahAbsen,
			"dibuka_oleh":    p.Dibuka_by,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"kode_matakuliah": kodeMatakuliah,
		"rekap":           rekap,
	})
}
