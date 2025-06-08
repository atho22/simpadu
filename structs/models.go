package structs

import "time"

type Tahun_akademik struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Kode_tahun_akademik string    `gorm:"size:20;not null" json:"kode_tahun_akademik"`
	Tahun               string    `gorm:"size:9;not null" json:"tahun"`
	Semester            string    `gorm:"type:enum('Ganjil','Genap');not null" json:"semester"`
	Status              string    `gorm:"type:enum('Aktif','Tidak Aktif');default:'Tidak Aktif'" json:"status"`
	Tanggal_mulai       time.Time `gorm:"not null" json:"tanggal_mulai"`
	Tanggal_selesai     time.Time `gorm:"not null" json:"tanggal_selesai"`
	Created_at          time.Time `json:"created_at"`
	Updated_at          time.Time `json:"updated_at"`
}

type Users struct {
	ID       uint   `gorm:"primaryKey"`
	User_id  string `gorm:"size:20;not null;unique"`
	Username string `gorm:"size:50;not null;unique"`
	Email    string `gorm:"size:100;not null;unique"`
	Password string `gorm:"size:255;not null"` // hashed with bcrypt
	Role     string `gorm:"size:20;not null"`  // e.g. "admin_prodi", "mahasiswa", "dosen"
	Status   string `gorm:"type:enum('Aktif','Tidak Aktif');default:'Aktif'"`
}

type Jurusan struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Kode_jurusan string    `gorm:"size:20;not null" json:"kode_jurusan"`
	Nama_jurusan string    `gorm:"size:50;not null" json:"nama_jurusan"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

type Prodi struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Kode_prodi   string    `gorm:"size:20;not null" json:"kode_prodi"`
	Nama_prodi   string    `gorm:"size:50;not null" json:"nama_prodi"`
	Kode_jurusan string    `gorm:"size:20;not null" json:"kode_jurusan"`
	Jenjang      string    `gorm:"size:20;not null" json:"jenjang"`
	User_id      string    `gorm:"size:20" json:"user_id"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}
