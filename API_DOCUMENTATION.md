# Dokumentasi API Simpadu

## URL Dasar
```
http://localhost:8080
```

## Autentikasi
Sebagian besar endpoint memerlukan autentikasi menggunakan token JWT. Sertakan token dalam header Authorization:
```
Authorization: Bearer <token_jwt_anda>
```

## Endpoint

### Autentikasi

#### Pendaftaran Mahasiswa
- **URL**: `/register-mahasiswa`
- **Method**: `POST`
- **Deskripsi**: Mendaftarkan pengguna mahasiswa baru
- **Body Permintaan**:
  ```json
  {
    "username": "string",
    "password": "string",
    "email": "string",
    "nim": "string"
  }
  ```
- **Respon**:
  ```json
  {
    "message": "User berhasil dibuat",
    "user_id": "string"
  }
  ```

#### Pendaftaran Admin Prodi
- **URL**: `/register-admin-prodi`
- **Method**: `POST`
- **Deskripsi**: Mendaftarkan admin program studi baru
- **Body Permintaan**:
  ```json
  {
    "user_id": "string",
    "nip": "string",
    "jabatan": "string",
    "kode_prodi": "string"
  }
  ```
- **Respon**:
  ```json
  {
    "message": "Admin prodi berhasil disimpan",
    "user_id": "string"
  }
  ```

#### Masuk
- **URL**: `/login`
- **Method**: `POST`
- **Deskripsi**: Autentikasi pengguna dan mendapatkan token JWT
- **Body Permintaan**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Respon**:
  ```json
  {
    "message": "Login berhasil",
    "token": "string"
  }
  ```

### Profil Pengguna

#### Dapatkan Profil
- **URL**: `/profile`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan informasi profil pengguna saat ini
- **Header**: Memerlukan token JWT
- **Respon**:
  ```json
  {
    "status": 200,
    "message": "Berhasil mendapatkan profil",
    "data": {
      "id": "number",
      "user_id": "string",
      "username": "string",
      "email": "string",
      "role": "string",
      "status": "string"
    }
  }
  ```

### Kelas

#### Dapatkan Semua Kelas
- **URL**: `/kelas`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan daftar semua kelas
- **Header**: Memerlukan token JWT
- **Respon**:
  ```json
  {
    "message": "Berhasil Mengambil Data Kelas",
    "data": [
      {
        "id": "number",
        "kode": "string",
        "nama": "string",
        "kode_prodi": "string",
        "tahun_akademik_id": "number",
        "created_by": "string",
        "created_at": "datetime",
        "updated_at": "datetime"
      }
    ]
  }
  ```

#### Buat Kelas
- **URL**: `/create-kelas`
- **Method**: `POST`
- **Deskripsi**: Membuat kelas baru
- **Header**: Memerlukan token JWT
- **Body Permintaan**:
  ```json
  {
    "kode": "string",
    "nama": "string",
    "prodi_id": "number",
    "tahun_akademik_id": "number"
  }
  ```
- **Respon**:
  ```json
  {
    "message": "Kelas Berhasil Dibuat",
    "data": {
      "id": "number",
      "kode": "string",
      "nama": "string",
      "kode_prodi": "string",
      "tahun_akademik_id": "number",
      "created_by": "string",
      "created_at": "datetime",
      "updated_at": "datetime"
    },
    "jumlah_mahasiswa": "number"
  }
  ```

#### Detail Kelas
- **URL**: `/kelas/:id`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan detail kelas berdasarkan ID
- **Header**: Memerlukan token JWT
- **Respon**:
  ```json
  {
    "message": "Berhasil Mengambil Data Kelas",
    "data": {
      "id": "number",
      "kode": "string",
      "nama": "string",
      "kode_prodi": "string",
      "nama_prodi": "string",
      "tahun_akademik": {
        "id": "number",
        "tahun": "string",
        "semester": "string",
        "status": "string"
      },
      "total_mahasiswa": "number",
      "created_by_nama": "string",
      "mahasiswa_list": [
        {
          "nim": "string",
          "nama": "string",
          "angkatan": "number",
          "assigned_at": "datetime",
          "assigned_by": "string"
        }
      ],
      "available_mahasiswa": [
        {
          "nim": "string",
          "nama": "string",
          "angkatan": "number"
        }
      ]
    }
  }
  ```

