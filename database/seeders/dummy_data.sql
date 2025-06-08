-- Seed Users
INSERT INTO users (user_id, username, email, password, role, status) VALUES
('U001', 'admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', 'Aktif'),
('U002', 'dosen1', 'dosen1@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'dosen', 'Aktif'),
('U003', 'dosen2', 'dosen2@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'dosen', 'Aktif'),
('U004', 'mahasiswa1', 'mahasiswa1@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'mahasiswa', 'Aktif'),
('U005', 'mahasiswa2', 'mahasiswa2@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'mahasiswa', 'Aktif'),
('U006', 'ketua_ti', 'ketua_ti@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin_prodi', 'Aktif'),
('U007', 'ketua_si', 'ketua_si@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin_prodi', 'Aktif'),
('U008', 'ketua_mi', 'ketua_mi@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin_prodi', 'Aktif');

-- Seed Tahun_akademik
INSERT INTO tahun_akademiks (kode_tahun_akademik, tahun, semester, status, tanggal_mulai, tanggal_selesai) VALUES
('2023-1', '2023/2024', 'Ganjil', 'Aktif', '2023-09-01', '2024-01-31'),
('2023-2', '2023/2024', 'Genap', 'Tidak Aktif', '2024-02-01', '2024-06-30');

-- Seed Jurusan
INSERT INTO jurusans (kode_jurusan, nama_jurusan) VALUES
('J001', 'Teknik Informatika'),
('J002', 'Sistem Informasi'),
('J003', 'Manajemen Informatika');

-- Seed Prodi
INSERT INTO prodis (kode_prodi, nama_prodi, kode_jurusan, jenjang, user_id) VALUES
('TI', 'Teknik Informatika', 'J001', 'S1', 'U006'),
('SI', 'Sistem Informasi', 'J002', 'S1', 'U007'),
('MI', 'Manajemen Informatika', 'J003', 'D3', 'U008');

-- Seed Struktural_prodi
INSERT INTO struktural_prodis (user_id, kode_prodi, nama, nip, jabatan, no_sk, tanggal_sk) VALUES
('U006', 'TI', 'Dr. John Doe', '19800101199001001', 'Ketua Prodi', 'SK001', '2023-01-01'),
('U007', 'SI', 'Dr. Jane Doe', '19800202199002002', 'Ketua Prodi', 'SK002', '2023-01-01'),
('U008', 'MI', 'Dr. Alice Smith', '19800303199003003', 'Ketua Prodi', 'SK003', '2023-01-01');

-- Seed Dosen
INSERT INTO dosens (user_id, kode_dosen, nip, nidn, kode_prodi) VALUES
('U002', 'D001', '19800101199001001', '0001018001', 'TI'),
('U003', 'D002', '19800202199002002', '0002018002', 'SI');

-- Seed Kelas
INSERT INTO kelass (kode_kelas, nama, kode_prodi, kode_tahun_akademik, created_by) VALUES
('TI-1A', 'Kelas 1A TI', 'TI', '2023-1', 'U006'),
('TI-1B', 'Kelas 1B TI', 'TI', '2023-1', 'U006'),
('SI-1A', 'Kelas 1A SI', 'SI', '2023-1', 'U007');

-- Seed Mahasiswa
INSERT INTO mahasiswas (user_id, nim, nama, angkatan, kode_prodi) VALUES
('U004', '2021001', 'Ahmad Rizki', 2021, 'TI'),
('U005', '2021002', 'Budi Santoso', 2021, 'SI');

-- Seed Kelas_mahasiswa
INSERT INTO kelas_mahasiswas (kelas_id, nim, created_by) VALUES
(1, '2021001', 'U006'),
(2, '2021002', 'U007');

-- Seed Mata_kuliah
INSERT INTO mata_kuliahs (kode_matakuliah, nama_matakuliah, kode_prodi, sks, semester, kode_tahun_akademik, created_by) VALUES
('IF101', 'Pemrograman Dasar', 'TI', 3, 1, '2023-1', 'U006'),
('IF102', 'Basis Data', 'TI', 3, 1, '2023-1', 'U006'),
('SI101', 'Sistem Informasi', 'SI', 3, 1, '2023-1', 'U007');

-- Seed Jadwal_matakuliah
INSERT INTO jadwal_matakuliahs (kode_matakuliah, kelas_id, kode_dosen, kode_ruangan, kode_prodi, jam_mulai, jam_selesai, hari, kode_tahun_akademik, created_by) VALUES
('IF101', 1, 'D001', 'LAB-1', 'TI', '2023-09-01 08:00:00', '2023-09-01 10:30:00', 'Senin', '2023-1', 'U006'),
('IF102', 1, 'D001', 'LAB-2', 'TI', '2023-09-01 13:00:00', '2023-09-01 15:30:00', 'Selasa', '2023-1', 'U006'),
('SI101', 3, 'D002', 'LAB-3', 'SI', '2023-09-01 08:00:00', '2023-09-01 10:30:00', 'Rabu', '2023-1', 'U007');

-- Seed Kode_ruangan
INSERT INTO kode_ruangans (kode_ruangan, nama_ruangan, kode_prodi, gedung, tipe, nomor, kapasitas, lokasi) VALUES
('LAB-1', 'Laboratorium 1', 'TI', 'A', 'LAB', 1, 30, 'Lantai 1'),
('LAB-2', 'Laboratorium 2', 'TI', 'A', 'LAB', 2, 30, 'Lantai 1'),
('LAB-3', 'Laboratorium 3', 'SI', 'B', 'LAB', 1, 40, 'Lantai 2');

-- Seed Pertemuan
INSERT INTO pertemuans (kode_pertemuan, kode_matakuliah, pertemuan_ke, dibuka, waktu_dibuka, waktu_ditutup, dibuka_by) VALUES
('P001', 'IF101', 1, true, '2023-09-01 08:00:00', '2023-09-01 10:30:00', 'U002'),
('P002', 'IF102', 1, true, '2023-09-01 13:00:00', '2023-09-01 15:30:00', 'U002'),
('P003', 'SI101', 1, true, '2023-09-01 08:00:00', '2023-09-01 10:30:00', 'U003');

-- Seed Absensi
INSERT INTO absensis (kode_matakuliah, kode_pertemuan, nim, status, keterangan) VALUES
('IF101', 'P001', '2021001', 'hadir', 'Hadir tepat waktu'),
('IF102', 'P002', '2021001', 'izin', 'Izin sakit'),
('SI101', 'P003', '2021002', 'hadir', 'Hadir tepat waktu');

-- Seed Penilaian
INSERT INTO penilaians (kode_matakuliah, kelas_id, nim, tugas, uts, uas, nilai_akhir, grade, created_by, created_at) VALUES
('IF101', 1, '2021001', 85, 80, 90, 85, 'A', 'U002', NOW()),
('IF102', 1, '2021001', 90, 85, 95, 90, 'A', 'U002', NOW()),
('SI101', 3, '2021002', 88, 85, 92, 88, 'A', 'U003', NOW()); 