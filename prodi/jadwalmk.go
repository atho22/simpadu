package prodi

import (
	"net/http"
	"simpadu/helper"
	"simpadu/structs"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (Idb *InDB) CreateJadwalMK(c *gin.Context) {
	var input_jadwal struct {
		Kode_matakuliah     string    `json:"kode_matakuliah" binding:"required"`
		Kelas_id            uint      `json:"kelas_id" binding:"required"`
		NIP                 string    `json:"NIP" binding:"required"`
		Hari                string    `json:"hari" binding:"required"`
		Kode_ruangan        string    `json:"kode_ruangan" binding:"required"`
		Kode_prodi          string    `json:"kode_prodi" binding:"required"`
		Jam_mulai           time.Time `json:"jam_mulai" binding:"required"`
		Kode_tahun_akademik string    `json:"kode_tahun_akademik" binding:"required"`
		Jam_selesai         time.Time `json:"jam_selesai" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input_jadwal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Request Tidak Valid"})
		return
	}
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	var kelas structs.Kelas
	if err := Idb.DB.Where("kode = ?", input_jadwal.Kelas_id).First(&kelas).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas tidak ditemukan"})
		return
	}
	var matakuliah structs.Mata_kuliah
	if err := Idb.DB.Where("kode_matakuliah = ?", input_jadwal.Kode_matakuliah).First(&matakuliah).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Matakuliah tidak ditemukan"})
		return
	}
	var ruangan structs.Kode_ruangan
	if err := Idb.DB.Where("kode_ruangan = ?", input_jadwal.Kode_ruangan).First(&ruangan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruangan tidak ditemukan"})
		return
	}
	var dosen structs.Dosen
	if err := Idb.DB.Where("NIP = ?", input_jadwal.NIP).First(&dosen).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosen tidak ditemukan"})
		return
	}
	var tahun_akademik structs.Tahun_akademik
	if err := Idb.DB.Where("kode_tahun_akademik = ?", kelas.Kode_tahun_akademik).First(&tahun_akademik).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tahun Akademik tidak ditemukan"})
		return
	}
	jadwal := structs.Jadwal_matakuliah{
		Kode_matakuliah:     input_jadwal.Kode_matakuliah,
		Kelas_id:            kelas.ID,
		Kode_dosen:          dosen.Kode_dosen,
		Kode_ruangan:        input_jadwal.Kode_ruangan,
		Kode_prodi:          kelas.Kode_prodi,
		Jam_mulai:           input_jadwal.Jam_mulai,
		Jam_selesai:         input_jadwal.Jam_selesai,
		Hari:                input_jadwal.Hari,
		Kode_tahun_akademik: tahun_akademik.Kode_tahun_akademik,
		Created_at:          time.Now(),
		Created_by:          claims["user_id"].(string),
	}
	if err := Idb.DB.Create(&jadwal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat jadwal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Jadwal berhasil dibuat",
		"status":  http.StatusOK,
		"data":    jadwal,
	})
}
