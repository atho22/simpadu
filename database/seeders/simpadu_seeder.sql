-- Seed Users
INSERT INTO users (user_id, username, email, password, role) VALUES
(1, 'admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin'),
(2, 'dosen', 'dosen@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'dosen'),
(3, 'mahasiswa', 'mahasiswa@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'mahasiswa');

-- Seed Program Studi
INSERT INTO prodi (kode, nama) VALUES
('TI', 'Teknik Informatika'),
('SI', 'Sistem Informasi'),
('MI', 'Manajemen Informatika');

-- Seed Ruangan
INSERT INTO ruangan (kode, nama, kapasitas) VALUES
('LAB-1', 'Laboratorium 1', 30),
('LAB-2', 'Laboratorium 2', 30),
('LAB-3', 'Laboratorium 3', 40),
('LAB-4', 'Laboratorium 4', 40),
('LAB-5', 'Laboratorium 5', 35);

-- Seed Jadwal
INSERT INTO jadwal (kode_mk, nama_mk, kelas, sks, dosen, ruangan, hari, jam_mulai, jam_selesai) VALUES
('IF101', 'Pemrograman Dasar', 'A', 3, 'Dr. John Doe', 'LAB-1', 'Senin', '08:00', '10:30'),
('IF102', 'Basis Data', 'B', 3, 'Dr. Jane Doe', 'LAB-2', 'Selasa', '13:00', '15:30'),
('IF103', 'Pemrograman Web', 'C', 3, 'Dr. Alice Smith', 'LAB-1', 'Rabu', '10:00', '12:30'),
('IF104', 'Jaringan Komputer', 'A', 3, 'Dr. Bob Johnson', 'LAB-3', 'Kamis', '08:00', '10:30'),
('IF105', 'Kecerdasan Buatan', 'B', 3, 'Dr. Carol White', 'LAB-4', 'Jumat', '13:00', '15:30');

-- Seed Absensi
INSERT INTO absensi (kode_mk, kelas, tanggal, status) VALUES
('IF101', 'A', CURRENT_DATE, 'hadir'),
('IF102', 'B', CURRENT_DATE, 'izin'),
('IF103', 'C', CURRENT_DATE, 'hadir'),
('IF104', 'A', CURRENT_DATE, 'sakit'),
('IF105', 'B', CURRENT_DATE, 'hadir');

-- Seed Penilaian
INSERT INTO penilaian (kode_mk, kelas, nim, nilai_tugas, nilai_uts, nilai_uas) VALUES
('IF101', 'A', '2021001', 85, 80, 90),
('IF102', 'B', '2021002', 90, 85, 95),
('IF103', 'C', '2021003', 88, 85, 92),
('IF104', 'A', '2021004', 92, 88, 94),
('IF105', 'B', '2021005', 87, 90, 93);

-- Seed Mahasiswa
INSERT INTO mahasiswa (nim, nama, prodi, angkatan) VALUES
('2021001', 'Ahmad Rizki', 'TI', 2021),
('2021002', 'Budi Santoso', 'SI', 2021),
('2021003', 'Citra Dewi', 'TI', 2021),
('2021004', 'Dian Pratama', 'MI', 2021),
('2021005', 'Eka Putri', 'SI', 2021);

-- Seed Dosen
INSERT INTO dosen (nip, nama, prodi) VALUES
('19800101199001001', 'Dr. John Doe', 'TI'),
('19800202199002002', 'Dr. Jane Doe', 'SI'),
('19800303199003003', 'Dr. Alice Smith', 'TI'),
('19800404199004004', 'Dr. Bob Johnson', 'MI'),
('19800505199005005', 'Dr. Carol White', 'SI');

-- Seed Mata Kuliah
INSERT INTO mata_kuliah (kode, nama, sks, prodi) VALUES
('IF101', 'Pemrograman Dasar', 3, 'TI'),
('IF102', 'Basis Data', 3, 'TI'),
('IF103', 'Pemrograman Web', 3, 'TI'),
('SI101', 'Sistem Informasi', 3, 'SI'),
('MI101', 'Manajemen Proyek', 3, 'MI');

-- Seed Kelas
INSERT INTO kelas (kode_mk, kelas, dosen, ruangan, hari, jam_mulai, jam_selesai) VALUES
('IF101', 'A', '19800101199001001', 'LAB-1', 'Senin', '08:00', '10:30'),
('IF102', 'B', '19800202199002002', 'LAB-2', 'Selasa', '13:00', '15:30'),
('IF103', 'C', '19800303199003003', 'LAB-1', 'Rabu', '10:00', '12:30'),
('SI101', 'A', '19800404199004004', 'LAB-3', 'Kamis', '08:00', '10:30'),
('MI101', 'B', '19800505199005005', 'LAB-4', 'Jumat', '13:00', '15:30');

-- Seed KRS
INSERT INTO krs (nim, kode_mk, kelas, semester, tahun_ajaran) VALUES
('2021001', 'IF101', 'A', 'Ganjil', '2023/2024'),
('2021002', 'IF102', 'B', 'Ganjil', '2023/2024'),
('2021003', 'IF103', 'C', 'Ganjil', '2023/2024'),
('2021004', 'SI101', 'A', 'Ganjil', '2023/2024'),
('2021005', 'MI101', 'B', 'Ganjil', '2023/2024'); 