### Tahun Akademik

#### Dapatkan Semua Tahun Akademik
- **URL**: `/tahun-akademik`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan daftar semua tahun akademik
- **Header**: Memerlukan token JWT
- **Respon**:
  ```json
  {
    "status": 200,
    "message": "Berhasil Mengambil Data Tahun Akademik",
    "data": [
      {
        "id": "number",
        "tahun": "string",
        "semester": "string",
        "status": "string",
        "tanggal_mulai": "datetime",
        "tanggal_selesai": "datetime",
        "created_at": "datetime",
        "updated_at": "datetime"
      }
    ]
  }
  ```

#### Buat Tahun Akademik
- **URL**: `/create-tahun-akademik`
- **Method**: `POST`
- **Deskripsi**: Membuat tahun akademik baru
- **Header**: Memerlukan token JWT
- **Body Permintaan**:
  ```json
  {
    "tahun": "string",
    "semester": "string",
    "status": "string",
    "tanggal_mulai": "datetime",
    "tanggal_selesai": "datetime"
  }
  ```
- **Respon**:
  ```json
  {
    "status": 200,
    "message": "Berhasil Membuat Tahun Akademik",
    "data": {
      "id": "number",
      "tahun": "string",
      "semester": "string",
      "status": "string",
      "tanggal_mulai": "datetime",
      "tanggal_selesai": "datetime",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  }
  ```

#### Dapatkan Tahun Akademik berdasarkan ID
- **URL**: `/tahun-akademik/:id`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan tahun akademik tertentu berdasarkan ID
- **Header**: Memerlukan token JWT
- **Respon**:
  ```json
  {
    "message": "Berhasil Mengambil Data Tahun Akademik",
    "data": {
      "id": "number",
      "tahun": "string",
      "semester": "string",
      "status": "string",
      "tanggal_mulai": "datetime",
      "tanggal_selesai": "datetime",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  }
  ```

### Mata Kuliah

#### Dapatkan Semua Mata Kuliah
- **URL**: `/matakuliah`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan daftar semua mata kuliah
- **Header**: Memerlukan token JWT
- **Respon**:
  ```json
  {
    "data": [
      {
        "id": "number",
        "kode_matakuliah": "string",
        "nama_matakuliah": "string",
        "kode_prodi": "string",
        "sks": "number",
        "semester": "number",
        "status": "string",
        "created_at": "datetime",
        "created_by": "string",
        "updated_at": "datetime"
      }
    ]
  }
  ```

#### Buat Mata Kuliah
- **URL**: `/create-matakuliah`
- **Method**: `POST`
- **Deskripsi**: Membuat mata kuliah baru
- **Header**: Memerlukan token JWT
- **Body Permintaan**:
  ```json
  {
    "kode_matakuliah": "string",
    "nama_matakuliah": "string",
    "sks": "number",
    "semester": "number",
    "kode_prodi": "string"
  }
  ```
- **Respon**:
  ```json
  {
    "message": "Matakuliah Berhasil Dibuat"
  }
  ```

#### Detail Mata Kuliah
- **URL**: `/matakuliah/:kode_matakuliah`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan detail mata kuliah berdasarkan kode
- **Header**: Memerlukan token JWT
- **Respon**:
  ```json
  {
    "data": {
      "id": "number",
      "kode_matakuliah": "string",
      "nama_matakuliah": "string",
      "kode_prodi": "string",
      "sks": "number",
      "semester": "number",
      "status": "string",
      "created_at": "datetime",
      "created_by": "string",
      "updated_at": "datetime"
    }
  }
  ```

