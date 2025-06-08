package mahasiswa

import (
	"net/http"
	"simpadu/structs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InDB struct {
	DB *gorm.DB
}

// GetIPKMahasiswa menampilkan IPK dan IPS mahasiswa
func (h *InDB) GetIPKMahasiswa(c *gin.Context) {
	nim := c.Param("nim")

	// Ambil data mahasiswa
	var mahasiswa structs.Mahasiswa
	if err := h.DB.Where("nim = ?", nim).First(&mahasiswa).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa tidak ditemukan"})
		return
	}

	// Ambil data prodi
	var prodi structs.Prodi
	if err := h.DB.Where("kode_prodi = ?", mahasiswa.Kode_prodi).First(&prodi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data prodi tidak ditemukan"})
		return
	}

	// Ambil semua nilai mahasiswa
	var nilai []structs.Penilaian
	if err := h.DB.Where("nim = ?", nim).Find(&nilai).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data nilai"})
		return
	}

	// Kelompokkan nilai berdasarkan semester
	semesterMap := make(map[string]*structs.SemesterNilai)
	var totalSKS int

	for _, n := range nilai {
		// Ambil data mata kuliah
		var mk structs.Mata_kuliah
		if err := h.DB.Where("kode_matakuliah = ?", n.Kode_matakuliah).First(&mk).Error; err != nil {
			continue
		}

		// Ambil data tahun akademik
		var ta structs.Tahun_akademik
		if err := h.DB.Where("kode_tahun_akademik = ?", mk.Kode_tahun_akademik).First(&ta).Error; err != nil {
			continue
		}

		// Buat key semester
		semesterKey := ta.Kode_tahun_akademik

		// Inisialisasi semester jika belum ada
		if _, exists := semesterMap[semesterKey]; !exists {
			semesterMap[semesterKey] = &structs.SemesterNilai{
				Tahun_akademik: ta.Kode_tahun_akademik,
				Semester:       ta.Semester,
				Mata_kuliah:    make([]structs.NilaiMK, 0),
			}
		}

		// Tambahkan ke semester
		semesterMap[semesterKey].Mata_kuliah = append(semesterMap[semesterKey].Mata_kuliah, structs.NilaiMK{
			Kode_matakuliah: mk.Kode_matakuliah,
			Nama_matakuliah: mk.Nama_matakuliah,
			SKS:             mk.SKS,
			Nilai_akhir:     n.Nilai_akhir,
			Grade:           n.Grade,
		})
		semesterMap[semesterKey].Total_sks += mk.SKS
		totalSKS += mk.SKS
	}

	// Hitung IPS untuk setiap semester
	semesters := calculateIPS(semesterMap)

	// Hitung IPK
	ipk := calculateIPK(semesters)

	response := structs.IPKResponse{
		NIM:        mahasiswa.NIM,
		Nama:       mahasiswa.Nama,
		Kode_prodi: mahasiswa.Kode_prodi,
		Nama_prodi: prodi.Nama_prodi,
		Angkatan:   mahasiswa.Angkatan,
		IPK:        ipk,
		Total_sks:  totalSKS,
		Semester:   semesters,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data IPK dan IPS berhasil diambil",
		"data":    response,
	})
}

// GetIPSMahasiswa menampilkan IPS mahasiswa untuk semester tertentu
func (h *InDB) GetIPSMahasiswa(c *gin.Context) {
	nim := c.Param("nim")
	tahunAkademik := c.Param("tahun_akademik")

	// Ambil data mahasiswa
	var mahasiswa structs.Mahasiswa
	if err := h.DB.Where("nim = ?", nim).First(&mahasiswa).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa tidak ditemukan"})
		return
	}

	// Ambil semua nilai mahasiswa untuk semester tertentu
	var nilai []structs.Penilaian
	if err := h.DB.Where("nim = ?", nim).Find(&nilai).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data nilai"})
		return
	}

	// Kelompokkan nilai berdasarkan semester
	semesterMap := make(map[string]*structs.SemesterNilai)

	for _, n := range nilai {
		// Ambil data mata kuliah
		var mk structs.Mata_kuliah
		if err := h.DB.Where("kode_matakuliah = ?", n.Kode_matakuliah).First(&mk).Error; err != nil {
			continue
		}

		// Ambil data tahun akademik
		var ta structs.Tahun_akademik
		if err := h.DB.Where("kode_tahun_akademik = ?", mk.Kode_tahun_akademik).First(&ta).Error; err != nil {
			continue
		}

		// Skip jika bukan semester yang diminta
		if ta.Kode_tahun_akademik != tahunAkademik {
			continue
		}

		// Buat key semester
		semesterKey := ta.Kode_tahun_akademik

		// Inisialisasi semester jika belum ada
		if _, exists := semesterMap[semesterKey]; !exists {
			semesterMap[semesterKey] = &structs.SemesterNilai{
				Tahun_akademik: ta.Kode_tahun_akademik,
				Semester:       ta.Semester,
				Mata_kuliah:    make([]structs.NilaiMK, 0),
			}
		}

		// Tambahkan ke semester
		semesterMap[semesterKey].Mata_kuliah = append(semesterMap[semesterKey].Mata_kuliah, structs.NilaiMK{
			Kode_matakuliah: mk.Kode_matakuliah,
			Nama_matakuliah: mk.Nama_matakuliah,
			SKS:             mk.SKS,
			Nilai_akhir:     n.Nilai_akhir,
			Grade:           n.Grade,
		})
		semesterMap[semesterKey].Total_sks += mk.SKS
	}

	// Hitung IPS untuk semester yang diminta
	semesters := calculateIPS(semesterMap)

	if len(semesters) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data nilai tidak ditemukan untuk semester tersebut"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data IPS berhasil diambil",
		"data":    semesters[0], // Ambil semester pertama karena hanya ada satu semester
	})
}

// calculateIPS menghitung IPS untuk setiap semester
func calculateIPS(semesterMap map[string]*structs.SemesterNilai) []structs.SemesterNilai {
	for _, semester := range semesterMap {
		var totalBobotSemester float64
		for _, mk := range semester.Mata_kuliah {
			bobot := calculateBobot(mk.Grade)
			totalBobotSemester += float64(mk.SKS) * bobot
		}
		if semester.Total_sks > 0 {
			semester.IPS = totalBobotSemester / float64(semester.Total_sks)
		}
	}

	// Konversi map ke slice
	semesters := make([]structs.SemesterNilai, 0, len(semesterMap))
	for _, semester := range semesterMap {
		semesters = append(semesters, *semester)
	}

	return semesters
}

// calculateIPK menghitung IPK dari semua IPS
func calculateIPK(semesters []structs.SemesterNilai) float64 {
	var totalIPS float64
	var semesterCount int

	for _, semester := range semesters {
		if semester.Total_sks > 0 {
			totalIPS += semester.IPS
			semesterCount++
		}
	}

	if semesterCount > 0 {
		return totalIPS / float64(semesterCount)
	}
	return 0.0
}

// calculateBobot menghitung bobot nilai berdasarkan grade
func calculateBobot(grade string) float64 {
	switch grade {
	case "A":
		return 4.0
	case "A-":
		return 3.75
	case "B+":
		return 3.5
	case "B":
		return 3.0
	case "B-":
		return 2.75
	case "C+":
		return 2.5
	case "C":
		return 2.0
	case "C-":
		return 1.75
	case "D+":
		return 1.5
	case "D":
		return 1.0
	case "E":
		return 0.0
	default:
		return 0.0
	}
}
