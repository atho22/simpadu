-- Insert test data for tahun akademik
INSERT INTO tahun_akademiks (tahun, semester, status, tanggal_mulai, tanggal_selesai, created_at, updated_at)
VALUES 
('2023/2024', 'Ganjil', 'Aktif', '2023-09-01', '2024-01-31', NOW(), NOW()),
('2023/2024', 'Genap', 'Tidak Aktif', '2024-02-01', '2024-06-30', NOW(), NOW());

-- Insert test data for jurusan
INSERT INTO jurusans (kode_jurusan, nama_jurusan, created_at, updated_at)
VALUES 
('J001', 'Teknik Informatika', NOW(), NOW()),
('J002', 'Sistem Informasi', NOW(), NOW());

-- Insert test data for prodi
INSERT INTO prodis (kode_prodi, nama_prodi, kode_jurusan, jenjang, created_at, updated_at)
VALUES 
('TI001', 'Teknik Informatika', 'J001', 'S1', NOW(), NOW()),
('SI001', 'Sistem Informasi', 'J002', 'S1', NOW(), NOW());

-- Insert test data for users (admin and students)
INSERT INTO users (user_id, username, email, password, role, status)
VALUES 
('ADM001', 'admin_ti', 'admin.ti@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'admin_prodi', 'Aktif'),
('MHS001', 'student1', 'student1@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'mahasiswa', 'Aktif'),
('MHS002', 'student2', 'student2@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'mahasiswa', 'Aktif'),
('MHS003', 'student3', 'student3@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'mahasiswa', 'Aktif'),
('MHS004', 'student4', 'student4@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'mahasiswa', 'Aktif');

-- Insert test data for mahasiswa
INSERT INTO mahasiswas (user_id, nim, nama, angkatan, kode_prodi, created_at, updated_at)
VALUES 
('MHS001', '2021001', 'John Doe', 2021, 'TI001', NOW(), NOW()),
('MHS002', '2021002', 'Jane Smith', 2021, 'TI001', NOW(), NOW()),
('MHS003', '2021003', 'Bob Johnson', 2021, 'TI001', NOW(), NOW()),
('MHS004', '2021004', 'Alice Brown', 2021, 'TI001', NOW(), NOW());

-- Insert test data for struktural prodi
INSERT INTO struktural_prodis (user_id, kode_prodi, nama, nip, jabatan, tanggal_sk, created_at, updated_at)
VALUES 
('ADM001', 'TI001', 'Dr. Admin TI', '198501012010011001', 'Ketua Prodi', NOW(), NOW(), NOW());

-- Insert test data for kelas
INSERT INTO kelas (kode, nama, kode_prodi, tahun_akademik_id, created_by, created_at, updated_at)
VALUES 
('KLS001', 'Kelas A', 'TI001', 1, 'ADM001', NOW(), NOW()),
('KLS002', 'Kelas B', 'TI001', 1, 'ADM001', NOW(), NOW());

-- Insert test data for kelas_mahasiswa (assign some students to kelas)
INSERT INTO kelas_mahasiswas (kelas_id, nim, created_by, created_at, updated_at)
VALUES 
(1, '2021001', 'ADM001', NOW(), NOW()),
(1, '2021002', 'ADM001', NOW(), NOW()),
(2, '2021003', 'ADM001', NOW(), NOW());

-- Note: Student with NIM 2021004 is left unassigned to test available_mahasiswa functionality 