## Jadwal Mata Kuliah

### 1. Membuat Jadwal Baru
```http
POST /jadwal/
Authorization: Bearer <token>
Content-Type: application/json

{
    "kode_matakuliah": "TI101",
    "kelas_id": 1,
    "kode_dosen": "DSN001",
    "kode_ruangan": "A-LAB-1",
    "kode_prodi": "TI",
    "jam_mulai": "2024-03-18T08:00:00Z",
    "jam_selesai": "2024-03-18T09:30:00Z",
    "hari": "Senin",
    "kode_tahun_akademik": "2023-2024-Ganjil"
}
```

Response Sukses (200):
```json
{
    "message": "Jadwal berhasil dibuat",
    "data": {
        "id": 1,
        "kode_matakuliah": "TI101",
        "kelas_id": 1,
        "kode_dosen": "DSN001",
        "kode_ruangan": "A-LAB-1",
        "kode_prodi": "TI",
        "jam_mulai": "2024-03-18T08:00:00Z",
        "jam_selesai": "2024-03-18T09:30:00Z",
        "hari": "Senin",
        "kode_tahun_akademik": "2023-2024-Ganjil",
        "created_at": "2024-03-18T07:00:00Z",
        "created_by": "ADM001"
    }
}
```

### 2. Melihat Jadwal per Kelas
```http
GET /jadwal/kelas/1
```

Response Sukses (200):
```json
{
    "data": [
        {
            "id": 1,
            "kode_matakuliah": "TI101",
            "kelas_id": 1,
            "kode_dosen": "DSN001",
            "kode_ruangan": "A-LAB-1",
            "jam_mulai": "2024-03-18T08:00:00Z",
            "jam_selesai": "2024-03-18T09:30:00Z",
            "hari": "Senin"
        }
    ]
}
```

### 3. Melihat Jadwal per Dosen
```http
GET /jadwal/dosen/DSN001
```

Response Sukses (200):
```json
{
    "data": [
        {
            "id": 1,
            "kode_matakuliah": "TI101",
            "kelas_id": 1,
            "kode_dosen": "DSN001",
            "kode_ruangan": "A-LAB-1",
            "jam_mulai": "2024-03-18T08:00:00Z",
            "jam_selesai": "2024-03-18T09:30:00Z",
            "hari": "Senin"
        }
    ]
}
```

### 4. Update Jadwal
```http
PUT /jadwal/1
Authorization: Bearer <token>
Content-Type: application/json

{
    "kode_ruangan": "A-LAB-2",
    "jam_mulai": "2024-03-18T10:00:00Z",
    "jam_selesai": "2024-03-18T11:30:00Z",
    "hari": "Senin"
}
```

Response Sukses (200):
```json
{
    "message": "Jadwal berhasil diupdate"
}
```

### 5. Hapus Jadwal
```http
DELETE /jadwal/1
Authorization: Bearer <token>
```

Response Sukses (200):
```json
{
    "message": "Jadwal berhasil dihapus"
}
```

### 6. Cek Ruangan Tersedia
```http
POST /jadwal/ruangan-tersedia
Content-Type: application/json

{
    "hari": "Senin",
    "jam_mulai": "2024-03-18T08:00:00Z",
    "jam_selesai": "2024-03-18T09:30:00Z",
    "kode_prodi": "TI"
}
```

Response Sukses (200):
```json
{
    "data": [
        {
            "kode_ruangan": "A-LAB-1",
            "nama_ruangan": "Laboratorium Komputer 1",
            "kapasitas": 30,
            "lokasi": "Lantai 1",
            "status": "Aktif"
        }
    ]
}
```

## Manajemen Ruangan

### 1. Membuat Ruangan Baru
```http
POST /ruangan/
Authorization: Bearer <token>
Content-Type: application/json

{
    "gedung": "A",
    "tipe": "LAB",
    "nomor": 1,
    "nama_ruangan": "Laboratorium Komputer 1",
    "kode_prodi": "TI",
    "kapasitas": 30,
    "lokasi": "Lantai 1"
}
```

