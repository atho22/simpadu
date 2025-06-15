# API Documentation - kelompok1

## Informasi Umum
- **Base URL** : `https://ti054c01.agussbn.my.id`
- **Format** : JSON (kecuali upload menggunakan form-data)
- **Autentikasi** : JWT via header Authorization & Role
- **File Upload** : tersedia di endpoint tertentu, file disimpan di folder `/uploads`
- **Response Autentikasi**: Untuk setiap endpoint yang memiliki autentikasi akan melalui 2 buah middleware autentikasi.

Jika **tidak ada token** akan mengembalikan response code (401).
Jika **token tidak valid** akan mengembalikan response code (403).
Jika **endpoint diakses oleh yang bukan rolenya** akan mengembalikan response json:

```json
{
  "message": "Access Denied!"
}
```

---

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
- **Header**: Memerlukan token JWT (hanya admin akademik).
- **Body Permintaan**:
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string",
    "nama": "string",
    "nip": "string",
    "prodi": "string",
    "jabatan": "string",
    "no_sk": "string",
    "tanggal_sk": "datetime"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Registrasi berhasil",
    "data": {
      "user_id": "string",
      "username": "string",
      "email": "string",
      "nama": "string",
      "prodi": "string",
      "nip": "string",
      "jabatan": "string"
    }
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
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Login berhasil",
    "token": "string",
    "user": {
      "user_id": "string",
      "email": "string",
      "role": "string",
      "nama": "string"
    }
  }
  ```

### Profil Pengguna

#### Dapatkan Profil
- **URL**: `/profile`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan informasi profil pengguna saat ini
- **Header**: Memerlukan token JWT
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Profile berhasil diambil",
    "data": {
      "user_id": "string",
      "username": "string",
      "email": "string",
      "status": "string",
      "nama": "string",
      "roles": [
        {
          "role_code": "string",
          "role_name": "string",
          "kode_prodi": "string",
          "status": "string"
        }
      ]
    }
  }
  ```

### Kelas

