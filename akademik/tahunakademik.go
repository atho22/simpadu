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