Response Sukses (200):
```json
{
    "message": "Ruangan berhasil dibuat",
    "data": {
        "kode_ruangan": "A-LAB-1",
        "nama_ruangan": "Laboratorium Komputer 1",
        "kode_prodi": "TI",
        "gedung": "A",
        "tipe": "LAB",
        "nomor": 1,
        "kapasitas": 30,
        "lokasi": "Lantai 1",
        "status": "Aktif"
    }
}
```

### 2. Melihat Ruangan per Gedung
```http
GET /ruangan/gedung/A
```

Response Sukses (200):
```json
{
    "data": [
        {
            "kode_ruangan": "A-LAB-1",
            "nama_ruangan": "Laboratorium Komputer 1",
            "kode_prodi": "TI",
            "gedung": "A",
            "tipe": "LAB",
            "nomor": 1,
            "kapasitas": 30,
            "lokasi": "Lantai 1",
            "status": "Aktif"
        }
    ]
}
```

### 3. Melihat Ruangan per Prodi
```http
GET /ruangan/prodi/TI
```

Response Sukses (200):
```json
{
    "data": [
        {
            "kode_ruangan": "A-LAB-1",
            "nama_ruangan": "Laboratorium Komputer 1",
            "kode_prodi": "TI",
            "gedung": "A",
            "tipe": "LAB",
            "nomor": 1,
            "kapasitas": 30,
            "lokasi": "Lantai 1",
            "status": "Aktif"
        }
    ]
}
```

### 4. Melihat Ruangan per Tipe
```http
GET /ruangan/tipe/LAB
```

Response Sukses (200):
```json
{
    "data": [
        {
            "kode_ruangan": "A-LAB-1",
            "nama_ruangan": "Laboratorium Komputer 1",
            "kode_prodi": "TI",
            "gedung": "A",
            "tipe": "LAB",
            "nomor": 1,
            "kapasitas": 30,
            "lokasi": "Lantai 1",
            "status": "Aktif"
        }
    ]
}
```

### 5. Update Ruangan
```http
PUT /ruangan/A-LAB-1
Authorization: Bearer <token>
Content-Type: application/json

{
    "nama_ruangan": "Laboratorium Komputer 1",
    "kapasitas": 35,
    "lokasi": "Lantai 1",
    "status": "Aktif"
}
```

Response Sukses (200):
```json
{
    "message": "Ruangan berhasil diupdate"
}
```

### 6. Hapus Ruangan
```http
DELETE /ruangan/A-LAB-1
Authorization: Bearer <token>
```

Response Sukses (200):
```json
{
    "message": "Ruangan berhasil dihapus"
}
```

### 7. Cek Status Ruangan
```http
GET /ruangan/status/A-LAB-1?hari=Senin
```

Response Sukses (200):
```json
{
    "ruangan": {
        "kode_ruangan": "A-LAB-1",
        "nama_ruangan": "Laboratorium Komputer 1",
        "kode_prodi": "TI",
        "gedung": "A",
        "tipe": "LAB",
        "nomor": 1,
        "kapasitas": 30,
        "lokasi": "Lantai 1",
        "status": "Aktif"
    },
    "jadwal": [
        {
            "id": 1,
            "kode_matakuliah": "TI101",
            "kelas_id": 1,
            "kode_dosen": "DSN001",
            "jam_mulai": "2024-03-18T08:00:00Z",
            "jam_selesai": "2024-03-18T09:30:00Z",
            "hari": "Senin"
        }
    ]
}
```

## Absensi

### 1. Buka Pertemuan
```http
POST /buka-pertemuan
Authorization: Bearer <token>
Content-Type: application/json

{
    "kode_matakuliah": "TI101",
    "kelas_id": 1,
    "pertemuan_ke": 1,
    "tanggal": "2024-03-18",
    "jam_mulai": "2024-03-18T08:00:00Z",
    "jam_selesai": "2024-03-18T09:30:00Z",
    "kode_ruangan": "A-LAB-1"
}
```

