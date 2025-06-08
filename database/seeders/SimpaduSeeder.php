<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Hash;

class SimpaduSeeder extends Seeder
{
    public function run()
    {
        // Seed Users
        DB::table('users')->insert([
            [
                'username' => 'admin',
                'password' => Hash::make('admin123'),
                'role' => 'admin',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'username' => 'dosen',
                'password' => Hash::make('dosen123'),
                'role' => 'dosen',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'username' => 'mahasiswa',
                'password' => Hash::make('mahasiswa123'),
                'role' => 'mahasiswa',
                'created_at' => now(),
                'updated_at' => now()
            ]
        ]);

        // Seed Prodi
        DB::table('prodi')->insert([
            [
                'kode' => 'TI',
                'nama' => 'Teknik Informatika',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'kode' => 'SI',
                'nama' => 'Sistem Informasi',
                'created_at' => now(),
                'updated_at' => now()
            ]
        ]);

        // Seed Ruangan
        DB::table('ruangan')->insert([
            [
                'kode' => 'LAB-1',
                'nama' => 'Laboratorium 1',
                'kapasitas' => 30,
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'kode' => 'LAB-2',
                'nama' => 'Laboratorium 2',
                'kapasitas' => 30,
                'created_at' => now(),
                'updated_at' => now()
            ]
        ]);

        // Seed Jadwal
        DB::table('jadwal')->insert([
            [
                'kode_mk' => 'IF101',
                'nama_mk' => 'Pemrograman Dasar',
                'kelas' => 'A',
                'sks' => 3,
                'dosen' => 'Dr. John Doe',
                'ruangan' => 'LAB-1',
                'hari' => 'Senin',
                'jam_mulai' => '08:00',
                'jam_selesai' => '10:30',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'kode_mk' => 'IF102',
                'nama_mk' => 'Basis Data',
                'kelas' => 'B',
                'sks' => 3,
                'dosen' => 'Dr. Jane Doe',
                'ruangan' => 'LAB-2',
                'hari' => 'Selasa',
                'jam_mulai' => '13:00',
                'jam_selesai' => '15:30',
                'created_at' => now(),
                'updated_at' => now()
            ]
        ]);

        // Seed Absensi
        DB::table('absensi')->insert([
            [
                'kode_mk' => 'IF101',
                'kelas' => 'A',
                'tanggal' => now(),
                'status' => 'hadir',
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'kode_mk' => 'IF102',
                'kelas' => 'B',
                'tanggal' => now(),
                'status' => 'izin',
                'created_at' => now(),
                'updated_at' => now()
            ]
        ]);

        // Seed Penilaian
        DB::table('penilaian')->insert([
            [
                'kode_mk' => 'IF101',
                'kelas' => 'A',
                'nim' => '2021001',
                'nilai_tugas' => 85,
                'nilai_uts' => 80,
                'nilai_uas' => 90,
                'created_at' => now(),
                'updated_at' => now()
            ],
            [
                'kode_mk' => 'IF102',
                'kelas' => 'B',
                'nim' => '2021002',
                'nilai_tugas' => 90,
                'nilai_uts' => 85,
                'nilai_uas' => 95,
                'created_at' => now(),
                'updated_at' => now()
            ]
        ]);
    }
} 