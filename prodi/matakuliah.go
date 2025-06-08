package prodi

import (
	"net/http"
	"simpadu/helper"
	"simpadu/structs"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (InDB *InDB) CreateMatakuliah(c *gin.Context) {
	var InputMatakuliah struct {
		Kode_matakuliah     string `json:"kode_matakuliah" binding:"required"`
		Nama_matakuliah     string `json:"nama_matakuliah" binding:"required"`
		Sks                 int    `json:"sks" binding:"required"`
		Kode_tahun_akademik string `json:"kode_tahun_akademik" binding:"required"`
		Semester            int    `json:"semester" binding:"required"`
		Kode_prodi          string `json:"kode_prodi" binding:"required"`
	}

	if err := c.ShouldBindJSON(&InputMatakuliah); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Request Tidak Valid"})
		return
	}
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}
	var existingMatakuliah structs.Mata_kuliah
	if err := InDB.DB.Where("kode_matakuliah = ?", InputMatakuliah.Kode_matakuliah).First(&existingMatakuliah).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Matakuliah Sudah Ada"})
		return
	}

	matakuliah := structs.Mata_kuliah{
		Kode_matakuliah:     InputMatakuliah.Kode_matakuliah,
		Nama_matakuliah:     InputMatakuliah.Nama_matakuliah,
		SKS:                 InputMatakuliah.Sks,
		Semester:            InputMatakuliah.Semester,
		Kode_prodi:          InputMatakuliah.Kode_prodi,
		Kode_tahun_akademik: InputMatakuliah.Kode_tahun_akademik,
		Status:              "Aktif",
		Created_by:          claims["user_id"].(string),
		Created_at:          time.Now(),
		Updated_at:          time.Now(),
	}

	if err := InDB.DB.Create(&matakuliah).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Membuat Matakuliah"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Matakuliah Berhasil Dibuat"})
}

func (InDB *InDB) GetAllMatakuliah(c *gin.Context) {
	var matakuliah []structs.Mata_kuliah
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}
	if err := InDB.DB.Find(&matakuliah).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Mengambil Data Matakuliah"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": matakuliah})
}

func (InDB *InDB) DetailMatakuliah(c *gin.Context) {
	var matakuliah structs.Mata_kuliah
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}
	if err := InDB.DB.Where("kode_matakuliah = ?", c.Param("kode_matakuliah")).First(&matakuliah).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Matakuliah Tidak Ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": matakuliah})
}
