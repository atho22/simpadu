package prodi

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

func (Idb *InDB) CreateKelas(c *gin.Context) {
	var InputKelas struct {
		Kode            string `json:"kode" binding:"required"`
		Nama            string `json:"nama" binding:"required"`
		ProdiID         uint   `json:"prodi_id" binding:"required"`
		TahunAkademikID uint   `json:"tahun_akademik_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&InputKelas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Request Tidak Valid"})
		return
	}

	// Validasi token seperti semula
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Cek Tahun Akademik
	var tahunAkademik structs.Tahun_akademik
	if err := Idb.DB.First(&tahunAkademik, InputKelas.TahunAkademikID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tahun Akademik Tidak Ditemukan"})
		return
	}

	// Cek Prodi
	var prodi structs.Prodi
	if err := Idb.DB.First(&prodi, InputKelas.ProdiID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prodi Tidak Ditemukan"})
		return
	}

	// Simpan data kelas
	kelas := structs.Kelas{
		Kode_kelas:          InputKelas.Kode,
		Nama:                InputKelas.Nama,
		Kode_tahun_akademik: tahunAkademik.Kode_tahun_akademik,
		Kode_prodi:          prodi.Kode_prodi,
		Created_by:          claims["user_id"].(string),
	}

	if err := Idb.DB.Create(&kelas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Membuat Kelas"})
		return
	}

	// Hitung mahasiswa (meskipun 0)
	var jumlahMahasiswa int64
	if err := Idb.DB.Table("kelas_mahasiswa").Where("kelas_id = ?", kelas.ID).Count(&jumlahMahasiswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung jumlah mahasiswa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Kelas Berhasil Dibuat",
		"data":             kelas,
		"jumlah_mahasiswa": jumlahMahasiswa,
	})
}

func (Idb *InDB) GetAllKelas(c *gin.Context) {
	var kelas []structs.Kelas

	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	if err := Idb.DB.Find(&kelas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Mengambil Data Kelas"})
		return
	}

	if len(kelas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak Ada Data Kelas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil Mengambil Data Kelas",
		"data":    kelas,
	})

}

func (Idb *InDB) DetailKelas(c *gin.Context) {
	var (
		kelas       structs.Kelas
		kelasDetail structs.Kelas_detail_response
	)
	kelas_id := c.Param("id")

	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Get basic kelas info
	if err := Idb.DB.Table("kelas").Where("id = ?", kelas_id).First(&kelas).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas Tidak Ditemukan"})
		return
	}

	// Get prodi info
	var prodi structs.Prodi
	if err := Idb.DB.Where("kode_prodi = ?", kelas.Kode_prodi).First(&prodi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data prodi"})
		return
	}

	// Get tahun akademik info
	var tahunAkademik structs.Tahun_akademik
	if err := Idb.DB.Where("kode_tahun_akademik = ?", kelas.Kode_tahun_akademik).First(&tahunAkademik).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tahun akademik"})
		return
	}

	// Get created by user info
	var createdByUser structs.Users
	if err := Idb.DB.Where("user_id = ?", kelas.Created_by).First(&createdByUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pembuat kelas"})
		return
	}

	// Get mahasiswa in kelas
	var mahasiswaList []structs.Mahasiswa_in_kelas
	if err := Idb.DB.Table("kelas_mahasiswas").
		Select("mahasiswas.nim, mahasiswas.nama, mahasiswas.angkatan, kelas_mahasiswas.created_at as assigned_at, users.username as assigned_by").
		Joins("JOIN mahasiswas ON kelas_mahasiswas.nim = mahasiswas.nim").
		Joins("JOIN users ON kelas_mahasiswas.created_by = users.user_id").
		Where("kelas_mahasiswas.kelas_id = ?", kelas_id).
		Scan(&mahasiswaList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data mahasiswa"})
		return
	}

	// Get available mahasiswa (not in any kelas for this tahun akademik)
	var availableMahasiswa []structs.Mahasiswa_simple
	if err := Idb.DB.Table("mahasiswas").
		Select("mahasiswas.nim, mahasiswas.nama, mahasiswas.angkatan").
		Where("mahasiswas.kode_prodi = ? AND mahasiswas.nim NOT IN (SELECT nim FROM kelas_mahasiswas WHERE kelas_id IN (SELECT id FROM kelas WHERE kode_tahun_akademik = ?))",
			kelas.Kode_prodi, kelas.Kode_tahun_akademik).
		Scan(&availableMahasiswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data mahasiswa yang tersedia"})
		return
	}

	// Count total mahasiswa
	var totalMahasiswa int64
	if err := Idb.DB.Table("kelas_mahasiswas").Where("kelas_id = ?", kelas_id).Count(&totalMahasiswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung jumlah mahasiswa"})
		return
	}

	// Construct response
	kelasDetail = structs.Kelas_detail_response{
		Kode:       kelas.Kode_kelas,
		Nama:       kelas.Nama,
		Kode_prodi: kelas.Kode_prodi,
		Nama_prodi: prodi.Nama_prodi,
		Tahun_akademik: structs.Tahun_akademik_simple{
			Tahun:    tahunAkademik.Tahun,
			Semester: tahunAkademik.Semester,
			Status:   tahunAkademik.Status,
		},
		Total_mahasiswa: int(totalMahasiswa),
		Created_by_nama: createdByUser.Username,
		Mahasiswa_list:  mahasiswaList,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil Mengambil Data Kelas",
		"data":    kelasDetail,
	})
}
