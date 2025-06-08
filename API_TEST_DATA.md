# API Test Data Documentation

## Authentication

### Login
```http
POST /login
Content-Type: application/json

{
    "username": "admin",
    "password": "admin123"
}
```

Expected Response:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "username": "admin",
        "role": "admin"
    }
}
```

## Program Studi

### Get All Prodi
```http
GET /prodi
Authorization: Bearer {token}
```

Expected Response:
```json
{
    "data": [
        {
            "kode": "TI",
            "nama": "Teknik Informatika"
        },
        {
            "kode": "SI",
            "nama": "Sistem Informasi"
        }
    ]
}
```

### Create Prodi
```http
POST /prodi
Authorization: Bearer {token}
Content-Type: application/json

{
    "kode": "MI",
    "nama": "Manajemen Informatika"
}
```

Expected Response:
```json
{
    "message": "Prodi berhasil ditambahkan",
    "data": {
        "kode": "MI",
        "nama": "Manajemen Informatika"
    }
}
```

## Ruangan

### Get All Ruangan
```http
GET /ruangan
Authorization: Bearer {token}
```

Expected Response:
```json
{
    "data": [
        {
            "kode": "LAB-1",
            "nama": "Laboratorium 1",
            "kapasitas": 30
        },
        {
            "kode": "LAB-2",
            "nama": "Laboratorium 2",
            "kapasitas": 30
        }
    ]
}
```

### Create Ruangan
```http
POST /ruangan
Authorization: Bearer {token}
Content-Type: application/json

{
    "kode": "LAB-3",
    "nama": "Laboratorium 3",
    "kapasitas": 40
}
```

Expected Response:
```json
{
    "message": "Ruangan berhasil ditambahkan",
    "data": {
        "kode": "LAB-3",
        "nama": "Laboratorium 3",
        "kapasitas": 40
    }
}
```

## Jadwal

### Get All Jadwal
```http
GET /jadwal
Authorization: Bearer {token}
```

Expected Response:
```json
{
    "data": [
        {
            "kode_mk": "IF101",
            "nama_mk": "Pemrograman Dasar",
            "kelas": "A",
            "sks": 3,
            "dosen": "Dr. John Doe",
            "ruangan": "LAB-1",
            "hari": "Senin",
            "jam_mulai": "08:00",
            "jam_selesai": "10:30"
        },
        {
            "kode_mk": "IF102",
            "nama_mk": "Basis Data",
            "kelas": "B",
            "sks": 3,
            "dosen": "Dr. Jane Doe",
            "ruangan": "LAB-2",
            "hari": "Selasa",
            "jam_mulai": "13:00",
            "jam_selesai": "15:30"
        }
    ]
}
```

### Create Jadwal
```http
POST /jadwal
Authorization: Bearer {token}
Content-Type: application/json

{
    "kode_mk": "IF103",
    "nama_mk": "Pemrograman Web",
    "kelas": "C",
    "sks": 3,
    "dosen": "Dr. Alice Smith",
    "ruangan": "LAB-1",
    "hari": "Rabu",
    "jam_mulai": "10:00",
    "jam_selesai": "12:30"
}
```

Expected Response:
```json
{
    "message": "Jadwal berhasil ditambahkan",
    "data": {
        "kode_mk": "IF103",
        "nama_mk": "Pemrograman Web",
        "kelas": "C",
        "sks": 3,
        "dosen": "Dr. Alice Smith",
        "ruangan": "LAB-1",
        "hari": "Rabu",
        "jam_mulai": "10:00",
        "jam_selesai": "12:30"
    }
}
```

## Absensi

### Get All Absensi
```http
GET /absensi
Authorization: Bearer {token}
```

Expected Response:
```json
{
    "data": [
        {
            "kode_mk": "IF101",
            "kelas": "A",
            "tanggal": "2024-03-20",
            "status": "hadir"
        },
        {
            "kode_mk": "IF102",
            "kelas": "B",
            "tanggal": "2024-03-20",
            "status": "izin"
        }
    ]
}
```

### Create Absensi
```http
POST /absensi
Authorization: Bearer {token}
Content-Type: application/json

{
    "kode_mk": "IF101",
    "kelas": "A",
    "tanggal": "2024-03-21",
    "status": "hadir"
}
```

Expected Response:
```json
{
    "message": "Absensi berhasil ditambahkan",
    "data": {
        "kode_mk": "IF101",
        "kelas": "A",
        "tanggal": "2024-03-21",
        "status": "hadir"
    }
}
```

## Penilaian

### Get All Penilaian
```http
GET /penilaian
Authorization: Bearer {token}
```

Expected Response:
```json
{
    "data": [
        {
            "kode_mk": "IF101",
            "kelas": "A",
            "nim": "2021001",
            "nilai_tugas": 85,
            "nilai_uts": 80,
            "nilai_uas": 90
        },
        {
            "kode_mk": "IF102",
            "kelas": "B",
            "nim": "2021002",
            "nilai_tugas": 90,
            "nilai_uts": 85,
            "nilai_uas": 95
        }
    ]
}
```

### Create Penilaian
```http
POST /penilaian
Authorization: Bearer {token}
Content-Type: application/json

{
    "kode_mk": "IF101",
    "kelas": "A",
    "nim": "2021003",
    "nilai_tugas": 88,
    "nilai_uts": 85,
    "nilai_uas": 92
}
```

Expected Response:
```json
{
    "message": "Penilaian berhasil ditambahkan",
    "data": {
        "kode_mk": "IF101",
        "kelas": "A",
        "nim": "2021003",
        "nilai_tugas": 88,
        "nilai_uts": 85,
        "nilai_uas": 92
    }
}
```

## Error Responses

### Unauthorized
```json
{
    "error": "Unauthorized",
    "message": "Token tidak valid atau expired"
}
```

### Validation Error
```json
{
    "error": "Validation Error",
    "message": "Data tidak valid",
    "errors": {
        "kode_mk": ["Kode mata kuliah harus diisi"],
        "kelas": ["Kelas harus diisi"]
    }
}
```

### Not Found
```json
{
    "error": "Not Found",
    "message": "Data tidak ditemukan"
}
```

## Testing Tools

Anda dapat menggunakan tools berikut untuk testing API:

1. **Postman**
   - Import collection dari file `Simpadu.postman_collection.json`
   - Set environment variables untuk `base_url` dan `token`

2. **cURL**
   ```bash
   # Login
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"admin123"}'

   # Get All Prodi
   curl http://localhost:8080/prodi \
     -H "Authorization: Bearer {token}"
   ```

3. **JavaScript Fetch**
   ```javascript
   // Login
   fetch('http://localhost:8080/login', {
     method: 'POST',
     headers: {
       'Content-Type': 'application/json'
     },
     body: JSON.stringify({
       username: 'admin',
       password: 'admin123'
     })
   })
   .then(response => response.json())
   .then(data => console.log(data));

   // Get All Prodi
   fetch('http://localhost:8080/prodi', {
     headers: {
       'Authorization': 'Bearer ' + token
     }
   })
   .then(response => response.json())
   .then(data => console.log(data));
   ```