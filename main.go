package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"simpadu/absensi"
	"simpadu/akademik"
	"simpadu/config"
	"simpadu/helper"
	"simpadu/jadwal"
	"simpadu/mahasiswa"
	"simpadu/penilaian"
	"simpadu/prodi"
	"simpadu/ruangan"
	"simpadu/structs"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Use environment variable for JWT key with a fallback for development
var jwtKey = getJWTKey()

// Structs

type TokenClaims struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Kode_prodi string `json:"kode_prodi"`
	Nama       string `json:"nama"`
	NIP        string `json:"nip"`
	NIM        string `json:"nim"`
	jwt.RegisteredClaims
}

// Global DB
var DB *gorm.DB

func main() {
	var err error

	// Validate JWT key is adequate
	if len(jwtKey) < 32 {
		log.Println("Warning: JWT key is less than 32 bytes. Consider using a stronger key.")
	}

	// Connect to database
	DB, err = config.ConnectDB()
	if err != nil {
		log.Fatal("Koneksi ke database gagal:", err)
	}

	// Migrate struct
	log.Println("Starting database migration...")
	if err := DB.AutoMigrate(&structs.Users{}, &structs.Tahun_akademik{}, &structs.Kelas{}, &structs.Mahasiswa{}, &structs.Kelas_mahasiswa{}, &structs.Jurusan{}, &structs.Prodi{}, &structs.Struktural_prodi{}, &structs.Mata_kuliah{}, &structs.Jadwal_matakuliah{}, &structs.Absensi{}, &structs.Penilaian{}, &structs.Kode_ruangan{}, &structs.Pertemuan{}, &structs.Dosen{}, &structs.Struktural_Jurusan{}, &structs.Staff_akademik{}); err != nil {
		log.Fatal("Database migration failed:", err)
	}
	log.Println("Database migration completed successfully")

	// Replace individual seed calls with SeedAll

	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Initialize handlers
	kelasHandler := &prodi.InDB{DB: DB}
	tahunAkademikHandler := &akademik.InDB{DB: DB}
	absensiHandler := &absensi.InDB{DB: DB}

	// Routes
	r.POST("/register-mahasiswa", RegisterUser)
	r.POST("/register-admin-prodi", RegisterAdminProdi)
	r.POST("/login", Login)

	// Kelas routes
	r.GET("/kelas", AuthMiddleware(), kelasHandler.GetAllKelas)
	r.GET("/kelas/:id", AuthMiddleware(), kelasHandler.DetailKelas)
	r.POST("/create-kelas", AuthMiddleware(), kelasHandler.CreateKelas)

	// Tahun Akademik routes
	r.GET("/tahun-akademik", AuthMiddleware(), tahunAkademikHandler.GetAllTahunAkademik)
	r.POST("/create-tahun-akademik", AuthMiddleware(), tahunAkademikHandler.CreateTahunAkademik)
	r.GET("/tahun-akademik/:id", AuthMiddleware(), tahunAkademikHandler.GetTahunAkademikByID)
	r.GET("/akademik-summary", AuthMiddleware(), tahunAkademikHandler.GetAcademicSummary)

	// Absensi routes
	r.POST("/buka-pertemuan", AuthMiddleware(), absensiHandler.BukaPertemuan)
	r.POST("/tutup-pertemuan/:kode_pertemuan", AuthMiddleware(), absensiHandler.TutupPertemuan)
	r.POST("/absen", AuthMiddleware(), absensiHandler.Absen)
	r.GET("/status-pertemuan/:kode_matakuliah", AuthMiddleware(), absensiHandler.GetStatusPertemuan)
	r.GET("/rekap-absensi/:kode_matakuliah", AuthMiddleware(), absensiHandler.GetRekapAbsensi)

	// Jadwal Mata Kuliah
	jadwalHandler := jadwal.InDB{DB: DB}
	jadwalRoutes := r.Group("/jadwal")
	{
		jadwalRoutes.POST("/", AuthMiddleware(), jadwalHandler.CreateJadwal)
		jadwalRoutes.GET("/kelas/:kelas_id", AuthMiddleware(), jadwalHandler.GetJadwalByKelas)
		jadwalRoutes.GET("/dosen/:kode_dosen", AuthMiddleware(), jadwalHandler.GetJadwalByDosen)
		jadwalRoutes.PUT("/:id", AuthMiddleware(), jadwalHandler.UpdateJadwal)
		jadwalRoutes.DELETE("/:id", AuthMiddleware(), jadwalHandler.DeleteJadwal)
		jadwalRoutes.POST("/ruangan-tersedia", AuthMiddleware(), jadwalHandler.GetRuanganTersedia)
	}

	// Penilaian routes
	penilaianHandler := penilaian.InDB{DB: DB}
	penilaianRoutes := r.Group("/penilaian")
	{
		penilaianRoutes.POST("/", AuthMiddleware(), penilaianHandler.CreatePenilaian)
		penilaianRoutes.GET("/kelas/:kode_matakuliah/:kelas_id", AuthMiddleware(), penilaianHandler.GetPenilaianByKelas)
		penilaianRoutes.GET("/mahasiswa/:nim", AuthMiddleware(), penilaianHandler.GetPenilaianByMahasiswa)
		penilaianRoutes.PUT("/:id", AuthMiddleware(), penilaianHandler.UpdatePenilaian)
		penilaianRoutes.GET("/rekap/:kode_matakuliah/:kelas_id", AuthMiddleware(), penilaianHandler.GetRekapNilai)
	}

	// Ruangan routes
	ruanganHandler := ruangan.InDB{DB: DB}
	ruanganRoutes := r.Group("/ruangan")
	{
		ruanganRoutes.POST("/", AuthMiddleware(), ruanganHandler.CreateRuangan)
		ruanganRoutes.GET("/gedung/:gedung", AuthMiddleware(), ruanganHandler.GetRuanganByGedung)
		ruanganRoutes.GET("/prodi/:kode_prodi", AuthMiddleware(), ruanganHandler.GetRuanganByProdi)
		ruanganRoutes.GET("/tipe/:tipe", AuthMiddleware(), ruanganHandler.GetRuanganByTipe)
		ruanganRoutes.PUT("/:kode_ruangan", AuthMiddleware(), ruanganHandler.UpdateRuangan)
		ruanganRoutes.DELETE("/:kode_ruangan", AuthMiddleware(), ruanganHandler.DeleteRuangan)
		ruanganRoutes.GET("/status/:kode_ruangan", AuthMiddleware(), ruanganHandler.GetRuanganStatus)
	}

	// Mahasiswa routes
	mahasiswaHandler := &mahasiswa.InDB{DB: DB}
	mahasiswaRoutes := r.Group("/mahasiswa")
	{
		mahasiswaRoutes.GET("/ipk/:nim", AuthMiddleware(), mahasiswaHandler.GetIPKMahasiswa)
		mahasiswaRoutes.GET("/ips/:nim/:tahun_akademik", AuthMiddleware(), mahasiswaHandler.GetIPSMahasiswa)
	}

	// Dokumentasi API (HTML view dengan CSS yang rapi dan profesional dan parsing Markdown)
	r.GET("/", func(c *gin.Context) {
		mdContent, err := os.ReadFile("API_DOCUMENTATION.md")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load documentation")
			return
		}

		// Markdown parsing options
		extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
		p := parser.NewWithExtensions(extensions)
		doc := p.Parse(mdContent)

		// HTML rendering options
		htmlFlags := html.CommonFlags | html.HrefTargetBlank
		opts := html.RendererOptions{Flags: htmlFlags}
		renderer := html.NewRenderer(opts)

		parsedHTML := markdown.Render(doc, renderer)

		// Elegant and professional HTML template
		htmlPage := `<!DOCTYPE html>
	<html lang="en">
	<head>
	  <meta charset="UTF-8" />
	  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
	  <title>API Documentation - Kelompok 1</title>
	  <link rel="icon" type="image/png" href="https://img.icons8.com/fluency/48/api-settings.png">
	  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
	  <link href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css" rel="stylesheet">
	  <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" rel="stylesheet">
	  <style>
		:root {
		  --primary-color: #6366f1;
		  --primary-dark: #4f46e5;
		  --secondary-color: #f8fafc;
		  --accent-color: #0ea5e9;
		  --text-primary: #1e293b;
		  --text-secondary: #64748b;
		  --border-color: #e2e8f0;
		  --success-color: #10b981;
		  --warning-color: #f59e0b;
		  --error-color: #ef4444;
		  --code-bg: #1e293b;
		  --sidebar-width: 280px;
		  --header-height: 80px;
		}
	
		* {
		  margin: 0;
		  padding: 0;
		  box-sizing: border-box;
		}
	
		body {
		  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
		  line-height: 1.6;
		  color: var(--text-primary);
		  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		  min-height: 100vh;
		}
	
		.container {
		  display: flex;
		  min-height: 100vh;
		}
	
		/* Sidebar */
		.sidebar {
		  width: var(--sidebar-width);
		  background: rgba(255, 255, 255, 0.95);
		  backdrop-filter: blur(10px);
		  border-right: 1px solid var(--border-color);
		  position: fixed;
		  height: 100vh;
		  overflow-y: auto;
		  z-index: 1000;
		  transition: transform 0.3s ease;
		}
	
		.sidebar-header {
		  padding: 2rem 1.5rem;
		  border-bottom: 1px solid var(--border-color);
		  background: var(--primary-color);
		  color: white;
		  text-align: center;
		}
	
		.sidebar-header .logo {
		  width: 48px;
		  height: 48px;
		  margin-bottom: 1rem;
		  filter: brightness(0) invert(1);
		}
	
		.sidebar-header h1 {
		  font-size: 1.25rem;
		  font-weight: 600;
		  margin-bottom: 0.5rem;
		}
	
		.sidebar-header p {
		  font-size: 0.875rem;
		  opacity: 0.9;
		}
	
		.nav-menu {
		  padding: 1.5rem 0;
		}
	
		.nav-item {
		  padding: 0.75rem 1.5rem;
		  cursor: pointer;
		  transition: all 0.2s ease;
		  border-left: 3px solid transparent;
		}
	
		.nav-item:hover {
		  background: var(--secondary-color);
		  border-left-color: var(--primary-color);
		}
	
		.nav-item.active {
		  background: var(--secondary-color);
		  border-left-color: var(--primary-color);
		  color: var(--primary-color);
		  font-weight: 500;
		}
	
		.nav-item i {
		  margin-right: 0.75rem;
		  width: 16px;
		}
	
		/* Main Content */
		.main-content {
		  flex: 1;
		  margin-left: var(--sidebar-width);
		  background: white;
		  min-height: 100vh;
		}
	
		.content-header {
		  background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
		  color: white;
		  padding: 3rem 2rem;
		  text-align: center;
		  position: relative;
		  overflow: hidden;
		}
	
		.content-header::before {
		  content: '';
		  position: absolute;
		  top: 0;
		  left: 0;
		  right: 0;
		  bottom: 0;
		  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grid" width="10" height="10" patternUnits="userSpaceOnUse"><path d="M 10 0 L 0 0 0 10" fill="none" stroke="rgba(255,255,255,0.1)" stroke-width="1"/></pattern></defs><rect width="100" height="100" fill="url(%23grid)"/></svg>');
		  opacity: 0.1;
		}
	
		.content-header h1 {
		  font-size: 2.5rem;
		  font-weight: 700;
		  margin-bottom: 1rem;
		  position: relative;
		  z-index: 1;
		}
	
		.content-header p {
		  font-size: 1.1rem;
		  opacity: 0.9;
		  position: relative;
		  z-index: 1;
		}
	
		.content-body {
		  padding: 2rem;
		}
	
		/* Info Cards */
		.info-grid {
		  display: grid;
		  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		  gap: 1.5rem;
		  margin-bottom: 3rem;
		}
	
		.info-card {
		  background: white;
		  border: 1px solid var(--border-color);
		  border-radius: 12px;
		  padding: 1.5rem;
		  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
		  transition: all 0.3s ease;
		}
	
		.info-card:hover {
		  transform: translateY(-2px);
		  box-shadow: 0 8px 25px -5px rgba(0, 0, 0, 0.1);
		}
	
		.info-card-header {
		  display: flex;
		  align-items: center;
		  margin-bottom: 1rem;
		}
	
		.info-card-icon {
		  width: 40px;
		  height: 40px;
		  background: var(--primary-color);
		  border-radius: 8px;
		  display: flex;
		  align-items: center;
		  justify-content: center;
		  margin-right: 1rem;
		}
	
		.info-card-icon i {
		  color: white;
		  font-size: 1.25rem;
		}
	
		.info-card h3 {
		  font-size: 1.125rem;
		  font-weight: 600;
		  color: var(--text-primary);
		}
	
		.info-list {
		  list-style: none;
		}
	
		.info-list li {
		  padding: 0.5rem 0;
		  border-bottom: 1px solid var(--border-color);
		  display: flex;
		  align-items: center;
		}
	
		.info-list li:last-child {
		  border-bottom: none;
		}
	
		.info-list li i {
		  margin-right: 0.75rem;
		  color: var(--success-color);
		  width: 16px;
		}
	
		/* Documentation Content */
		.documentation-content {
		  background: white;
		  border-radius: 12px;
		  padding: 2rem;
		  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
		  margin-top: 2rem;
		}
	
		.documentation-content h1,
		.documentation-content h2,
		.documentation-content h3,
		.documentation-content h4,
		.documentation-content h5,
		.documentation-content h6 {
		  color: var(--text-primary);
		  margin-top: 2rem;
		  margin-bottom: 1rem;
		  padding-bottom: 0.5rem;
		  border-bottom: 2px solid var(--border-color);
		}
	
		.documentation-content h1 {
		  font-size: 2rem;
		  color: var(--primary-color);
		}
	
		.documentation-content h2 {
		  font-size: 1.5rem;
		}
	
		.documentation-content pre {
		  background: var(--code-bg);
		  color: #f8f8f2;
		  padding: 1.5rem;
		  border-radius: 8px;
		  overflow-x: auto;
		  font-family: 'Fira Code', Consolas, 'Courier New', monospace;
		  font-size: 0.875rem;
		  line-height: 1.5;
		  position: relative;
		}
	
		.documentation-content code {
		  font-family: 'Fira Code', Consolas, 'Courier New', monospace;
		  background: #f1f5f9;
		  padding: 0.25rem 0.5rem;
		  border-radius: 4px;
		  font-size: 0.875rem;
		  color: var(--primary-color);
		}
	
		.documentation-content pre code {
		  background: transparent;
		  padding: 0;
		  color: inherit;
		}
	
		.documentation-content a {
		  color: var(--accent-color);
		  text-decoration: none;
		  font-weight: 500;
		  transition: color 0.2s ease;
		}
	
		.documentation-content a:hover {
		  color: var(--primary-color);
		  text-decoration: underline;
		}
	
		.documentation-content blockquote {
		  border-left: 4px solid var(--primary-color);
		  background: var(--secondary-color);
		  padding: 1rem 1.5rem;
		  margin: 1.5rem 0;
		  border-radius: 0 8px 8px 0;
		}
	
		.documentation-content table {
		  width: 100%;
		  border-collapse: collapse;
		  margin: 1.5rem 0;
		  background: white;
		  border-radius: 8px;
		  overflow: hidden;
		  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
		}
	
		.documentation-content th,
		.documentation-content td {
		  padding: 1rem;
		  text-align: left;
		  border-bottom: 1px solid var(--border-color);
		}
	
		.documentation-content th {
		  background: var(--secondary-color);
		  font-weight: 600;
		  color: var(--text-primary);
		}
	
		.documentation-content tr:hover {
		  background: var(--secondary-color);
		}
	
		/* Copy Button */
		.copy-button {
		  position: absolute;
		  top: 1rem;
		  right: 1rem;
		  background: rgba(255, 255, 255, 0.1);
		  color: white;
		  border: 1px solid rgba(255, 255, 255, 0.2);
		  padding: 0.5rem 1rem;
		  border-radius: 6px;
		  cursor: pointer;
		  font-size: 0.75rem;
		  font-weight: 500;
		  transition: all 0.2s ease;
		  backdrop-filter: blur(4px);
		}
	
		.copy-button:hover {
		  background: rgba(255, 255, 255, 0.2);
		  transform: translateY(-1px);
		}
	
		.copy-button.copied {
		  background: var(--success-color);
		  border-color: var(--success-color);
		}
	
		/* Status badges */
		.status-badge {
		  display: inline-block;
		  padding: 0.25rem 0.75rem;
		  border-radius: 20px;
		  font-size: 0.75rem;
		  font-weight: 500;
		  text-transform: uppercase;
		  letter-spacing: 0.5px;
		}
	
		.status-success {
		  background: rgba(16, 185, 129, 0.1);
		  color: var(--success-color);
		  border: 1px solid var(--success-color);
		}
	
		.status-warning {
		  background: rgba(245, 158, 11, 0.1);
		  color: var(--warning-color);
		  border: 1px solid var(--warning-color);
		}
	
		.status-error {
		  background: rgba(239, 68, 68, 0.1);
		  color: var(--error-color);
		  border: 1px solid var(--error-color);
		}
	
		/* Responsive Design */
		@media (max-width: 768px) {
		  .sidebar {
			transform: translateX(-100%);
		  }
	
		  .sidebar.open {
			transform: translateX(0);
		  }
	
		  .main-content {
			margin-left: 0;
		  }
	
		  .content-header h1 {
			font-size: 2rem;
		  }
	
		  .content-body {
			padding: 1rem;
		  }
	
		  .info-grid {
			grid-template-columns: 1fr;
		  }
		}
	
		/* Scrollbar Styling */
		::-webkit-scrollbar {
		  width: 8px;
		}
	
		::-webkit-scrollbar-track {
		  background: var(--secondary-color);
		}
	
		::-webkit-scrollbar-thumb {
		  background: var(--border-color);
		  border-radius: 4px;
		}
	
		::-webkit-scrollbar-thumb:hover {
		  background: var(--text-secondary);
		}
	
		/* Loading Animation */
		.loading {
		  display: inline-block;
		  width: 20px;
		  height: 20px;
		  border: 3px solid rgba(255, 255, 255, 0.3);
		  border-radius: 50%;
		  border-top-color: white;
		  animation: spin 1s ease-in-out infinite;
		}
	
		@keyframes spin {
		  to { transform: rotate(360deg); }
		}
	  </style>
	</head>
	<body>
	  <div class="container">
		<!-- Sidebar -->
		<div class="sidebar" id="sidebar">
		  <div class="sidebar-header">
			<img src="https://img.icons8.com/fluency/48/api-settings.png" alt="API Icon" class="logo"/>
			<h1>API Docs</h1>
			<p>Kelompok1 System</p>
		  </div>
		  <nav class="nav-menu">
			<div class="nav-item active" onclick="scrollToSection('overview')">
			  <i class="fas fa-home"></i>
			  Overview
			</div>
			<div class="nav-item" onclick="scrollToSection('authentication')">
			  <i class="fas fa-lock"></i>
			  Authentication
			</div>
			<div class="nav-item" onclick="scrollToSection('endpoints')">
			  <i class="fas fa-plug"></i>
			  Endpoints
			</div>
			<div class="nav-item" onclick="scrollToSection('examples')">
			  <i class="fas fa-code"></i>
			  Examples
			</div>
			<div class="nav-item" onclick="scrollToSection('errors')">
			  <i class="fas fa-exclamation-triangle"></i>
			  Error Handling
			</div>
		  </nav>
		</div>
	
		<!-- Main Content -->
		<div class="main-content">
		  <div class="content-header">
			<h1>API Documentation</h1>
			<p>Sistem Kepegawaian - Comprehensive REST API Guide</p>
		  </div>
	
		  <div class="content-body">
			<!-- Info Cards -->
			<div class="info-grid">
			  <div class="info-card">
				<div class="info-card-header">
				  <div class="info-card-icon">
					<i class="fas fa-server"></i>
				  </div>
				  <h3>Base Information</h3>
				</div>
				<ul class="info-list">
				  <li>
					<i class="fas fa-check"></i>
					<strong>Base URL:</strong> <a href="https://ti054c01.agussbn.my.id" target="_blank">https://ti054c01.agussbn.my.id</a>
				  </li>
				  <li>
					<i class="fas fa-check"></i>
					<strong>Format:</strong> JSON (form-data untuk upload)
				  </li>
				  <li>
					<i class="fas fa-check"></i>
					<strong>Version:</strong> v1.0
				  </li>
				</ul>
			  </div>
	
			  <div class="info-card">
				<div class="info-card-header">
				  <div class="info-card-icon">
					<i class="fas fa-shield-alt"></i>
				  </div>
				  <h3>Authentication</h3>
				</div>
				<ul class="info-list">
				  <li>
					<i class="fas fa-check"></i>
					<strong>Type:</strong> JWT Token
				  </li>
				  <li>
					<i class="fas fa-check"></i>
					<strong>Header:</strong> Authorization
				  </li>
				  <li>
					<i class="fas fa-check"></i>
					<strong>Role-based:</strong> Access Control
				  </li>
				</ul>
			  </div>
	
			  <div class="info-card">
				<div class="info-card-header">
				  <div class="info-card-icon">
					<i class="fas fa-cloud-upload-alt"></i>
				  </div>
				  <h3>File Upload</h3>
				</div>
				<ul class="info-list">
				  <li>
					<i class="fas fa-check"></i>
					<strong>Available:</strong> Specific endpoints
				  </li>
				  <li>
					<i class="fas fa-check"></i>
					<strong>Storage:</strong> /uploads folder
				  </li>
				  <li>
					<i class="fas fa-check"></i>
					<strong>Format:</strong> Multipart form-data
				  </li>
				</ul>
			  </div>
			</div>
	
			<!-- Authentication Status Info -->
			<div class="info-card">
			  <div class="info-card-header">
				<div class="info-card-icon">
				  <i class="fas fa-info-circle"></i>
				</div>
				<h3>Response Codes</h3>
			  </div>
			  <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem; margin-top: 1rem;">
				<div style="padding: 1rem; background: rgba(239, 68, 68, 0.1); border-radius: 8px; border-left: 4px solid var(--error-color);">
				  <div style="display: flex; align-items: center; margin-bottom: 0.5rem;">
					<span class="status-badge status-error">401</span>
					<span style="margin-left: 0.5rem; font-weight: 500;">No Token</span>
				  </div>
				  <p style="font-size: 0.875rem; color: var(--text-secondary);">Token tidak ditemukan</p>
				</div>
				<div style="padding: 1rem; background: rgba(245, 158, 11, 0.1); border-radius: 8px; border-left: 4px solid var(--warning-color);">
				  <div style="display: flex; align-items: center; margin-bottom: 0.5rem;">
					<span class="status-badge status-warning">403</span>
					<span style="margin-left: 0.5rem; font-weight: 500;">Invalid Token</span>
				  </div>
				  <p style="font-size: 0.875rem; color: var(--text-secondary);">Token tidak valid</p>
				</div>
				<div style="padding: 1rem; background: rgba(239, 68, 68, 0.1); border-radius: 8px; border-left: 4px solid var(--error-color);">
				  <div style="display: flex; align-items: center; margin-bottom: 0.5rem;">
					<span class="status-badge status-error">Access Denied</span>
				  </div>
				  <p style="font-size: 0.875rem; color: var(--text-secondary);">Role tidak memiliki akses</p>
				  <div style="margin-top: 1rem; background: var(--code-bg); color: white; padding: 1rem; border-radius: 6px; position: relative;">
					<code>{ "message": "Access Denied!" }</code>
					<button class="copy-button" onclick="copyToClipboard(this, '{ \"message\": \"Access Denied!\" }')">
					  <i class="fas fa-copy"></i>
					</button>
				  </div>
				</div>
			  </div>
			</div>
	
			<!-- Documentation Content -->
			<div class="documentation-content" id="documentation">
			  ` + string(parsedHTML) + `
			</div>
		  </div>
		</div>
	  </div>
	
	  <script>
		// Copy to clipboard functionality
		function copyToClipboard(button, text) {
		  const textToCopy = text || button.previousElementSibling.textContent;
		  navigator.clipboard.writeText(textToCopy).then(() => {
			const originalContent = button.innerHTML;
			button.innerHTML = '<i class="fas fa-check"></i>';
			button.classList.add('copied');
			setTimeout(() => {
			  button.innerHTML = originalContent;
			  button.classList.remove('copied');
			}, 2000);
		  }).catch(err => {
			console.error('Failed to copy: ', err);
		  });
		}
	
		// Smooth scroll to section
		function scrollToSection(sectionId) {
		  const element = document.getElementById(sectionId);
		  if (element) {
			element.scrollIntoView({ behavior: 'smooth' });
		  }
		  
		  // Update active nav item
		  document.querySelectorAll('.nav-item').forEach(item => {
			item.classList.remove('active');
		  });
		  const activeItem = document.querySelector('[onclick="scrollToSection(\'' + sectionId + '\')"]');
		  if (activeItem) {
			activeItem.classList.add('active');
		  }
		}
	
		// Mobile sidebar toggle
		function toggleSidebar() {
		  const sidebar = document.getElementById('sidebar');
		  sidebar.classList.toggle('open');
		}
	
		// Add copy buttons to all pre elements
		document.addEventListener('DOMContentLoaded', function() {
		  const preElements = document.querySelectorAll('pre');
		  preElements.forEach(pre => {
			if (!pre.querySelector('.copy-button')) {
			  const button = document.createElement('button');
			  button.className = 'copy-button';
			  button.innerHTML = '<i class="fas fa-copy"></i>';
			  button.onclick = function() {
				copyToClipboard(this, pre.textContent);
			  };
			  pre.appendChild(button);
			}
		  });
		});
	
		// Intersection Observer for active navigation
		const observer = new IntersectionObserver((entries) => {
		  entries.forEach(entry => {
			if (entry.isIntersecting) {
			  const id = entry.target.id;
			  document.querySelectorAll('.nav-item').forEach(item => {
				item.classList.remove('active');
			  });
			  const activeItem = document.querySelector('[onclick="scrollToSection(\'' + id + '\')"]');
			  if (activeItem) {
				activeItem.classList.add('active');
			  }
			}
		  });
		}, { threshold: 0.5 });
	
		// Observe sections
		document.addEventListener('DOMContentLoaded', function() {
		  const sections = document.querySelectorAll('[id]');
		  sections.forEach(section => {
			observer.observe(section);
		  });
		});
	  </script>
	</body>
	</html>`

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, htmlPage)
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

// getJWTKey returns the JWT key from environment or uses a fallback
func getJWTKey() []byte {
	key := os.Getenv("simpadu_jwt_key")
	if key == "" {
		// Only for development, warn about using default key
		log.Println("Warning: Using default JWT key. Set simpadu_jwt_key environment variable in production.")
		return []byte("secret-key-boleh-diubah")
	}
	return []byte(key)
}

func RegisterUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		NIM      string `json:"nim" binding:"required"`
		Nama     string `json:"nama" binding:"required"`
		Angkatan uint   `json:"angkatan" binding:"required"`
		Prodi    string `json:"prodi" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format registrasi tidak valid"})
		return
	}

	// Validasi token
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Validasi role (admin prodi atau admin akademik)
	if claims["role"] != "admin_prodi" && claims["role"] != "admin_akademik" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	// Validasi prodi jika admin prodi
	if claims["role"] == "admin_prodi" && claims["kode_prodi"] != input.Prodi {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda hanya dapat mendaftarkan mahasiswa di prodi Anda"})
		return
	}

	// Validasi username unik
	var existingUser structs.Users
	if err := DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
		return
	}

	// Validasi email unik
	if err := DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
		return
	}

	// Validasi NIM unik
	var existingMahasiswa structs.Mahasiswa
	if err := DB.Where("nim = ?", input.NIM).First(&existingMahasiswa).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIM sudah terdaftar"})
		return
	}

	// Validasi prodi exists
	var prodi structs.Prodi
	if err := DB.Where("kode_prodi = ?", input.Prodi).First(&prodi).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prodi tidak ditemukan"})
		return
	}

	// Generate user_id
	userID := helper.GenerateUserID("MHS")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// Buat user baru
	user := structs.Users{
		User_id:  userID,
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "mahasiswa",
		Status:   "Aktif",
	}

	// Buat mahasiswa baru
	mahasiswa := structs.Mahasiswa{
		User_id:    userID,
		NIM:        input.NIM,
		Nama:       input.Nama,
		Angkatan:   input.Angkatan,
		Kode_prodi: input.Prodi,
	}

	// Simpan ke database dalam satu transaksi
	tx := DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan user"})
		return
	}

	if err := tx.Create(&mahasiswa).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan mahasiswa"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Registrasi berhasil",
		"data": gin.H{
			"user_id":  userID,
			"username": input.Username,
			"email":    input.Email,
			"nim":      input.NIM,
			"nama":     input.Nama,
			"prodi":    input.Prodi,
		},
	})
}

func RegisterAdminProdi(c *gin.Context) {
	var input struct {
		Username  string    `json:"username" binding:"required"`
		Email     string    `json:"email" binding:"required,email"`
		Password  string    `json:"password" binding:"required,min=6"`
		Nama      string    `json:"nama" binding:"required"`
		NIP       string    `json:"nip" binding:"required"`
		Prodi     string    `json:"prodi" binding:"required"`
		Jabatan   string    `json:"jabatan" binding:"required"`
		NoSK      string    `json:"no_sk" binding:"required"`
		TanggalSK time.Time `json:"tanggal_sk" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format registrasi tidak valid"})
		return
	}

	// Validasi token
	AuthorizationHeader := c.GetHeader("Authorization")
	claims := helper.ExtractToken(strings.Replace(AuthorizationHeader, "Bearer ", "", -1))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Validasi role (hanya admin akademik)
	if claims["role"] != "admin_akademik" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin akademik yang dapat mendaftarkan admin prodi"})
		return
	}

	// Validasi username unik
	var existingUser structs.Users
	if err := DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
		return
	}

	// Validasi email unik
	if err := DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
		return
	}

	// Validasi NIP unik
	var existingAdmin structs.Struktural_prodi
	if err := DB.Where("NIP = ?", input.NIP).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIP sudah terdaftar"})
		return
	}

	// Validasi prodi exists
	var prodi structs.Prodi
	if err := DB.Where("kode_prodi = ?", input.Prodi).First(&prodi).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prodi tidak ditemukan"})
		return
	}

	// Generate user_id
	userID := helper.GenerateUserID("ADM")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// Buat user baru
	user := structs.Users{
		User_id:  userID,
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "admin_prodi",
		Status:   "Aktif",
	}

	// Buat admin prodi baru
	adminProdi := structs.Struktural_prodi{
		User_id:    userID,
		Kode_prodi: input.Prodi,
		Nama:       input.Nama,
		NIP:        input.NIP,
		Jabatan:    input.Jabatan,
		NoSK:       input.NoSK,
		TanggalSK:  input.TanggalSK,
	}

	// Simpan ke database dalam satu transaksi
	tx := DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan user"})
		return
	}

	if err := tx.Create(&adminProdi).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan admin prodi"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Registrasi berhasil",
		"data": gin.H{
			"user_id":  userID,
			"username": input.Username,
			"email":    input.Email,
			"nama":     input.Nama,
			"nip":      input.NIP,
			"prodi":    input.Prodi,
			"jabatan":  input.Jabatan,
		},
	})
}
func Login(c *gin.Context) {
	var LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&LoginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format login tidak valid"})
		return
	}

	var user structs.Users
	if err := DB.Where("username = ?", LoginRequest.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username tidak ditemukan"})
		return
	}

	// Validasi password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// Mulai bangun JWT claims
	claims := TokenClaims{
		UserID: user.User_id,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // Token berlaku 1 jam
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Tambahkan informasi tambahan berdasarkan role
	switch user.Role {
	case "mahasiswa":
		var mahasiswa structs.Mahasiswa
		if err := DB.Where("user_id = ?", user.User_id).First(&mahasiswa).Error; err == nil {
			claims.Kode_prodi = mahasiswa.Kode_prodi
			claims.Nama = mahasiswa.Nama
			claims.NIM = mahasiswa.NIM
		}
	case "admin_prodi":
		var adminProdi structs.Struktural_prodi
		if err := DB.Where("user_id = ?", user.User_id).First(&adminProdi).Error; err == nil {
			claims.Kode_prodi = adminProdi.Kode_prodi
			claims.Nama = adminProdi.Nama
			claims.NIP = adminProdi.NIP
		}
	case "dosen":
		var dosen structs.Dosen
		if err := DB.Where("user_id = ?", user.User_id).First(&dosen).Error; err == nil {
			claims.Kode_prodi = dosen.Kode_prodi
			claims.NIP = dosen.NIP
			// Untuk dosen, kita perlu mengambil nama dari tabel struktural_prodi
			var struktural structs.Struktural_prodi
			if err := DB.Where("NIP = ?", dosen.NIP).First(&struktural).Error; err == nil {
				claims.Nama = struktural.Nama
			}
		}
	case "admin_akademik":
		claims.Nama = "Admin Akademik"
	case "staff_keuangan":
		claims.Nama = "Staff Keuangan"
	case "staff_tu":
		claims.Nama = "Staff Tata Usaha"
	case "staff_rektorat":
		claims.Nama = "Staff Rektorat"
	default:
		claims.Nama = "Pengguna Sistem"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"user_id": claims.UserID,
			"email":   claims.Email,
			"role":    claims.Role,
			"nama":    claims.Nama,
		},
	})
}

