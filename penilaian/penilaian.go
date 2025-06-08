package penilaian

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

// CreatePenilaian creates a new grade entry
func (Idb *InDB) CreatePenilaian(c *gin.Context) {
	var input struct {
		Kode_matakuliah string  `json:"kode_matakuliah" binding:"required"`
		Kelas_id        uint    `json:"kelas_id" binding:"required"`
		NIM             string  `json:"nim" binding:"required"`
		Tugas           float64 `json:"tugas"`
		UTS             float64 `json:"uts"`
		UAS             float64 `json:"uas"`
		Keterangan      string  `json:"keterangan"`
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

	// Validasi role (dosen atau admin prodi)
	if claims["role"] != "dosen" && claims["role"] != "admin_prodi" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	// Validasi dosen adalah pengajar mata kuliah
	if claims["role"] == "dosen" {
		var jadwal structs.Jadwal_matakuliah
		if err := Idb.DB.Where("kode_matakuliah = ? AND kelas_id = ? AND kode_dosen = ?",
			input.Kode_matakuliah, input.Kelas_id, claims["user_id"]).First(&jadwal).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda bukan pengajar mata kuliah ini"})
			return
		}
	}

	// Validasi mahasiswa terdaftar di kelas
	var kelasMahasiswa structs.Kelas_mahasiswa
	if err := Idb.DB.Where("kelas_id = ? AND nim = ?", input.Kelas_id, input.NIM).First(&kelasMahasiswa).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mahasiswa tidak terdaftar di kelas ini"})
		return
	}

	// Hitung nilai akhir
	nilaiAkhir := (input.Tugas * 0.3) + (input.UTS * 0.3) + (input.UAS * 0.4)
	var grade string
	switch {
	case nilaiAkhir >= 85:
		grade = "A"
	case nilaiAkhir >= 80:
		grade = "A-"
	case nilaiAkhir >= 75:
		grade = "B+"
	case nilaiAkhir >= 70:
		grade = "B"
	case nilaiAkhir >= 65:
		grade = "B-"
	case nilaiAkhir >= 60:
		grade = "C+"
	case nilaiAkhir >= 55:
		grade = "C"
	case nilaiAkhir >= 50:
		grade = "C-"
	case nilaiAkhir >= 40:
		grade = "D"
	default:
		grade = "E"
	}

	// Simpan penilaian
	penilaian := structs.Penilaian{
		Kode_matakuliah: input.Kode_matakuliah,
		Kelas_id:        input.Kelas_id,
		NIM:             input.NIM,
		Tugas:           input.Tugas,
		UTS:             input.UTS,
		UAS:             input.UAS,
		Nilai_akhir:     nilaiAkhir,
		Grade:           grade,
		Keterangan:      input.Keterangan,
		Created_at:      time.Now(),
		Created_by:      claims["user_id"].(string),
	}

	if err := Idb.DB.Create(&penilaian).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan penilaian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penilaian berhasil disimpan",
		"data":    penilaian,
	})
}

// GetPenilaianByKelas gets all grades for a specific class
func (Idb *InDB) GetPenilaianByKelas(c *gin.Context) {
	kodeMatakuliah := c.Param("kode_matakuliah")
	kelasID := c.Param("kelas_id")

	var penilaian []structs.Penilaian
	if err := Idb.DB.Where("kode_matakuliah = ? AND kelas_id = ?", kodeMatakuliah, kelasID).
		Order("nim").
		Find(&penilaian).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penilaian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": penilaian,
	})
}

// GetPenilaianByMahasiswa gets grades for a specific student
func (Idb *InDB) GetPenilaianByMahasiswa(c *gin.Context) {
	nim := c.Param("nim")

	var penilaian []structs.Penilaian
	if err := Idb.DB.Where("nim = ?", nim).
		Order("kode_matakuliah").
		Find(&penilaian).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penilaian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": penilaian,
	})
}

