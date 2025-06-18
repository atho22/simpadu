package akademik

import (
	"net/http"
	"simpadu/helper"
	"simpadu/structs"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InDB struct {
	DB *gorm.DB
}

func (Idb *InDB) GetAllTahunAkademik(c *gin.Context) {
	var tahunAkademik []structs.Tahun_akademik
	if err := Idb.DB.Table("tahun_akademiks").Scan(&tahunAkademik).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Mengambil Data Tahun Akademik"})
		return
	}
	if len(tahunAkademik) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak Ada Data Tahun Akademik"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Berhasil Mengambil Data Tahun Akademik",
		"data":    tahunAkademik,
	})
}

func (Idb *InDB) GetTahunAkademikByID(c *gin.Context) {
	id := c.Param("id")
	var tahunAkademik structs.Tahun_akademik
	if err := Idb.DB.Table("tahun_akademiks").Where("id = ?", id).Scan(&tahunAkademik).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Mengambil Data Tahun Akademik"})
		return
	}
	if tahunAkademik.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tahun Akademik Tidak Ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil Mengambil Data Tahun Akademik",
		"data":    tahunAkademik,
	})
}

func (Idb *InDB) CreateTahunAkademik(c *gin.Context) {
	var (
		InputTahunAkademik structs.Tahun_akademik
	)
	if err := c.ShouldBindJSON(&InputTahunAkademik); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Request Tidak Valid"})
		return
	}

	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}
	tahun_akademik := structs.Tahun_akademik{
		Tahun:           InputTahunAkademik.Tahun,
		Semester:        InputTahunAkademik.Semester,
		Status:          InputTahunAkademik.Status,
		Tanggal_mulai:   InputTahunAkademik.Tanggal_mulai,
		Tanggal_selesai: InputTahunAkademik.Tanggal_selesai,
	}

	// Validasi input dulu
	if InputTahunAkademik.Tahun == "" || InputTahunAkademik.Semester == "" || InputTahunAkademik.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak boleh kosong"})
		return
	}

	if InputTahunAkademik.Tanggal_mulai.After(InputTahunAkademik.Tanggal_selesai) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal mulai tidak boleh lebih besar dari tanggal selesai"})
		return
	}

	// Baru insert ke database
	if err := Idb.DB.Table("tahun_akademiks").Create(&tahun_akademik).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Membuat Tahun Akademik"})
		return
	}

	// Jika berhasil
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Berhasil Membuat Tahun Akademik",
		"data":    tahun_akademik,
	})
	
}
func (Idb *InDB) UpdateTahunAkademik(c *gin.Context) {
	var (
		InputTahunAkademik structs.Tahun_akademik
	)
	if err := c.ShouldBindJSON(&InputTahunAkademik); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Request Tidak Valid"})
		return
	}

	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	tahun_akademik := structs.Tahun_akademik{
		Tahun:           InputTahunAkademik.Tahun,
		Semester:        InputTahunAkademik.Semester,
		Status:          InputTahunAkademik.Status,
		Tanggal_mulai:   InputTahunAkademik.Tanggal_mulai,
		Tanggal_selesai: InputTahunAkademik.Tanggal_selesai,
	}

	if err := Idb.DB.Table("tahun_akademiks").Where("id = ?", id).Updates(&tahun_akademik).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Mengupdate Tahun Akademik"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Berhasil Mengupdate Tahun Akademik",
		"data":    tahun_akademik,
	})
}

func (Idb *InDB) DeleteTahunAkademik(c *gin.Context) {
	id := c.Param("id")

	if err := Idb.DB.Table("tahun_akademiks").Where("id = ?", id).Delete(&structs.Tahun_akademik{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Menghapus Tahun Akademik"})
		return
	}

	// Update status kelas_mahasiswa
	if err := Idb.DB.Table("kelas_mahasiswas").Where("kode_tahun_akademik = ?", id).Update("kode_tahun_akademik", "").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Mengupdate Status Kelas Mahasiswa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Berhasil Menghapus Tahun Akademik",
	})
}

func (Idb *InDB) GetAcademicSummary(c *gin.Context) {
	var currentTahunAkademik structs.Tahun_akademik
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	if err := Idb.DB.Table("tahun_akademiks").Where("status = ?", "Aktif").First(&currentTahunAkademik).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil tahun akademik aktif: " + err.Error()})
		return
	}

	var totalKelas int64
	if err := Idb.DB.Table("kelas").Where("kode_tahun_akademik = ?", currentTahunAkademik.Kode_tahun_akademik).Count(&totalKelas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total kelas: " + err.Error()})
		return
	}

	var totalMahasiswa int64
	query := Idb.DB.Table("mahasiswas").
		Select("COUNT(DISTINCT mahasiswas.nim)").
		Joins("JOIN kelas_mahasiswas ON mahasiswas.nim = kelas_mahasiswas.nim").
		Joins("JOIN kelas ON kelas_mahasiswas.kelas_id = kelas.id").
		Where("kelas.kode_tahun_akademik = ?", currentTahunAkademik.Kode_tahun_akademik)

	if err := query.Count(&totalMahasiswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total mahasiswa: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          http.StatusOK,
		"message":         "Berhasil mengambil ringkasan akademik",
		"tahun_akademik":  currentTahunAkademik.Tahun,
		"semester":        currentTahunAkademik.Semester,
		"total_kelas":     totalKelas,
		"total_mahasiswa": totalMahasiswa,
	})
}