func Profile(c *gin.Context) {
	// Ambil user_id dari context (sudah diset oleh AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// Ambil role dari context
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	// Ambil data user
	var user structs.Users
	if err := DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Siapkan response
	response := structs.User_profile_response{
		User_id:  user.User_id,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
		Roles:    make([]structs.User_role_detail, 0),
	}

	// Tambahkan informasi tambahan berdasarkan role
	switch role {
	case "mahasiswa":
		var mahasiswa structs.Mahasiswa
		if err := DB.Where("user_id = ?", userID).First(&mahasiswa).Error; err == nil {
			response.Nama = mahasiswa.Nama
			response.Roles = append(response.Roles, structs.User_role_detail{
				Role_code:  "mahasiswa",
				Role_name:  "Mahasiswa",
				Kode_prodi: mahasiswa.Kode_prodi,
				Status:     user.Status,
			})
		}
	case "admin_prodi":
		var adminProdi structs.Struktural_prodi
		if err := DB.Where("user_id = ?", userID).First(&adminProdi).Error; err == nil {
			response.Nama = adminProdi.Nama
			response.Roles = append(response.Roles, structs.User_role_detail{
				Role_code:  "admin_prodi",
				Role_name:  "Admin Prodi",
				Kode_prodi: adminProdi.Kode_prodi,
				Status:     user.Status,
			})
		}
	case "dosen":
		var dosen structs.Dosen
		if err := DB.Where("user_id = ?", userID).First(&dosen).Error; err == nil {
			// Untuk dosen, kita perlu mengambil nama dari tabel struktural_prodi
			var struktural structs.Struktural_prodi
			if err := DB.Where("n_ip = ?", dosen.NIP).First(&struktural).Error; err == nil {
				response.Nama = struktural.Nama
			}
			response.Roles = append(response.Roles, structs.User_role_detail{
				Role_code:  "dosen",
				Role_name:  "Dosen",
				Kode_prodi: dosen.Kode_prodi,
				Status:     user.Status,
			})
		}
	case "admin_akademik":
		response.Roles = append(response.Roles, structs.User_role_detail{
			Role_code: "admin_akademik",
			Role_name: "Admin Akademik",
			Status:    user.Status,
		})
	case "admin":
		response.Roles = append(response.Roles, structs.User_role_detail{
			Role_code: "admin",
			Role_name: "Administrator",
			Status:    user.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile berhasil diambil",
		"data":    response,
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak disediakan"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi metode enkripsi
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode token tidak valid")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		// Ambil klaim dan inject ke context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Validasi claims yang diperlukan
			requiredClaims := []string{"user_id", "role", "email"}
			for _, claim := range requiredClaims {
				if _, exists := claims[claim]; !exists {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak memiliki informasi yang diperlukan"})
					c.Abort()
					return
				}
			}

			// Validasi role
			role := claims["role"].(string)
			validRoles := map[string]bool{
				"admin":          true,
				"admin_akademik": true,
				"admin_prodi":    true,
				"dosen":          true,
				"mahasiswa":      true,
			}
			if !validRoles[role] {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Role tidak valid"})
				c.Abort()
				return
			}

			// Set claims ke context
			c.Set("user_id", claims["user_id"])
			c.Set("role", claims["role"])
			c.Set("email", claims["email"])
			c.Set("kode_prodi", claims["kode_prodi"])
			c.Set("nama", claims["nama"])
			c.Set("n_ip", claims["n_ip"])
			c.Set("nim", claims["nim"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal memproses token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func MenuAccessMiddleware(menu string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsRaw, exists := c.Get("claims")
		if !exists {
			c.JSON(401, gin.H{"message": "claims token tidak ditemukan"})
			c.Abort()
			return
		}

		claims, ok := claimsRaw.(*helper.TokenClaims)
		if !ok {
			c.JSON(401, gin.H{"message": "claims token tidak valid"})
			c.Abort()
			return
		}

		role := claims.Role // pastikan kamu sudah simpan role di token claims
		if !helper.CanAccessMenu(role, menu) {
			c.JSON(403, gin.H{"message": "akses menu ini tidak diizinkan untuk role " + role})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RandString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