// UpdatePenilaian updates a grade entry
func (Idb *InDB) UpdatePenilaian(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Tugas      float64 `json:"tugas"`
		UTS        float64 `json:"uts"`
		UAS        float64 `json:"uas"`
		Keterangan string  `json:"keterangan"`
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
	if claims["role"] != "dosen" && claims["role"] != "admin_prodi" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	// Ambil data penilaian yang akan diupdate
	var penilaian structs.Penilaian
	if err := Idb.DB.First(&penilaian, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data penilaian tidak ditemukan"})
		return
	}

	// Validasi dosen adalah pengajar mata kuliah
	if claims["role"] == "dosen" {
		var jadwal structs.Jadwal_matakuliah
		if err := Idb.DB.Where("kode_matakuliah = ? AND kelas_id = ? AND kode_dosen = ?",
			penilaian.Kode_matakuliah, penilaian.Kelas_id, claims["user_id"]).First(&jadwal).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda bukan pengajar mata kuliah ini"})
			return
		}
	}

	// Hitung nilai akhir
	nilaiAkhir := (input.Tugas * 0.3) + (input.UTS * 0.3) + (input.UAS * 0.4)
	var grade string
	switch {
	case nilaiAkhir >= 85:
		grade = "A"
	case nilaiAkhir >= 80:
		grade = "A-"
	case nilaiAkhir >= 75:
		grade = "B+"
	case nilaiAkhir >= 70:
		grade = "B"
	case nilaiAkhir >= 65:
		grade = "B-"
	case nilaiAkhir >= 60:
		grade = "C+"
	case nilaiAkhir >= 55:
		grade = "C"
	case nilaiAkhir >= 50:
		grade = "C-"
	case nilaiAkhir >= 40:
		grade = "D"
	default:
		grade = "E"
	}

	// Update penilaian
	updates := map[string]interface{}{
		"tugas":       input.Tugas,
		"uts":         input.UTS,
		"uas":         input.UAS,
		"nilai_akhir": nilaiAkhir,
		"grade":       grade,
		"keterangan":  input.Keterangan,
		"updated_at":  time.Now(),
		"updated_by":  claims["user_id"].(string),
	}

	if err := Idb.DB.Model(&structs.Penilaian{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate penilaian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penilaian berhasil diupdate",
	})
}

// GetRekapNilai gets grade summary for a class
func (Idb *InDB) GetRekapNilai(c *gin.Context) {
	kodeMatakuliah := c.Param("kode_matakuliah")
	kelasID := c.Param("kelas_id")

	var rekap struct {
		Kode_matakuliah string              `json:"kode_matakuliah"`
		Nama_matakuliah string              `json:"nama_matakuliah"`
		Kelas_id        uint                `json:"kelas_id"`
		Nama_kelas      string              `json:"nama_kelas"`
		Total_mahasiswa int                 `json:"total_mahasiswa"`
		Rata_rata       float64             `json:"rata_rata"`
		Distribusi      map[string]int      `json:"distribusi"`
		Detail_nilai    []structs.Penilaian `json:"detail_nilai"`
	}

	// Ambil data mata kuliah dan kelas
	var matakuliah structs.Mata_kuliah
	if err := Idb.DB.Where("kode_matakuliah = ?", kodeMatakuliah).First(&matakuliah).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata kuliah tidak ditemukan"})
		return
	}

	var kelas structs.Kelas
	if err := Idb.DB.First(&kelas, kelasID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas tidak ditemukan"})
		return
	}

	// Ambil semua nilai
	var penilaian []structs.Penilaian
	if err := Idb.DB.Where("kode_matakuliah = ? AND kelas_id = ?", kodeMatakuliah, kelasID).
		Find(&penilaian).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penilaian"})
		return
	}

	// Hitung statistik
	var totalNilai float64
	distribusi := make(map[string]int)
	for _, p := range penilaian {
		totalNilai += p.Nilai_akhir
		distribusi[p.Grade]++
	}

	rekap.Kode_matakuliah = kodeMatakuliah
	rekap.Nama_matakuliah = matakuliah.Nama_matakuliah
	rekap.Kelas_id = kelas.ID
	rekap.Nama_kelas = kelas.Nama
	rekap.Total_mahasiswa = len(penilaian)
	if len(penilaian) > 0 {
		rekap.Rata_rata = totalNilai / float64(len(penilaian))
	}
	rekap.Distribusi = distribusi
	rekap.Detail_nilai = penilaian

	c.JSON(http.StatusOK, gin.H{
		"data": rekap,
	})
}