#### Dapatkan Semua Kelas
- **URL**: `/kelas`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan daftar semua kelas.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Berhasil Mengambil Data Kelas",
    "data": [
      {
        "id": "number",
        "kode_kelas": "string",
        "nama": "string",
        "kode_prodi": "string",
        "kode_tahun_akademik": "string",
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
- **Deskripsi**: Membuat kelas baru.
- **Header**: Memerlukan token JWT.
- **Body Permintaan**:
  ```json
  {
    "kode": "string",
    "nama": "string",
    "prodi_id": "number",
    "tahun_akademik_id": "number"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Kelas Berhasil Dibuat",
    "data": {
      "id": "number",
      "kode_kelas": "string",
      "nama": "string",
      "kode_prodi": "string",
      "kode_tahun_akademik": "string",
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
- **Deskripsi**: Mendapatkan detail kelas berdasarkan ID.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Berhasil Mengambil Data Kelas",
    "data": {
      "kode": "string",
      "nama": "string",
      "kode_prodi": "string",
      "nama_prodi": "string",
      "tahun_akademik": {
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
- **Deskripsi**: Mendapatkan daftar semua tahun akademik.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "status": 200,
    "message": "Berhasil Mengambil Data Tahun Akademik",
    "data": [
      {
        "id": "number",
        "kode_tahun_akademik": "string",
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
- **Deskripsi**: Membuat tahun akademik baru.
- **Header**: Memerlukan token JWT.
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
- **Respon Sukses (200)**:
  ```json
  {
    "status": 200,
    "message": "Berhasil Membuat Tahun Akademik",
    "data": {
      "id": "number",
      "kode_tahun_akademik": "string",
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
- **Deskripsi**: Mendapatkan tahun akademik tertentu berdasarkan ID.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Berhasil Mengambil Data Tahun Akademik",
    "data": {
      "id": "number",
      "kode_tahun_akademik": "string",
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

#### Dapatkan Ringkasan Akademik
- **URL**: `/akademik-summary`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan ringkasan akademik (total kelas dan mahasiswa) untuk tahun akademik aktif.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "status": 200,
    "message": "Berhasil mengambil ringkasan akademik",
    "tahun_akademik": "string",
    "semester": "string",
    "total_kelas": "number",
    "total_mahasiswa": "number"
  }
  ```

### Mata Kuliah

#### Dapatkan Semua Mata Kuliah
- **URL**: `/matakuliah`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan daftar semua mata kuliah.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
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
        "kode_tahun_akademik": "string",
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
- **Deskripsi**: Membuat mata kuliah baru.
- **Header**: Memerlukan token JWT.
- **Body Permintaan**:
  ```json
  {
    "kode_matakuliah": "string",
    "nama_matakuliah": "string",
    "sks": "number",
    "semester": "number",
    "kode_prodi": "string",
    "kode_tahun_akademik": "string"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Matakuliah Berhasil Dibuat"
  }
  ```

#### Detail Mata Kuliah
- **URL**: `/matakuliah/:kode_matakuliah`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan detail mata kuliah berdasarkan kode.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "data": {
      "id": "number",
      "kode_matakuliah": "string",
      "nama_matakuliah": "string",
      "kode_prodi": "string",
      "sks": "number",
      "semester": "number",
      "kode_tahun_akademik": "string",
      "status": "string",
      "created_at": "datetime",
      "created_by": "string",
      "updated_at": "datetime"
    }
  }
  ```

## Jadwal Mata Kuliah

### 1. Membuat Jadwal Baru
- **URL**: `/jadwal/`
- **Method**: `POST`
- **Deskripsi**: Membuat jadwal mata kuliah baru.
- **Header**: Memerlukan token JWT (hanya admin prodi atau admin akademik).
- **Body Permintaan**:
  ```json
  {
    "kode_matakuliah": "string",
    "kelas_id": "number",
    "kode_dosen": "string",
    "kode_ruangan": "string",
    "kode_prodi": "string",
    "jam_mulai": "datetime",
    "jam_selesai": "datetime",
    "hari": "string",
    "kode_tahun_akademik": "string"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Jadwal berhasil dibuat",
    "data": {
      "id": "number",
      "kode_matakuliah": "string",
      "kelas_id": "number",
      "kode_dosen": "string",
      "kode_ruangan": "string",
      "kode_prodi": "string",
      "jam_mulai": "datetime",
      "jam_selesai": "datetime",
      "hari": "string",
      "kode_tahun_akademik": "string",
      "created_at": "datetime",
      "created_by": "string"
    }
  }
  ```

### 2. Melihat Jadwal per Kelas
- **URL**: `/jadwal/kelas/:kelas_id`
- **Method**: `GET`
- **Deskripsi**: Melihat jadwal mata kuliah untuk kelas tertentu.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "data": [
      {
        "id": "number",
        "kode_matakuliah": "string",
        "kelas_id": "number",
        "kode_dosen": "string",
        "kode_ruangan": "string",
        "kode_prodi": "string",
        "jam_mulai": "datetime",
        "jam_selesai": "datetime",
        "hari": "string",
        "kode_tahun_akademik": "string",
        "created_at": "datetime",
        "created_by": "string",
        "updated_at": "datetime"
      }
    ]
  }
  ```

### 3. Melihat Jadwal per Dosen
- **URL**: `/jadwal/dosen/:kode_dosen`
- **Method**: `GET`
- **Deskripsi**: Melihat jadwal mata kuliah untuk dosen tertentu.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "data": [
      {
        "id": "number",
        "kode_matakuliah": "string",
        "kelas_id": "number",
        "kode_dosen": "string",
        "kode_ruangan": "string",
        "kode_prodi": "string",
        "jam_mulai": "datetime",
        "jam_selesai": "datetime",
        "hari": "string",
        "kode_tahun_akademik": "string",
        "created_at": "datetime",
        "created_by": "string",
        "updated_at": "datetime"
      }
    ]
  }
  ```

### 4. Update Jadwal
- **URL**: `/jadwal/:id`
- **Method**: `PUT`
- **Deskripsi**: Memperbarui jadwal mata kuliah.
- **Header**: Memerlukan token JWT (hanya admin prodi atau admin akademik).
- **Body Permintaan**:
  ```json
  {
    "kode_ruangan": "string",
    "jam_mulai": "datetime",
    "jam_selesai": "datetime",
    "hari": "string"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Jadwal berhasil diupdate"
  }
  ```

### 5. Hapus Jadwal
- **URL**: `/jadwal/:id`
- **Method**: `DELETE`
- **Deskripsi**: Menghapus jadwal mata kuliah.
- **Header**: Memerlukan token JWT (hanya admin prodi atau admin akademik).
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Jadwal berhasil dihapus"
  }
  ```

### 6. Cek Ruangan Tersedia
- **URL**: `/jadwal/ruangan-tersedia`
- **Method**: `POST`
- **Deskripsi**: Memeriksa ketersediaan ruangan pada waktu tertentu.
- **Header**: Memerlukan token JWT.
- **Body Permintaan**:
  ```json
  {
    "hari": "string",
    "jam_mulai": "datetime",
    "jam_selesai": "datetime",
    "kode_prodi": "string"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "data": [
      {
        "kode_ruangan": "string",
        "nama_ruangan": "string",
        "kapasitas": "number",
        "lokasi": "string",
        "status": "string"
      }
    ]
  }
  ```

## Manajemen Ruangan

### 1. Membuat Ruangan Baru
- **URL**: `/ruangan/`
- **Method**: `POST`
- **Deskripsi**: Membuat ruangan baru.
- **Header**: Memerlukan token JWT (admin prodi atau admin akademik).
- **Body Permintaan**:
  ```json
  {
    "gedung": "string",
    "tipe": "string",
    "nomor": "number",
    "nama_ruangan": "string",
    "kode_prodi": "string",
    "kapasitas": "number",
    "lokasi": "string"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Ruangan berhasil dibuat",
    "data": {
      "kode_ruangan": "string",
      "nama_ruangan": "string",
      "kode_prodi": "string",
      "gedung": "string",
      "tipe": "string",
      "nomor": "number",
      "kapasitas": "number",
      "lokasi": "string",
      "status": "string",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  }
  ```

### 2. Melihat Ruangan per Gedung
- **URL**: `/ruangan/gedung/:gedung`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan daftar semua ruangan di gedung tertentu.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "data": [
      {
        "kode_ruangan": "string",
        "nama_ruangan": "string",
        "kode_prodi": "string",
        "gedung": "string",
        "tipe": "string",
        "nomor": "number",
        "kapasitas": "number",
        "lokasi": "string",
        "status": "string",
        "created_at": "datetime",
        "updated_at": "datetime"
      }
    ]
  }
  ```

### 3. Melihat Ruangan per Prodi
- **URL**: `/ruangan/prodi/:kode_prodi`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan daftar semua ruangan untuk program studi tertentu.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "data": [
      {
        "kode_ruangan": "string",
        "nama_ruangan": "string",
        "kode_prodi": "string",
        "gedung": "string",
        "tipe": "string",
        "nomor": "number",
        "kapasitas": "number",
        "lokasi": "string",
        "status": "string",
        "created_at": "datetime",
        "updated_at": "datetime"
      }
    ]
  }
  ```

### 4. Melihat Ruangan per Tipe
- **URL**: `/ruangan/tipe/:tipe`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan daftar ruangan berdasarkan tipe (LAB, KLS, STU).
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "data": [
      {
        "kode_ruangan": "string",
        "nama_ruangan": "string",
        "kode_prodi": "string",
        "gedung": "string",
        "tipe": "string",
        "nomor": "number",
        "kapasitas": "number",
        "lokasi": "string",
        "status": "string",
        "created_at": "datetime",
        "updated_at": "datetime"
      }
    ]
  }
  ```

### 5. Update Ruangan
- **URL**: `/ruangan/:kode_ruangan`
- **Method**: `PUT`
- **Deskripsi**: Memperbarui informasi ruangan.
- **Header**: Memerlukan token JWT (admin prodi atau admin akademik).
- **Body Permintaan**:
  ```json
  {
    "nama_ruangan": "string",
    "kapasitas": "number",
    "lokasi": "string",
    "status": "string"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Ruangan berhasil diupdate"
  }
  ```

### 6. Hapus Ruangan
- **URL**: `/ruangan/:kode_ruangan`
- **Method**: `DELETE`
- **Deskripsi**: Menghapus ruangan. Ruangan tidak dapat dihapus jika masih digunakan dalam jadwal.
- **Header**: Memerlukan token JWT (admin prodi atau admin akademik).
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Ruangan berhasil dihapus"
  }
  ```

### 7. Cek Status Ruangan
- **URL**: `/ruangan/status/:kode_ruangan?hari={string}`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan status ruangan dan jadwalnya untuk hari tertentu (opsional).
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "ruangan": {
      "kode_ruangan": "string",
      "nama_ruangan": "string",
      "kode_prodi": "string",
      "gedung": "string",
      "tipe": "string",
      "nomor": "number",
      "kapasitas": "number",
      "lokasi": "string",
      "status": "string"
    },
    "jadwal": [
      {
        "id": "number",
        "kode_matakuliah": "string",
        "kelas_id": "number",
        "kode_dosen": "string",
        "jam_mulai": "datetime",
        "jam_selesai": "datetime",
        "hari": "string"
      }
    ]
  }
  ```

## Absensi

### 1. Buka Pertemuan
- **URL**: `/buka-pertemuan`
- **Method**: `POST`
- **Deskripsi**: Membuka pertemuan baru untuk absensi.
- **Header**: Memerlukan token JWT (hanya dosen).
- **Body Permintaan**:
  ```json
  {
    "kode_matakuliah": "string",
    "pertemuan_ke": "number",
    "durasi": "number" // durasi dalam menit
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Pertemuan berhasil dibuka",
    "data": {
      "kode_pertemuan": "string",
      "waktu_dibuka": "datetime",
      "waktu_ditutup": "datetime"
    }
  }
  ```

### 2. Tutup Pertemuan
- **URL**: `/tutup-pertemuan/:kode_pertemuan`
- **Method**: `POST`
- **Deskripsi**: Menutup pertemuan yang sedang aktif.
- **Header**: Memerlukan token JWT (hanya dosen).
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Pertemuan berhasil ditutup"
  }
  ```

### 3. Absen
- **URL**: `/absen`
- **Method**: `POST`
- **Deskripsi**: Melakukan absensi untuk pertemuan yang sedang aktif.
- **Header**: Memerlukan token JWT (hanya mahasiswa).
- **Body Permintaan**:
  ```json
  {
    "kode_pertemuan": "string"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Absensi berhasil"
  }
  ```

### 4. Cek Status Pertemuan
- **URL**: `/status-pertemuan/:kode_matakuliah`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan status pertemuan terakhir untuk mata kuliah tertentu.
- **Header**: Tidak memerlukan token JWT (digunakan untuk publik).
- **Respon Sukses (200)**:
  ```json
  {
    "status": "string", // "dibuka" atau "belum_dibuka"
    "pertemuan_ke": "number",
    "kode_pertemuan": "string",
    "sisa_waktu": "number", // detik
    "waktu_dibuka": "datetime",
    "waktu_ditutup": "datetime"
  }
  ```

### 5. Rekap Absensi
- **URL**: `/rekap-absensi/:kode_matakuliah`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan rekap absensi untuk mata kuliah tertentu.
- **Header**: Memerlukan token JWT (hanya dosen).
- **Respon Sukses (200)**:
  ```json
  {
    "rekap_absensi": [
      {
        "nim": "string",
        "nama_mahasiswa": "string",
        "total_hadir": "number",
        "total_pertemuan": "number",
        "persentase_kehadiran": "number",
        "detail_absensi": [
          {
            "kode_pertemuan": "string",
            "pertemuan_ke": "number",
            "waktu_absen": "datetime",
            "status": "string" // "Hadir" atau "Tidak Hadir"
          }
        ]
      }
    ]
  }
  ```

## Penilaian

### 1. Membuat Penilaian Baru
- **URL**: `/penilaian/`
- **Method**: `POST`
- **Deskripsi**: Membuat entri nilai baru.
- **Header**: Memerlukan token JWT (hanya dosen atau admin prodi).
- **Body Permintaan**:
  ```json
  {
    "kode_matakuliah": "string",
    "kelas_id": "number",
    "nim": "string",
    "tugas": "number",
    "uts": "number",
    "uas": "number",
    "keterangan": "string"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Penilaian berhasil disimpan",
    "data": {
      "id": "number",
      "kode_matakuliah": "string",
      "kelas_id": "number",
      "nim": "string",
      "tugas": "number",
      "uts": "number",
      "uas": "number",
      "nilai_akhir": "number",
      "grade": "string",
      "keterangan": "string",
      "created_at": "datetime",
      "created_by": "string"
    }
  }
  ```

### 2. Melihat Penilaian per Kelas
- **URL**: `/penilaian/kelas/:kode_matakuliah/:kelas_id`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan semua nilai untuk kelas tertentu.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "data": [
      {
        "id": "number",
        "kode_matakuliah": "string",
        "kelas_id": "number",
        "nim": "string",
        "tugas": "number",
        "uts": "number",
        "uas": "number",
        "nilai_akhir": "number",
        "grade": "string",
        "keterangan": "string",
        "created_at": "datetime",
        "created_by": "string",
        "updated_at": "datetime"
      }
    ]
  }
  ```

