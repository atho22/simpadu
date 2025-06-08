-- Disable foreign key checks temporarily
SET FOREIGN_KEY_CHECKS = 0;

-- Truncate all tables
TRUNCATE TABLE penilaian;
TRUNCATE TABLE absensi;
TRUNCATE TABLE pertemuan;
TRUNCATE TABLE jadwal_matakuliah;
TRUNCATE TABLE mata_kuliah;
TRUNCATE TABLE kelas_mahasiswa;
TRUNCATE TABLE kelas;
TRUNCATE TABLE mahasiswa;
TRUNCATE TABLE dosen;
TRUNCATE TABLE struktural_prodi;
TRUNCATE TABLE prodi;
TRUNCATE TABLE jurusan;
TRUNCATE TABLE tahun_akademik;
TRUNCATE TABLE kode_ruangan;
TRUNCATE TABLE users;

-- Enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1; 