Response Sukses (200):
```json
{
    "message": "Pertemuan berhasil dibuka",
    "data": {
        "kode_pertemuan": "PTM20240318001",
        "kode_matakuliah": "TI101",
        "kelas_id": 1,
        "pertemuan_ke": 1,
        "tanggal": "2024-03-18",
        "jam_mulai": "2024-03-18T08:00:00Z",
        "jam_selesai": "2024-03-18T09:30:00Z",
        "kode_ruangan": "A-LAB-1",
        "status": "Aktif",
        "created_by": "DSN001"
    }
}
```

### 2. Tutup Pertemuan
```http
POST /tutup-pertemuan/PTM20240318001
Authorization: Bearer <token>
```

Response Sukses (200):
```json
{
    "message": "Pertemuan berhasil ditutup",
    "data": {
        "kode_pertemuan": "PTM20240318001",
        "status": "Selesai",
        "jumlah_hadir": 25,
        "jumlah_tidak_hadir": 5
    }
}
```

### 3. Absen
```http
POST /absen
Authorization: Bearer <token>
Content-Type: application/json

{
    "kode_pertemuan": "PTM20240318001",
    "nim": "2021001",
    "status": "Hadir",
    "keterangan": "Hadir tepat waktu"
}
```

Response Sukses (200):
```json
{
    "message": "Absensi berhasil dicatat",
    "data": {
        "id": 1,
        "kode_pertemuan": "PTM20240318001",
        "nim": "2021001",
        "status": "Hadir",
        "keterangan": "Hadir tepat waktu",
        "waktu_absen": "2024-03-18T08:05:00Z"
    }
}
```

### 4. Cek Status Pertemuan
```http
GET /status-pertemuan/TI101
Authorization: Bearer <token>
```

Response Sukses (200):
```json
{
    "data": {
        "kode_pertemuan": "PTM20240318001",
        "kode_matakuliah": "TI101",
        "kelas_id": 1,
        "pertemuan_ke": 1,
        "tanggal": "2024-03-18",
        "jam_mulai": "2024-03-18T08:00:00Z",
        "jam_selesai": "2024-03-18T09:30:00Z",
        "kode_ruangan": "A-LAB-1",
        "status": "Aktif",
        "jumlah_hadir": 25,
        "jumlah_tidak_hadir": 5,
        "created_by": "DSN001"
    }
}
```

### 5. Rekap Absensi
```http
GET /rekap-absensi/TI101
Authorization: Bearer <token>
```

Response Sukses (200):
```json
{
    "data": {
        "kode_matakuliah": "TI101",
        "nama_matakuliah": "Pemrograman Dasar",
        "kelas_id": 1,
        "nama_kelas": "TI-1A",
        "dosen": "Dr. John Doe",
        "total_pertemuan": 14,
        "rekap_mahasiswa": [
            {
                "nim": "2021001",
                "nama": "Jane Doe",
                "total_hadir": 12,
                "total_tidak_hadir": 2,
                "persentase_kehadiran": 85.71,
                "detail_pertemuan": [
                    {
                        "pertemuan_ke": 1,
                        "tanggal": "2024-03-18",
                        "status": "Hadir",
                        "keterangan": "Hadir tepat waktu"
                    }
                ]
            }
        ]
    }
}
```

### Catatan Penting Absensi

1. **Status Pertemuan**:
   - Aktif: Pertemuan sedang berlangsung
   - Selesai: Pertemuan telah ditutup
   - Batal: Pertemuan dibatalkan

2. **Status Kehadiran**:
   - Hadir: Mahasiswa hadir dalam pertemuan
   - Tidak Hadir: Mahasiswa tidak hadir
   - Izin: Mahasiswa hadir dengan izin
   - Sakit: Mahasiswa tidak hadir karena sakit