### 3. Melihat Penilaian per Mahasiswa
- **URL**: `/penilaian/mahasiswa/:nim`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan semua nilai untuk mahasiswa tertentu.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "data": [
      {
        "id": "number",
        "kode_matakuliah": "string",
        "kelas_id": "number",
        "nim": "string",
        "tugas": "number",
        "uts": "number",
        "uas": "number",
        "nilai_akhir": "number",
        "grade": "string",
        "keterangan": "string",
        "created_at": "datetime",
        "created_by": "string",
        "updated_at": "datetime"
      }
    ]
  }
  ```

### 4. Update Penilaian
- **URL**: `/penilaian/:id`
- **Method**: `PUT`
- **Deskripsi**: Memperbarui entri nilai.
- **Header**: Memerlukan token JWT (hanya dosen atau admin prodi).
- **Body Permintaan**:
  ```json
  {
    "tugas": "number",
    "uts": "number",
    "uas": "number",
    "keterangan": "string"
  }
  ```
- **Respon Sukses (200)**:
  ```json
  {
    "message": "Penilaian berhasil diupdate"
  }
  ```

### 5. Rekap Nilai
- **URL**: `/penilaian/rekap/:kode_matakuliah/:kelas_id`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan rekap nilai dan distribusi nilai untuk mata kuliah dan kelas tertentu.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "data": {
      "kode_matakuliah": "string",
      "nama_matakuliah": "string",
      "kelas_id": "number",
      "nama_kelas": "string",
      "total_mahasiswa": "number",
      "rata_rata": "number",
      "distribusi": {
        "A": "number",
        "A-": "number",
        "B+": "number",
        "B": "number",
        "B-": "number",
        "C+": "number",
        "C": "number",
        "C-": "number",
        "D": "number",
        "E": "number"
      },
      "detail_nilai": [
        {
          "id": "number",
          "kode_matakuliah": "string",
          "kelas_id": "number",
          "nim": "string",
          "tugas": "number",
          "uts": "number",
          "uas": "number",
          "nilai_akhir": "number",
          "grade": "string",
          "keterangan": "string"
        }
      ]
    }
  }
  ```

## Mahasiswa

### 1. Dapatkan IPK Mahasiswa
- **URL**: `/mahasiswa/ipk/:nim`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan IPK (Indeks Prestasi Kumulatif) mahasiswa.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "message": "IPK berhasil diambil",
    "data": {
      "nim": "string",
      "nama": "string",
      "ipk": "number"
    }
  }
  ```

### 2. Dapatkan IPS Mahasiswa
- **URL**: `/mahasiswa/ips/:nim/:tahun_akademik`
- **Method**: `GET`
- **Deskripsi**: Mendapatkan IPS (Indeks Prestasi Semester) mahasiswa untuk tahun akademik tertentu.
- **Header**: Memerlukan token JWT.
- **Respon Sukses (200)**:
  ```json
  {
    "message": "IPS berhasil diambil",
    "data": {
      "nim": "string",
      "nama": "string",
      "tahun_akademik": "string",
      "ips": "number"
    }
  }
  ```