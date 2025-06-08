-- Create Users table
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    status ENUM('Aktif', 'Tidak Aktif') DEFAULT 'Aktif'
);

-- Create Tahun_akademik table
CREATE TABLE tahun_akademik (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode_tahun_akademik VARCHAR(20) NOT NULL,
    tahun VARCHAR(9) NOT NULL,
    semester ENUM('Ganjil', 'Genap') NOT NULL,
    status ENUM('Aktif', 'Tidak Aktif') DEFAULT 'Tidak Aktif',
    tanggal_mulai DATETIME NOT NULL,
    tanggal_selesai DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Jurusan table
CREATE TABLE jurusan (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode_jurusan VARCHAR(20) NOT NULL,
    nama_jurusan VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Prodi table
CREATE TABLE prodi (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode_prodi VARCHAR(20) NOT NULL,
    nama_prodi VARCHAR(50) NOT NULL,
    kode_jurusan VARCHAR(20) NOT NULL,
    jenjang VARCHAR(20) NOT NULL,
    user_id VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Struktural_prodi table
CREATE TABLE struktural_prodi (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL UNIQUE,
    kode_prodi VARCHAR(20) NOT NULL,
    nama VARCHAR(100) NOT NULL,
    nip VARCHAR(30) UNIQUE,
    jabatan VARCHAR(50),
    no_sk VARCHAR(50),
    tanggal_sk DATETIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Dosen table
CREATE TABLE dosen (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    kode_dosen VARCHAR(20) NOT NULL,
    nip VARCHAR(20),
    nidn VARCHAR(20),
    kode_prodi VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Kelas table
CREATE TABLE kelas (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode_kelas VARCHAR(10) NOT NULL,
    nama VARCHAR(50) NOT NULL,
    kode_prodi VARCHAR(20) NOT NULL,
    kode_tahun_akademik VARCHAR(20) NOT NULL,
    created_by VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Mahasiswa table
CREATE TABLE mahasiswa (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    nim VARCHAR(20) NOT NULL,
    nama VARCHAR(100) NOT NULL,
    angkatan INT UNSIGNED NOT NULL,
    kode_prodi VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Kelas_mahasiswa table
CREATE TABLE kelas_mahasiswa (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kelas_id BIGINT UNSIGNED NOT NULL,
    nim VARCHAR(20) NOT NULL,
    created_by VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Mata_kuliah table
CREATE TABLE mata_kuliah (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode_matakuliah VARCHAR(20) NOT NULL,
    nama_matakuliah VARCHAR(100) NOT NULL,
    kode_prodi VARCHAR(20) NOT NULL,
    sks INT NOT NULL,
    semester INT NOT NULL,
    kode_tahun_akademik VARCHAR(20) NOT NULL,
    status ENUM('Aktif', 'Tidak Aktif') DEFAULT 'Aktif',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(20) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Jadwal_matakuliah table
CREATE TABLE jadwal_matakuliah (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode_matakuliah VARCHAR(20) NOT NULL,
    kelas_id BIGINT UNSIGNED NOT NULL,
    kode_dosen VARCHAR(20) NOT NULL,
    kode_ruangan VARCHAR(20) NOT NULL,
    kode_prodi VARCHAR(20) NOT NULL,
    jam_mulai DATETIME NOT NULL,
    jam_selesai DATETIME NOT NULL,
    hari VARCHAR(20) NOT NULL,
    kode_tahun_akademik VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(20) NOT NULL
);

-- Create Kode_ruangan table
CREATE TABLE kode_ruangan (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode_ruangan VARCHAR(20) NOT NULL UNIQUE,
    nama_ruangan VARCHAR(100) NOT NULL,
    kode_prodi VARCHAR(20) NOT NULL,
    gedung VARCHAR(20) NOT NULL,
    tipe VARCHAR(20) NOT NULL,
    nomor INT NOT NULL,
    kapasitas INT NOT NULL,
    lokasi VARCHAR(100) NOT NULL,
    status ENUM('Aktif', 'Maintenance', 'Tidak Aktif') DEFAULT 'Aktif',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create Pertemuan table
CREATE TABLE pertemuan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode_pertemuan VARCHAR(20) NOT NULL,
    kode_matakuliah VARCHAR(20) NOT NULL,
    pertemuan_ke INT NOT NULL,
    dibuka BOOLEAN NOT NULL,
    waktu_dibuka DATETIME NOT NULL,
    waktu_ditutup DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    dibuka_by VARCHAR(20) NOT NULL
);

-- Create Absensi table
CREATE TABLE absensi (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode_matakuliah VARCHAR(20) NOT NULL,
    kode_pertemuan VARCHAR(20) NOT NULL,
    nim VARCHAR(20) NOT NULL,
    status ENUM('hadir', 'izin', 'sakit', 'alfa') DEFAULT 'alfa',
    keterangan VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Penilaian table
CREATE TABLE penilaian (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode_matakuliah VARCHAR(20) NOT NULL,
    kelas_id BIGINT UNSIGNED NOT NULL,
    nim VARCHAR(20) NOT NULL,
    tugas DECIMAL(5,2) NOT NULL,
    uts DECIMAL(5,2) NOT NULL,
    uas DECIMAL(5,2) NOT NULL,
    nilai_akhir DECIMAL(5,2) NOT NULL,
    grade VARCHAR(2) NOT NULL,
    keterangan VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(20) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(20)
); 