3. **Validasi**:
   - Hanya dosen pengajar yang dapat membuka/menutup pertemuan
   - Absensi hanya dapat dilakukan saat pertemuan aktif
   - Mahasiswa hanya dapat absen untuk kelas yang diikutinya
   - Tidak dapat membuka pertemuan baru jika ada pertemuan yang masih aktif
   - Waktu absen tidak boleh melebihi jam selesai pertemuan

4. **Format Kode Pertemuan**:
   - Format: `PTM<YYYYMMDD><XXX>`
   - Contoh: PTM20240318001 (Pertemuan tanggal 18 Maret 2024, nomor urut 001)

5. **Rekap Absensi**:
   - Menampilkan statistik kehadiran per mahasiswa
   - Menghitung persentase kehadiran
   - Menampilkan detail kehadiran per pertemuan
   - Dapat diakses oleh dosen dan admin prodi

## Penilaian

### 1. Input Nilai
```http
POST /penilaian/
Authorization: Bearer <token>
Content-Type: application/json

{
    "kode_matakuliah": "TI101",
    "kelas_id": 1,
    "nim": "2021001",
    "tugas": 85,
    "uts": 80,
    "uas": 90,
    "keterangan": "Nilai tugas sudah termasuk bonus"
}
```

Response Sukses (200):
```json
{
    "message": "Penilaian berhasil disimpan",
    "data": {
        "id": 1,
        "kode_matakuliah": "TI101",
        "kelas_id": 1,
        "nim": "2021001",
        "tugas": 85,
        "uts": 80,
        "uas": 90,
        "nilai_akhir": 85.5,
        "grade": "A",
        "keterangan": "Nilai tugas sudah termasuk bonus",
        "created_at": "2024-03-18T08:00:00Z",
        "created_by": "DSN001"
    }
}
```

### 2. Lihat Nilai per Kelas
```http
GET /penilaian/kelas/TI101/1
```

Response Sukses (200):
```json
{
    "data": [
        {
            "id": 1,
            "kode_matakuliah": "TI101",
            "kelas_id": 1,
            "nim": "2021001",
            "tugas": 85,
            "uts": 80,
            "uas": 90,
            "nilai_akhir": 85.5,
            "grade": "A",
            "keterangan": "Nilai tugas sudah termasuk bonus",
            "created_at": "2024-03-18T08:00:00Z",
            "created_by": "DSN001"
        }
    ]
}
```

### 3. Lihat Nilai per Mahasiswa
```http
GET /penilaian/mahasiswa/2021001
```

Response Sukses (200):
```json
{
    "data": [
        {
            "id": 1,
            "kode_matakuliah": "TI101",
            "kelas_id": 1,
            "nim": "2021001",
            "tugas": 85,
            "uts": 80,
            "uas": 90,
            "nilai_akhir": 85.5,
            "grade": "A",
            "keterangan": "Nilai tugas sudah termasuk bonus",
            "created_at": "2024-03-18T08:00:00Z",
            "created_by": "DSN001"
        }
    ]
}
```

### 4. Update Nilai
```http
PUT /penilaian/1
Authorization: Bearer <token>
Content-Type: application/json

{
    "tugas": 90,
    "uts": 85,
    "uas": 95,
    "keterangan": "Nilai tugas sudah termasuk bonus dan revisi"
}
```

Response Sukses (200):
```json
{
    "message": "Penilaian berhasil diupdate"
}
```

### 5. Rekap Nilai
```http
GET /penilaian/rekap/TI101/1
```

