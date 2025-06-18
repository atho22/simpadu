package seeders

import (
	"fmt"
	"simpadu/structs"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedTahunAkademik seeds academic year data
func SeedTahunAkademik(db *gorm.DB) error {
	tahunAkademik := []structs.Tahun_akademik{
		{
			Kode_tahun_akademik: "TA2024001",
			Tahun:               "2024/2025",
			Semester:            "Ganjil",
			Status:              "Aktif",
			Tanggal_mulai:       time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC),
			Tanggal_selesai:     time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
		{
			Kode_tahun_akademik: "TA2024002",
			Tahun:               "2024/2025",
			Semester:            "Genap",
			Status:              "Tidak Aktif",
			Tanggal_mulai:       time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			Tanggal_selesai:     time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC),
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
		{
			Kode_tahun_akademik: "TA2023001",
			Tahun:               "2023/2024",
			Semester:            "Ganjil",
			Status:              "Tidak Aktif",
			Tanggal_mulai:       time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			Tanggal_selesai:     time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
		{
			Kode_tahun_akademik: "TA2023002",
			Tahun:               "2023/2024",
			Semester:            "Genap",
			Status:              "Tidak Aktif",
			Tanggal_mulai:       time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			Tanggal_selesai:     time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC),
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
		{
			Kode_tahun_akademik: "TA2025001",
			Tahun:               "2025/2026",
			Semester:            "Ganjil",
			Status:              "Tidak Aktif",
			Tanggal_mulai:       time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			Tanggal_selesai:     time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
	}
	return db.Create(&tahunAkademik).Error
}

// SeedUsers seeds user data
func SeedUsers(db *gorm.DB) error {
	// Truncate (delete) the users table first to avoid duplicate entry errors
	if err := db.Exec("TRUNCATE TABLE users").Error; err != nil {
		return err
	}

	// Hash password once since it's the same for all users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	users := []structs.Users{
		{
			User_id:  "USR001",
			Username: "admin_sistem",
			Email:    "admin@university.ac.id",
			Password: string(hashedPassword),
			Role:     "admin_akademik",
			Status:   "Aktif",
		},
		{
			User_id:  "USR002",
			Username: "admin_prodi_ti",
			Email:    "admin.ti@university.ac.id",
			Password: string(hashedPassword),
			Role:     "admin_prodi",
			Status:   "Aktif",
		},
		{
			User_id:  "DSN001",
			Username: "dosen_1",
			Email:    "dosen1@university.ac.id",
			Password: string(hashedPassword),
			Role:     "dosen",
			Status:   "Aktif",
		},
		{
			User_id:  "MHS001",
			Username: "mahasiswa_1",
			Email:    "mahasiswa1@student.university.ac.id",
			Password: string(hashedPassword),
			Role:     "mahasiswa",
			Status:   "Aktif",
		},
		{
			User_id:  "STF001",
			Username: "staff_akademik",
			Email:    "staff@university.ac.id",
			Password: string(hashedPassword),
			Role:     "staff_akademik",
			Status:   "Aktif",
		},
	}
	return db.Create(&users).Error
}

// SeedJurusan seeds department data
func SeedJurusan(db *gorm.DB) error {
	jurusan := []structs.Jurusan{
		{
			Kode_jurusan: "JUR001",
			Nama_jurusan: "Teknik Informatika dan Komputer",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_jurusan: "JUR002",
			Nama_jurusan: "Teknik Elektro",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_jurusan: "JUR003",
			Nama_jurusan: "Teknik Mesin",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_jurusan: "JUR004",
			Nama_jurusan: "Teknik Sipil",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_jurusan: "JUR005",
			Nama_jurusan: "Ekonomi dan Bisnis",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
	}
	return db.Create(&jurusan).Error
}

// SeedProdi seeds study program data
func SeedProdi(db *gorm.DB) error {
	prodi := []structs.Prodi{
		{
			Kode_prodi:   "PRD001",
			Nama_prodi:   "Teknik Informatika",
			Kode_jurusan: "JUR001",
			Jenjang:      "S1",
			User_id:      "USR002",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_prodi:   "PRD002",
			Nama_prodi:   "Sistem Informasi",
			Kode_jurusan: "JUR001",
			Jenjang:      "S1",
			User_id:      "",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_prodi:   "PRD003",
			Nama_prodi:   "Teknik Elektro",
			Kode_jurusan: "JUR002",
			Jenjang:      "S1",
			User_id:      "",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_prodi:   "PRD004",
			Nama_prodi:   "Teknik Mesin",
			Kode_jurusan: "JUR003",
			Jenjang:      "S1",
			User_id:      "",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_prodi:   "PRD005",
			Nama_prodi:   "Manajemen",
			Kode_jurusan: "JUR005",
			Jenjang:      "S1",
			User_id:      "",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
	}
	return db.Create(&prodi).Error
}

// SeedStrukturalJurusan seeds department structural data
func SeedStrukturalJurusan(db *gorm.DB) error {
	struktural := []structs.Struktural_Jurusan{
		{
			User_id:      "STRJ001",
			Kode_jurusan: "JUR001",
			Nama:         "Dr. Ahmad Wijaya",
			Nama_lengkap: "Dr. Ahmad Wijaya, S.T., M.Kom.",
			Alamat:       "Jl. Merdeka No. 123, Jakarta",
			No_hp:        "081234567890",
			Email:        "ahmad.wijaya@university.ac.id",
			NIP:          "198501012010011001",
			Jabatan:      "Ketua Jurusan",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			User_id:      "STRJ002",
			Kode_jurusan: "JUR001",
			Nama:         "Dr. Siti Rahayu",
			Nama_lengkap: "Dr. Siti Rahayu, S.T., M.T.",
			Alamat:       "Jl. Sudirman No. 456, Jakarta",
			No_hp:        "081234567891",
			Email:        "siti.rahayu@university.ac.id",
			NIP:          "198502022010012002",
			Jabatan:      "Sekretaris Jurusan",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			User_id:      "STRJ003",
			Kode_jurusan: "JUR002",
			Nama:         "Prof. Budi Santoso",
			Nama_lengkap: "Prof. Budi Santoso, S.T., M.T., Ph.D.",
			Alamat:       "Jl. Thamrin No. 789, Jakarta",
			No_hp:        "081234567892",
			Email:        "budi.santoso@university.ac.id",
			NIP:          "197503032005011003",
			Jabatan:      "Ketua Jurusan",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			User_id:      "STRJ004",
			Kode_jurusan: "JUR003",
			Nama:         "Dr. Ir. Andi Pratama",
			Nama_lengkap: "Dr. Ir. Andi Pratama, M.T.",
			Alamat:       "Jl. Gatot Subroto No. 321, Jakarta",
			No_hp:        "081234567893",
			Email:        "andi.pratama@university.ac.id",
			NIP:          "198004042008011004",
			Jabatan:      "Ketua Jurusan",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			User_id:      "STRJ005",
			Kode_jurusan: "JUR004",
			Nama:         "Dr. Maya Sari",
			Nama_lengkap: "Dr. Maya Sari, S.T., M.T.",
			Alamat:       "Jl. Kuningan No. 654, Jakarta",
			No_hp:        "081234567894",
			Email:        "maya.sari@university.ac.id",
			NIP:          "198605052012012005",
			Jabatan:      "Ketua Jurusan",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
	}
	return db.Create(&struktural).Error
}

// SeedStrukturalProdi seeds study program structural data
func SeedStrukturalProdi(db *gorm.DB) error {
	struktural := []structs.Struktural_prodi{
		{
			User_id:      "STRP001",
			Kode_prodi:   "PRD001",
			Nama:         "Dr. Reza Firmansyah",
			Nama_lengkap: "Dr. Reza Firmansyah, S.Kom., M.Kom.",
			Alamat:       "Jl. Kebon Jeruk No. 111, Jakarta",
			No_hp:        "081234567895",
			Email:        "reza.firmansyah@university.ac.id",
			NIP:          "198706062015011006",
			Jabatan:      "Ketua Program Studi",
			NoSK:         "SK/001/2024",
			TanggalSK:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			User_id:      "STRP002",
			Kode_prodi:   "PRD001",
			Nama:         "Lisa Permata",
			Nama_lengkap: "Lisa Permata, S.Kom., M.Kom.",
			Alamat:       "Jl. Cempaka No. 222, Jakarta",
			No_hp:        "081234567896",
			Email:        "lisa.permata@university.ac.id",
			NIP:          "198907072018012007",
			Jabatan:      "Sekretaris Program Studi",
			NoSK:         "SK/002/2024",
			TanggalSK:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			User_id:      "STRP003",
			Kode_prodi:   "PRD002",
			Nama:         "Dr. Hendra Wijaya",
			Nama_lengkap: "Dr. Hendra Wijaya, S.Si., M.Kom.",
			Alamat:       "Jl. Melati No. 333, Jakarta",
			No_hp:        "081234567897",
			Email:        "hendra.wijaya@university.ac.id",
			NIP:          "198408082016011008",
			Jabatan:      "Ketua Program Studi",
			NoSK:         "SK/003/2024",
			TanggalSK:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			User_id:      "STRP004",
			Kode_prodi:   "PRD003",
			Nama:         "Prof. Indra Gunawan",
			Nama_lengkap: "Prof. Indra Gunawan, S.T., M.T., Ph.D.",
			Alamat:       "Jl. Mawar No. 444, Jakarta",
			No_hp:        "081234567898",
			Email:        "indra.gunawan@university.ac.id",
			NIP:          "197509092000011009",
			Jabatan:      "Ketua Program Studi",
			NoSK:         "SK/004/2024",
			TanggalSK:    time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			User_id:      "STRP005",
			Kode_prodi:   "PRD005",
			Nama:         "Dr. Ratna Sari",
			Nama_lengkap: "Dr. Ratna Sari, S.E., M.M.",
			Alamat:       "Jl. Dahlia No. 555, Jakarta",
			No_hp:        "081234567899",
			Email:        "ratna.sari@university.ac.id",
			NIP:          "198110102019012010",
			Jabatan:      "Ketua Program Studi",
			NoSK:         "SK/005/2024",
			TanggalSK:    time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}
	return db.Create(&struktural).Error
}

// SeedStaffAkademik seeds academic staff data
func SeedStaffAkademik(db *gorm.DB) error {
	staff := []structs.Staff_akademik{
		{
			User_id:    "STF001",
			Nama:       "Sari Dewi",
			Alamat:     "Jl. Kenanga No. 777, Jakarta",
			No_hp:      "081234560001",
			Email:      "sari.dewi@university.ac.id",
			Status:     "Aktif",
			Bagian:     "Akademik",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			User_id:    "STF002",
			Nama:       "Bambang Setiawan",
			Alamat:     "Jl. Anggrek No. 888, Jakarta",
			No_hp:      "081234560002",
			Email:      "bambang.setiawan@university.ac.id",
			Status:     "Aktif",
			Bagian:     "Kemahasiswaan",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			User_id:    "STF003",
			Nama:       "Rina Kartika",
			Alamat:     "Jl. Tulip No. 999, Jakarta",
			No_hp:      "081234560003",
			Email:      "rina.kartika@university.ac.id",
			Status:     "Aktif",
			Bagian:     "Keuangan",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			User_id:    "STF004",
			Nama:       "Dedi Kurniawan",
			Alamat:     "Jl. Sakura No. 101, Jakarta",
			No_hp:      "081234560004",
			Email:      "dedi.kurniawan@university.ac.id",
			Status:     "Aktif",
			Bagian:     "IT Support",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			User_id:    "STF005",
			Nama:       "Fitri Handayani",
			Alamat:     "Jl. Lily No. 102, Jakarta",
			No_hp:      "081234560005",
			Email:      "fitri.handayani@university.ac.id",
			Status:     "Tidak Aktif",
			Bagian:     "Perpustakaan",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
	}
	return db.Create(&staff).Error
}

// SeedDosen seeds lecturer data
func SeedDosen(db *gorm.DB) error {
	dosen := []structs.Dosen{
		{
			User_id:    "DSN001",
			Kode_dosen: "DOS001",
			NIP:        "198801012015041001",
			Nama:       "Dr. Agus Setiawan",
			Nidn:       "0101018801",
			Status:     "Aktif",
			Alamat:     "Jl. Flamboyan No. 201, Jakarta",
			No_hp:      "081234570001",
			Email:      "agus.setiawan@university.ac.id",
			Kode_prodi: "PRD001",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			User_id:    "DSN002",
			Kode_dosen: "DOS002",
			NIP:        "198902022016042002",
			Nama:       "Dr. Wulan Sari",
			Nidn:       "0202028902",
			Status:     "Aktif",
			Alamat:     "Jl. Kamboja No. 202, Jakarta",
			No_hp:      "081234570002",
			Email:      "wulan.sari@university.ac.id",
			Kode_prodi: "PRD001",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			User_id:    "DSN003",
			Kode_dosen: "DOS003",
			NIP:        "198703032017043003",
			Nama:       "Prof. Eko Prasetyo",
			Nidn:       "0303038703",
			Status:     "Aktif",
			Alamat:     "Jl. Teratai No. 203, Jakarta",
			No_hp:      "081234570003",
			Email:      "eko.prasetyo@university.ac.id",
			Kode_prodi: "PRD002",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			User_id:    "DSN004",
			Kode_dosen: "DOS004",
			NIP:        "198504042018044004",
			Nama:       "Dr. Nina Permatasari",
			Nidn:       "0404048504",
			Status:     "Aktif",
			Alamat:     "Jl. Bougenville No. 204, Jakarta",
			No_hp:      "081234570004",
			Email:      "nina.permatasari@university.ac.id",
			Kode_prodi: "PRD003",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			User_id:    "DSN005",
			Kode_dosen: "DOS005",
			NIP:        "198605052019045005",
			Nama:       "Dr. Rudi Hartono",
			Nidn:       "0505058605",
			Status:     "Tidak Aktif",
			Alamat:     "Jl. Azalea No. 205, Jakarta",
			No_hp:      "081234570005",
			Email:      "rudi.hartono@university.ac.id",
			Kode_prodi: "PRD004",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
	}
	return db.Create(&dosen).Error
}

// SeedKelas seeds class data
func SeedKelas(db *gorm.DB) error {
	kelas := []structs.Kelas{
		{
			Kode_kelas:          "TI-A",
			Nama:                "Teknik Informatika A",
			Kode_prodi:          "PRD001",
			Kode_tahun_akademik: "TA2024001",
			Created_by:          "USR002",
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
		{
			Kode_kelas:          "TI-B",
			Nama:                "Teknik Informatika B",
			Kode_prodi:          "PRD001",
			Kode_tahun_akademik: "TA2024001",
			Created_by:          "USR002",
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
		{
			Kode_kelas:          "SI-A",
			Nama:                "Sistem Informasi A",
			Kode_prodi:          "PRD002",
			Kode_tahun_akademik: "TA2024001",
			Created_by:          "USR002",
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
		{
			Kode_kelas:          "TE-A",
			Nama:                "Teknik Elektro A",
			Kode_prodi:          "PRD003",
			Kode_tahun_akademik: "TA2024001",
			Created_by:          "USR002",
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
		{
			Kode_kelas:          "MN-A",
			Nama:                "Manajemen A",
			Kode_prodi:          "PRD005",
			Kode_tahun_akademik: "TA2024001",
			Created_by:          "USR002",
			Created_at:          time.Now(),
			Updated_at:          time.Now(),
		},
	}
	return db.Create(&kelas).Error
}

// SeedMahasiswa seeds student data
func SeedMahasiswa(db *gorm.DB) error {
	mahasiswa := []structs.Mahasiswa{
		{
			User_id:       "MHS001",
			NIM:           "2024010001",
			Nama:          "Ahmad Rizki Pratama",
			Angkatan:      2024,
			Status:        "Aktif",
			Jenis_kelamin: "Laki-laki",
			Alamat:        "Jl. Mangga No. 301, Jakarta",
			No_hp:         "081234580001",
			Email:         "ahmad.rizki@student.university.ac.id",
			Kode_prodi:    "PRD001",
			Created_at:    time.Now(),
			Updated_at:    time.Now(),
		},
		{
			User_id:       "MHS002",
			NIM:           "2024010002",
			Nama:          "Siti Nurhaliza",
			Angkatan:      2024,
			Status:        "Aktif",
			Jenis_kelamin: "Perempuan",
			Alamat:        "Jl. Jeruk No. 302, Jakarta",
			No_hp:         "081234580002",
			Email:         "siti.nurhaliza@student.university.ac.id",
			Kode_prodi:    "PRD001",
			Created_at:    time.Now(),
			Updated_at:    time.Now(),
		},
		{
			User_id:       "MHS003",
			NIM:           "2024020001",
			Nama:          "Budi Santoso",
			Angkatan:      2024,
			Status:        "Aktif",
			Jenis_kelamin: "Laki-laki",
			Alamat:        "Jl. Apel No. 303, Jakarta",
			No_hp:         "081234580003",
			Email:         "budi.santoso@student.university.ac.id",
			Kode_prodi:    "PRD002",
			Created_at:    time.Now(),
			Updated_at:    time.Now(),
		},
		{
			User_id:       "MHS004",
			NIM:           "2024030001",
			Nama:          "Maya Puspita",
			Angkatan:      2024,
			Status:        "Aktif",
			Jenis_kelamin: "Perempuan",
			Alamat:        "Jl. Nanas No. 304, Jakarta",
			No_hp:         "081234580004",
			Email:         "maya.puspita@student.university.ac.id",
			Kode_prodi:    "PRD003",
			Created_at:    time.Now(),
			Updated_at:    time.Now(),
		},
		{
			User_id:       "MHS005",
			NIM:           "2023010001",
			Nama:          "Rian Firmansyah",
			Angkatan:      2023,
			Status:        "Cuti",
			Jenis_kelamin: "Laki-laki",
			Alamat:        "Jl. Durian No. 305, Jakarta",
			No_hp:         "081234580005",
			Email:         "rian.firmansyah@student.university.ac.id",
			Kode_prodi:    "PRD001",
			Created_at:    time.Now(),
			Updated_at:    time.Now(),
		},
	}
	return db.Create(&mahasiswa).Error
}

// SeedKelasMahasiswa seeds class-student relationship data
func SeedKelasMahasiswa(db *gorm.DB) error {
	kelasMahasiswa := []structs.Kelas_mahasiswa{
		{
			Kelas_id:   1,
			NIM:        "2024010001",
			Created_by: "USR002",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			Kelas_id:   1,
			NIM:        "2024010002",
			Created_by: "USR002",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			Kelas_id:   3,
			NIM:        "2024020001",
			Created_by: "USR002",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			Kelas_id:   4,
			NIM:        "2024030001",
			Created_by: "USR002",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
		{
			Kelas_id:   1,
			NIM:        "2023010001",
			Created_by: "USR002",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
	}
	return db.Create(&kelasMahasiswa).Error
}

// SeedMataKuliah seeds course data
func SeedMataKuliah(db *gorm.DB) error {
	mataKuliah := []structs.Mata_kuliah{
		{
			Kode_matakuliah:     "MK001",
			Nama_matakuliah:     "Pemrograman Dasar",
			Kode_prodi:          "PRD001",
			SKS:                 3,
			Semester:            1,
			Kode_tahun_akademik: "TA2024001",
			Status:              "Aktif",
			Created_at:          time.Now(),
			Created_by:          "DSN001",
			Updated_at:          time.Now(),
		},
		{
			Kode_matakuliah:     "MK002",
			Nama_matakuliah:     "Struktur Data",
			Kode_prodi:          "PRD001",
			SKS:                 3,
			Semester:            2,
			Kode_tahun_akademik: "TA2024001",
			Status:              "Aktif",
			Created_at:          time.Now(),
			Created_by:          "DSN002",
			Updated_at:          time.Now(),
		},
		{
			Kode_matakuliah:     "MK003",
			Nama_matakuliah:     "Basis Data",
			Kode_prodi:          "PRD002",
			SKS:                 3,
			Semester:            3,
			Kode_tahun_akademik: "TA2024001",
			Status:              "Aktif",
			Created_at:          time.Now(),
			Created_by:          "DSN003",
			Updated_at:          time.Now(),
		},
		{
			Kode_matakuliah:     "MK004",
			Nama_matakuliah:     "Rangkaian Listrik",
			Kode_prodi:          "PRD003",
			SKS:                 3,
			Semester:            1,
			Kode_tahun_akademik: "TA2024001",
			Status:              "Aktif",
			Created_at:          time.Now(),
			Created_by:          "DSN004",
			Updated_at:          time.Now(),
		},
		{
			Kode_matakuliah:     "MK005",
			Nama_matakuliah:     "Manajemen Strategis",
			Kode_prodi:          "PRD005",
			SKS:                 3,
			Semester:            5,
			Kode_tahun_akademik: "TA2024001",
			Status:              "Tidak Aktif",
			Created_at:          time.Now(),
			Created_by:          "DSN005",
			Updated_at:          time.Now(),
		},
	}
	return db.Create(&mataKuliah).Error
}

// SeedKodeRuangan seeds room code data
func SeedKodeRuangan(db *gorm.DB) error {
	ruangan := []structs.Kode_ruangan{
		{
			Kode_ruangan: "A-KLS-101",
			Nama_ruangan: "Ruang Kuliah A101",
			Kode_prodi:   "PRD001",
			Gedung:       "A",
			Tipe:         "KLS",
			Nomor:        101,
			Kapasitas:    40,
			Lokasi:       "Lantai 1",
			Status:       "Aktif",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_ruangan: "A-LAB-201",
			Nama_ruangan: "Laboratorium Komputer A201",
			Kode_prodi:   "PRD001",
			Gedung:       "A",
			Tipe:         "LAB",
			Nomor:        201,
			Kapasitas:    30,
			Lokasi:       "Lantai 2",
			Status:       "Aktif",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_ruangan: "B-KLS-102",
			Nama_ruangan: "Ruang Kuliah B102",
			Kode_prodi:   "PRD002",
			Gedung:       "B",
			Tipe:         "KLS",
			Nomor:        102,
			Kapasitas:    35,
			Lokasi:       "Lantai 1",
			Status:       "Aktif",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_ruangan: "C-LAB-301",
			Nama_ruangan: "Laboratorium Elektronika C301",
			Kode_prodi:   "PRD003",
			Gedung:       "C",
			Tipe:         "LAB",
			Nomor:        301,
			Kapasitas:    25,
			Lokasi:       "Lantai 3",
			Status:       "Maintenance",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
		{
			Kode_ruangan: "D-STU-401",
			Nama_ruangan: "Studio Desain D401",
			Kode_prodi:   "PRD005",
			Gedung:       "D",
			Tipe:         "STU",
			Nomor:        401,
			Kapasitas:    20,
			Lokasi:       "Lantai 4",
			Status:       "Tidak Aktif",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		},
	}
	return db.Create(&ruangan).Error
}

// SeedJadwalMatakuliah seeds course schedule data
func SeedJadwalMatakuliah(db *gorm.DB) error {
	jadwal := []structs.Jadwal_matakuliah{
		{
			Kode_matakuliah:     "MK001",
			Kelas_id:            1,
			Kode_dosen:          "DOS001",
			Kode_ruangan:        "A-KLS-101",
			Kode_prodi:          "PRD001",
			Jam_mulai:           time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
			Jam_selesai:         time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
			Hari:                "Senin",
			Kode_tahun_akademik: "TA2024001",
			Created_at:          time.Now(),
			Created_by:          "USR002",
		},
		{
			Kode_matakuliah:     "MK002",
			Kelas_id:            1,
			Kode_dosen:          "DOS002",
			Kode_ruangan:        "A-LAB-201",
			Kode_prodi:          "PRD001",
			Jam_mulai:           time.Date(2024, 1, 1, 10, 45, 0, 0, time.UTC),
			Jam_selesai:         time.Date(2024, 1, 1, 13, 15, 0, 0, time.UTC),
			Hari:                "Selasa",
			Kode_tahun_akademik: "TA2024001",
			Created_at:          time.Now(),
			Created_by:          "USR002",
		},
		{
			Kode_matakuliah:     "MK003",
			Kelas_id:            3,
			Kode_dosen:          "DOS003",
			Kode_ruangan:        "B-KLS-102",
			Kode_prodi:          "PRD002",
			Jam_mulai:           time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
			Jam_selesai:         time.Date(2024, 1, 1, 16, 30, 0, 0, time.UTC),
			Hari:                "Rabu",
			Kode_tahun_akademik: "TA2024001",
			Created_at:          time.Now(),
			Created_by:          "USR002",
		},
		{
			Kode_matakuliah:     "MK004",
			Kelas_id:            4,
			Kode_dosen:          "DOS004",
			Kode_ruangan:        "C-LAB-301",
			Kode_prodi:          "PRD003",
			Jam_mulai:           time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
			Jam_selesai:         time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
			Hari:                "Kamis",
			Kode_tahun_akademik: "TA2024001",
			Created_at:          time.Now(),
			Created_by:          "USR002",
		},
		{
			Kode_matakuliah:     "MK005",
			Kelas_id:            5,
			Kode_dosen:          "DOS005",
			Kode_ruangan:        "D-STU-401",
			Kode_prodi:          "PRD005",
			Jam_mulai:           time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
			Jam_selesai:         time.Date(2024, 1, 1, 15, 30, 0, 0, time.UTC),
			Hari:                "Jumat",
			Kode_tahun_akademik: "TA2024001",
			Created_at:          time.Now(),
			Created_by:          "USR002",
		},
	}
	return db.Create(&jadwal).Error
}

// SeedPertemuan seeds meeting data
func SeedPertemuan(db *gorm.DB) error {
	pertemuan := []structs.Pertemuan{
		{
			Kode_pertemuan:  "PTM001001",
			Kode_matakuliah: "MK001",
			Pertemuan_ke:    1,
			Dibuka:          true,
			Waktu_dibuka:    time.Now().Add(-2 * time.Hour),
			Waktu_ditutup:   time.Now().Add(-1 * time.Hour),
			Created_at:      time.Now(),
			Updated_at:      time.Now(),
			Dibuka_by:       "DOS001",
		},
		{
			Kode_pertemuan:  "PTM001002",
			Kode_matakuliah: "MK001",
			Pertemuan_ke:    2,
			Dibuka:          false,
			Waktu_dibuka:    time.Now().Add(24 * time.Hour),
			Waktu_ditutup:   time.Now().Add(26 * time.Hour),
			Created_at:      time.Now(),
			Updated_at:      time.Now(),
			Dibuka_by:       "DOS001",
		},
		{
			Kode_pertemuan:  "PTM002001",
			Kode_matakuliah: "MK002",
			Pertemuan_ke:    1,
			Dibuka:          true,
			Waktu_dibuka:    time.Now().Add(-3 * time.Hour),
			Waktu_ditutup:   time.Now().Add(-2 * time.Hour),
			Created_at:      time.Now(),
			Updated_at:      time.Now(),
			Dibuka_by:       "DOS002",
		},
		{
			Kode_pertemuan:  "PTM003001",
			Kode_matakuliah: "MK003",
			Pertemuan_ke:    1,
			Dibuka:          true,
			Waktu_dibuka:    time.Now().Add(-4 * time.Hour),
			Waktu_ditutup:   time.Now().Add(-3 * time.Hour),
			Created_at:      time.Now(),
			Updated_at:      time.Now(),
			Dibuka_by:       "DOS003",
		},
		{
			Kode_pertemuan:  "PTM004001",
			Kode_matakuliah: "MK004",
			Pertemuan_ke:    1,
			Dibuka:          false,
			Waktu_dibuka:    time.Now().Add(48 * time.Hour),
			Waktu_ditutup:   time.Now().Add(50 * time.Hour),
			Created_at:      time.Now(),
			Updated_at:      time.Now(),
			Dibuka_by:       "DOS004",
		},
	}
	return db.Create(&pertemuan).Error
}

// SeedAbsensi seeds attendance data
func SeedAbsensi(db *gorm.DB) error {
	absensi := []structs.Absensi{
		{
			Kode_matakuliah: "MK001",
			Kode_pertemuan:  "PTM001001",
			Nim:             "2024010001",
			Status:          "hadir",
			Keterangan:      "Tepat waktu",
			Created_at:      time.Now(),
		},
		{
			Kode_matakuliah: "MK001",
			Kode_pertemuan:  "PTM001001",
			Nim:             "2024010002",
			Status:          "hadir",
			Keterangan:      "Tepat waktu",
			Created_at:      time.Now(),
		},
		{
			Kode_matakuliah: "MK002",
			Kode_pertemuan:  "PTM002001",
			Nim:             "2024010001",
			Status:          "izin",
			Keterangan:      "Sakit",
			Created_at:      time.Now(),
		},
		{
			Kode_matakuliah: "MK003",
			Kode_pertemuan:  "PTM003001",
			Nim:             "2024020001",
			Status:          "hadir",
			Keterangan:      "Terlambat 10 menit",
			Created_at:      time.Now(),
		},
		{
			Kode_matakuliah: "MK001",
			Kode_pertemuan:  "PTM001001",
			Nim:             "2023010001",
			Status:          "alfa",
			Keterangan:      "Tidak hadir tanpa keterangan",
			Created_at:      time.Now(),
		},
	}
	return db.Create(&absensi).Error
}

// SeedPenilaian seeds grade data
func SeedPenilaian(db *gorm.DB) error {
	penilaian := []structs.Penilaian{
		{
			Kode_matakuliah: "MK001",
			NIM:             "2024010001",
			Kelas_id:        1,
			Tugas:           85,
			UTS:             80,
			UAS:             90,
			Nilai_akhir:     85,
			Grade:           "A",
			Keterangan:      "Lulus",
			Created_at:      time.Now(),
			Created_by:      "DOS001",
			Updated_at:      time.Now(),
		},
		{
			Kode_matakuliah: "MK001",
			NIM:             "2024010002",
			Kelas_id:        1,
			Tugas:           90,
			UTS:             85,
			UAS:             95,
			Nilai_akhir:     90,
			Grade:           "A",
			Keterangan:      "Lulus",
			Created_at:      time.Now(),
			Created_by:      "DOS001",
			Updated_at:      time.Now(),
		},
		{
			Kode_matakuliah: "MK002",
			NIM:             "2024010001",
			Kelas_id:        1,
			Tugas:           75,
			UTS:             80,
			UAS:             85,
			Nilai_akhir:     80,
			Grade:           "B",
			Keterangan:      "Lulus",
			Created_at:      time.Now(),
			Created_by:      "DOS002",
			Updated_at:      time.Now(),
		},
		{
			Kode_matakuliah: "MK003",
			NIM:             "2024020001",
			Kelas_id:        3,
			Tugas:           88,
			UTS:             85,
			UAS:             90,
			Nilai_akhir:     88,
			Grade:           "A",
			Keterangan:      "Lulus",
			Created_at:      time.Now(),
			Created_by:      "DOS003",
			Updated_at:      time.Now(),
		},
		{
			Kode_matakuliah: "MK004",
			NIM:             "2024030001",
			Kelas_id:        4,
			Tugas:           82,
			UTS:             85,
			UAS:             88,
			Nilai_akhir:     85,
			Grade:           "A",
			Keterangan:      "Lulus",
			Created_at:      time.Now(),
			Created_by:      "DOS004",
			Updated_at:      time.Now(),
		},
	}
	return db.Create(&penilaian).Error
}

// RunAllSeeds runs all seed functions
func SeedAll(db *gorm.DB) error {
	// Run seeds in order based on dependencies
	if err := SeedTahunAkademik(db); err != nil {
		return fmt.Errorf("failed to seed tahun akademik: %v", err)
	}
	if err := SeedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}
	if err := SeedJurusan(db); err != nil {
		return fmt.Errorf("failed to seed jurusan: %v", err)
	}
	if err := SeedProdi(db); err != nil {
		return fmt.Errorf("failed to seed prodi: %v", err)
	}
	if err := SeedStrukturalJurusan(db); err != nil {
		return fmt.Errorf("failed to seed struktural jurusan: %v", err)
	}
	if err := SeedStrukturalProdi(db); err != nil {
		return fmt.Errorf("failed to seed struktural prodi: %v", err)
	}
	if err := SeedStaffAkademik(db); err != nil {
		return fmt.Errorf("failed to seed staff akademik: %v", err)
	}
	if err := SeedDosen(db); err != nil {
		return fmt.Errorf("failed to seed dosen: %v", err)
	}
	if err := SeedKelas(db); err != nil {
		return fmt.Errorf("failed to seed kelas: %v", err)
	}
	if err := SeedMahasiswa(db); err != nil {
		return fmt.Errorf("failed to seed mahasiswa: %v", err)
	}
	if err := SeedKelasMahasiswa(db); err != nil {
		return fmt.Errorf("failed to seed kelas mahasiswa: %v", err)
	}
	if err := SeedMataKuliah(db); err != nil {
		return fmt.Errorf("failed to seed mata kuliah: %v", err)
	}
	if err := SeedKodeRuangan(db); err != nil {
		return fmt.Errorf("failed to seed kode ruangan: %v", err)
	}
	if err := SeedJadwalMatakuliah(db); err != nil {
		return fmt.Errorf("failed to seed jadwal matakuliah: %v", err)
	}
	if err := SeedPertemuan(db); err != nil {
		return fmt.Errorf("failed to seed pertemuan: %v", err)
	}
	if err := SeedAbsensi(db); err != nil {
		return fmt.Errorf("failed to seed absensi: %v", err)
	}
	if err := SeedPenilaian(db); err != nil {
		return fmt.Errorf("failed to seed penilaian: %v", err)
	}

	fmt.Println("✅ All seeding completed successfully")
	return nil
}