type Struktural_prodi struct {
	ID         uint   `gorm:"primaryKey"`
	User_id    string `gorm:"size:20;not null;unique"` // relasi ke users
	Kode_prodi string `gorm:"size:20;not null"`        // relasi ke prodi
	Nama       string `gorm:"size:100;not null"`
	N_ip       string `gorm:"size:30;unique"` // diubah dari NIP ke n_ip
	Jabatan    string `gorm:"size:50"`        // opsional: ketua, sekretaris , admin prodi
	NoSK       string `gorm:"size:50"`
	TanggalSK  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Dosen struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	User_id    string    `gorm:"size:20;not null" json:"user_id"`
	Kode_dosen string    `gorm:"size:20;not null" json:"kode_dosen"`
	NIP        string    `gorm:"size:20" json:"nip"`
	Nidn       string    `gorm:"size:20" json:"nidn"`
	Kode_prodi string    `gorm:"size:20;not null" json:"kode_prodi"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type Kelas struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Kode_kelas          string    `gorm:"size:10;not null" json:"kode_kelas"`
	Nama                string    `gorm:"size:50;not null" json:"nama"`
	Kode_prodi          string    `gorm:"size:20;not null" json:"kode_prodi"`
	Kode_tahun_akademik string    `gorm:"size:20;not null" json:"kode_tahun_akademik"`
	Created_by          string    `gorm:"size:20;not null" json:"created_by"`
	Created_at          time.Time `json:"created_at"`
	Updated_at          time.Time `json:"updated_at"`
}

type Mahasiswa struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	User_id    string    `gorm:"size:20;not null" json:"user_id"`
	NIM        string    `gorm:"size:20;not null" json:"nim"`
	Nama       string    `gorm:"size:100;not null" json:"nama"`
	Angkatan   uint      `gorm:"not null" json:"angkatan"`
	Kode_prodi string    `gorm:"size:20;not null" json:"kode_prodi"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type Kelas_mahasiswa struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Kelas_id   uint      `gorm:"not null" json:"kelas_id"`
	NIM        string    `gorm:"size:20;not null" json:"nim"`
	Created_by string    `gorm:"size:20;not null" json:"created_by"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type Mata_kuliah struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Kode_matakuliah     string    `gorm:"size:20;not null" json:"kode_matakuliah"`
	Nama_matakuliah     string    `gorm:"size:100;not null" json:"nama_matakuliah"`
	Kode_prodi          string    `gorm:"size:20;not null" json:"kode_prodi"`
	SKS                 int       `gorm:"not null" json:"sks"`
	Semester            int       `gorm:"not null" json:"semester"`
	Kode_tahun_akademik string    `gorm:"not null" json:"kode_tahun_akademik"`
	Status              string    `gorm:"type:enum('Aktif','Tidak Aktif');default:'Aktif'" json:"status"`
	Created_at          time.Time `json:"created_at"`
	Created_by          string    `gorm:"size:20;not null" json:"created_by"`
	Updated_at          time.Time `json:"updated_at"`
}
type Jadwal_matakuliah struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Kode_matakuliah     string    `gorm:"size:20;not null" json:"kode_matakuliah"`
	Kelas_id            uint      `gorm:"not null" json:"kelas_id"`
	Kode_dosen          string    `gorm:"size:20;not null" json:"kode_dosen"`
	Kode_ruangan        string    `gorm:"size:20;not null" json:"kode_ruangan"`
	Kode_prodi          string    `gorm:"size:20;not null" json:"kode_prodi"`
	Jam_mulai           time.Time `gorm:"not null" json:"jam_mulai"`
	Jam_selesai         time.Time `gorm:"not null" json:"jam_selesai"`
	Hari                string    `gorm:"size:20;not null" json:"hari"`
	Kode_tahun_akademik string    `gorm:"size:20;not null" json:"kode_tahun_akademik"`
	Created_at          time.Time `json:"created_at"`
	Created_by          string    `gorm:"size:20;not null" json:"created_by"`
}
type Kode_ruangan struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Kode_ruangan string    `gorm:"size:20;not null;unique" json:"kode_ruangan"` // Format: GEDUNG-TIPE-NOMOR (contoh: A-LAB-1, B-KLS-2)
	Nama_ruangan string    `gorm:"size:100;not null" json:"nama_ruangan"`
	Kode_prodi   string    `gorm:"size:20;not null" json:"kode_prodi"`
	Gedung       string    `gorm:"size:20;not null" json:"gedung"` // A, B, C, dll
	Tipe         string    `gorm:"size:20;not null" json:"tipe"`   // LAB, KLS, STU, dll
	Nomor        int       `gorm:"not null" json:"nomor"`          // 1, 2, 3, dll
	Kapasitas    int       `gorm:"not null" json:"kapasitas"`
	Lokasi       string    `gorm:"size:100;not null" json:"lokasi"` // Lantai 1, Lantai 2, dll
	Status       string    `gorm:"type:enum('Aktif','Maintenance','Tidak Aktif');default:'Aktif'" json:"status"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

type Pertemuan struct {
	ID              int       `gorm:"primaryKey" json:"id"`
	Kode_pertemuan  string    `gorm:"size:20;not null" json:"kode_pertemuan"`
	Kode_matakuliah string    `gorm:"size:20;not null" json:"kode_matakuliah"`
	Pertemuan_ke    int       `gorm:"not null" json:"pertemuan_ke"`
	Dibuka          bool      `gorm:"not null" json:"dibuka"`
	Waktu_dibuka    time.Time `gorm:"not null" json:"waktu_dibuka"`
	Waktu_ditutup   time.Time `gorm:"not null" json:"waktu_ditutup"`
	Created_at      time.Time `json:"created_at"`
	Updated_at      time.Time `json:"updated_at"`
	Dibuka_by       string    `gorm:"size:20;not null" json:"dibuka_by"`
}
type Absensi struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Kode_matakuliah string    `gorm:"size:20;not null" json:"kode_matakuliah"`
	Kode_pertemuan  string    `gorm:"size:20;not null" json:"kode_pertemuan"`
	Nim             string    `gorm:"size:20;not null" json:"nim"`
	Status          string    `gorm:"type:enum('hadir','izin','sakit','alfa');default:'alfa'" json:"status"`
	Keterangan      string    `gorm:"size:255" json:"keterangan"`
	Created_at      time.Time `json:"created_at"`
}

// =================================================================
// RESPONSE STRUCTS (untuk API responses)
// =================================================================

// Response untuk Admin Akademik - Overview semua prodi dan kelas
type Dashboard_admin_akademik struct {
	Total_prodi     int             `json:"total_prodi"`
	Total_kelas     int             `json:"total_kelas"`
	Total_mahasiswa int             `json:"total_mahasiswa"`
	Tahun_aktif     Tahun_akademik  `json:"tahun_aktif"`
	Prodi_list      []Prodi_summary `json:"prodi_list"`
}

type Prodi_summary struct {
	Kode_prodi      string `json:"kode_prodi"`
	Nama_prodi      string `json:"nama_prodi"`
	Nama_jurusan    string `json:"nama_jurusan"`
	Admin_nama      string `json:"admin_nama"`
	Total_kelas     int    `json:"total_kelas"`
	Total_mahasiswa int    `json:"total_mahasiswa"`
}

// Response untuk Admin Prodi - Detail kelas dan mahasiswa di prodi mereka
type Dashboard_admin_prodi struct {
	Prodi_info      Prodi_detail    `json:"prodi_info"`
	Total_kelas     int             `json:"total_kelas"`
	Total_mahasiswa int             `json:"total_mahasiswa"`
	Tahun_aktif     Tahun_akademik  `json:"tahun_aktif"`
	Kelas_list      []Kelas_summary `json:"kelas_list"`
}

type Prodi_detail struct {
	Kode_prodi   string `json:"kode_prodi"`
	Nama_prodi   string `json:"nama_prodi"`
	Nama_jurusan string `json:"nama_jurusan"`
	Admin_nama   string `json:"admin_nama"`
	Admin_email  string `json:"admin_email"`
}