Response Sukses (200):
```json
{
    "data": {
        "kode_matakuliah": "TI101",
        "nama_matakuliah": "Pemrograman Dasar",
        "kelas_id": 1,
        "nama_kelas": "TI-1A",
        "total_mahasiswa": 30,
        "rata_rata": 78.5,
        "distribusi": {
            "A": 5,
            "A-": 3,
            "B+": 7,
            "B": 8,
            "B-": 4,
            "C+": 2,
            "C": 1
        },
        "detail_nilai": [
            {
                "id": 1,
                "kode_matakuliah": "TI101",
                "kelas_id": 1,
                "nim": "2021001",
                "tugas": 85,
                "uts": 80,
                "uas": 90,
                "nilai_akhir": 85.5,
                "grade": "A",
                "keterangan": "Nilai tugas sudah termasuk bonus"
            }
        ]
    }
}
```

### Catatan Penting Penilaian

1. **Komponen Nilai**:
   - Tugas: 30% dari nilai akhir
   - UTS: 30% dari nilai akhir
   - UAS: 40% dari nilai akhir

2. **Skala Nilai**:
   - A: ≥ 85
   - A-: ≥ 80
   - B+: ≥ 75
   - B: ≥ 70
   - B-: ≥ 65
   - C+: ≥ 60
   - C: ≥ 55
   - C-: ≥ 50
   - D: ≥ 40
   - E: < 40

3. **Validasi**:
   - Hanya dosen pengajar yang dapat input/update nilai
   - Admin prodi dapat melihat semua nilai
   - Mahasiswa hanya dapat melihat nilai mereka sendiri
   - Nilai harus antara 0-100
   - Mahasiswa harus terdaftar di kelas yang bersangkutan

4. **Rekap Nilai**:
   - Menampilkan statistik nilai per kelas
   - Menghitung rata-rata kelas
   - Menampilkan distribusi nilai
   - Menampilkan detail nilai per mahasiswa

## Respon Error

API menggunakan kode status HTTP standar:

- `200 OK`: Permintaan berhasil
- `400 Bad Request`: Format permintaan tidak valid
- `401 Unauthorized`: Autentikasi diperlukan atau gagal
- `403 Forbidden`: Tidak memiliki izin yang cukup
- `404 Not Found`: Sumber daya tidak ditemukan
- `409 Conflict`: Sumber daya sudah ada
- `500 Internal Server Error`: Kesalahan server

Respon error mengikuti format berikut:
```json
{
  "error": "Deskripsi pesan error"
}
```

## Keamanan

- Semua endpoint sensitif memerlukan autentikasi JWT
- Password di-hash menggunakan bcrypt
- Token JWT ditandatangani menggunakan algoritma HS256
- Variabel lingkungan `simpadu_jwt_key` digunakan untuk penandatanganan JWT (disarankan 32+ bytes)
- Token JWT memiliki masa berlaku 1 jam

## Catatan Tambahan

### Format Data
- Semua tanggal menggunakan format ISO 8601
- Semester hanya menerima nilai "Ganjil" atau "Genap"
- Status hanya menerima nilai "Aktif" atau "Tidak Aktif"

### Validasi
- Username dan email harus unik
- Password minimal 6 karakter
- NIM harus terdaftar di sistem
- Tanggal mulai tidak boleh lebih besar dari tanggal selesai
- Kode mata kuliah harus unik
- SKS harus berupa angka positif
- Semester mata kuliah harus antara 1-8 

### Catatan Penting

1. Format waktu menggunakan ISO 8601 (UTC)
2. Format kode ruangan: `GEDUNG-TIPE-NOMOR` (contoh: A-LAB-1, B-KLS-2)
3. Tipe ruangan yang tersedia: LAB (Laboratorium), KLS (Kelas), STU (Studio)
4. Status ruangan: Aktif, Maintenance, Tidak Aktif
5. Validasi jadwal:
   - Jam mulai harus sebelum jam selesai
   - Tidak boleh bentrok dengan jadwal lain di ruangan yang sama
   - Dosen tidak boleh memiliki jadwal yang bentrok
   - Kelas tidak boleh memiliki jadwal yang bentrok
6. Ruangan tidak dapat dihapus jika masih digunakan dalam jadwal 