type Kelas_summary struct {
	ID              uint   `json:"id"`
	Kode            string `json:"kode"`
	Nama            string `json:"nama"`
	Tahun           string `json:"tahun"`
	Semester        string `json:"semester"`
	Total_mahasiswa int    `json:"total_mahasiswa"`
	Created_by_nama string `json:"created_by_nama"`
}

// Response detail kelas untuk Admin Prodi (ketika melihat/edit kelas)
type Kelas_detail_response struct {
	ID                  uint                  `json:"id"`
	Kode                string                `json:"kode"`
	Nama                string                `json:"nama"`
	Kode_prodi          string                `json:"kode_prodi"`
	Nama_prodi          string                `json:"nama_prodi"`
	Tahun_akademik      Tahun_akademik_simple `json:"tahun_akademik"`
	Total_mahasiswa     int                   `json:"total_mahasiswa"`
	Created_by_nama     string                `json:"created_by_nama"`
	Mahasiswa_list      []Mahasiswa_in_kelas  `json:"mahasiswa_list"`
	Available_mahasiswa []Mahasiswa_simple    `json:"available_mahasiswa,omitempty"` // untuk add mahasiswa
}

type Tahun_akademik_simple struct {
	ID       uint   `json:"id"`
	Tahun    string `json:"tahun"`
	Semester string `json:"semester"`
	Status   string `json:"status"`
}

type Mahasiswa_in_kelas struct {
	NIM         string    `json:"nim"`
	Nama        string    `json:"nama"`
	Angkatan    uint      `json:"angkatan"`
	Assigned_at time.Time `json:"assigned_at"`
	Assigned_by string    `json:"assigned_by_nama"`
}

type Mahasiswa_simple struct {
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Angkatan uint   `json:"angkatan"`
}

// Response untuk list mahasiswa (untuk assign ke kelas)
type Mahasiswa_list_response struct {
	Total_count  int                `json:"total_count"`
	Current_page int                `json:"current_page"`
	Per_page     int                `json:"per_page"`
	Mahasiswa    []Mahasiswa_simple `json:"mahasiswa"`
}

// Response untuk user profile dan role info
type User_profile_response struct {
	User_id  string             `json:"user_id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Nama     string             `json:"nama"`
	Status   string             `json:"status"`
	Roles    []User_role_detail `json:"roles"`
}

type User_role_detail struct {
	Role_code  string `json:"role_code"`
	Role_name  string `json:"role_name"`
	Kode_prodi string `json:"kode_prodi,omitempty"`
	Nama_prodi string `json:"nama_prodi,omitempty"`
	Status     string `json:"status"`
}

// Response untuk create/update operations
type Create_kelas_response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Kelas_id uint `json:"kelas_id"`
	} `json:"data,omitempty"`
}

type Assign_mahasiswa_response struct {
	Success        bool     `json:"success"`
	Message        string   `json:"message"`
	Assigned_count int      `json:"assigned_count"`
	Failed_nims    []string `json:"failed_nims,omitempty"`
}

// Penilaian represents a student's grade for a course
type Penilaian struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Kode_matakuliah string    `gorm:"size:20;not null" json:"kode_matakuliah"`
	Kelas_id        uint      `gorm:"not null" json:"kelas_id"`
	NIM             string    `gorm:"size:20;not null" json:"nim"`
	Tugas           float64   `gorm:"not null" json:"tugas"`
	UTS             float64   `gorm:"not null" json:"uts"`
	UAS             float64   `gorm:"not null" json:"uas"`
	Nilai_akhir     float64   `gorm:"not null" json:"nilai_akhir"`
	Grade           string    `gorm:"size:2;not null" json:"grade"`
	Keterangan      string    `gorm:"size:255" json:"keterangan"`
	Created_at      time.Time `gorm:"not null" json:"created_at"`
	Created_by      string    `gorm:"size:20;not null" json:"created_by"`
	Updated_at      time.Time `json:"updated_at"`
	Updated_by      string    `gorm:"size:20" json:"updated_by"`
}

// Response untuk IPK dan IPS mahasiswa
type IPKResponse struct {
	NIM        string          `json:"nim"`
	Nama       string          `json:"nama"`
	Kode_prodi string          `json:"kode_prodi"`
	Nama_prodi string          `json:"nama_prodi"`
	Angkatan   uint            `json:"angkatan"`
	IPK        float64         `json:"ipk"`
	Total_sks  int             `json:"total_sks"`
	Semester   []SemesterNilai `json:"semester"`
}

type SemesterNilai struct {
	Tahun_akademik string    `json:"tahun_akademik"`
	Semester       string    `json:"semester"`
	IPS            float64   `json:"ips"`
	Total_sks      int       `json:"total_sks"`
	Mata_kuliah    []NilaiMK `json:"mata_kuliah"`
}

type NilaiMK struct {
	Kode_matakuliah string  `json:"kode_matakuliah"`
	Nama_matakuliah string  `json:"nama_matakuliah"`
	SKS             int     `json:"sks"`
	Nilai_akhir     float64 `json:"nilai_akhir"`
	Grade           string  `json:"grade"`
}
