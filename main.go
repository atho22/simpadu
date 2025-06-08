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
	N_ip       string `json:"n_ip"`
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
	if err := DB.AutoMigrate(&structs.Users{}, &structs.Tahun_akademik{}, &structs.Kelas{}, &structs.Mahasiswa{}, &structs.Kelas_mahasiswa{}, &structs.Jurusan{}, &structs.Prodi{}, &structs.Struktural_prodi{}, &structs.Mata_kuliah{}, &structs.Jadwal_matakuliah{}, &structs.Absensi{}, &structs.Penilaian{}, &structs.Kode_ruangan{}, &structs.Pertemuan{}, &structs.Dosen{}); err != nil {
		log.Fatal("Database migration failed:", err)
	}
	log.Println("Database migration completed successfully")

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
	r.GET("/profile", AuthMiddleware(), Profile)

	// Kelas routes
	r.GET("/kelas", AuthMiddleware(), kelasHandler.GetAllKelas)
	r.GET("/kelas/:id", AuthMiddleware(), kelasHandler.DetailKelas)
	r.POST("/create-kelas", AuthMiddleware(), kelasHandler.CreateKelas)

	// Tahun Akademik routes
	r.GET("/tahun-akademik", AuthMiddleware(), tahunAkademikHandler.GetAllTahunAkademik)
	r.POST("/create-tahun-akademik", AuthMiddleware(), tahunAkademikHandler.CreateTahunAkademik)
	r.GET("/tahun-akademik/:id", AuthMiddleware(), tahunAkademikHandler.GetTahunAkademikByID)

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
		N_ip      string    `json:"n_ip" binding:"required"`
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
	if err := DB.Where("n_ip = ?", input.N_ip).First(&existingAdmin).Error; err == nil {
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
		N_ip:       input.N_ip,
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
			"n_ip":     input.N_ip,
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
			claims.N_ip = adminProdi.N_ip
		}
	case "dosen":
		var dosen structs.Dosen
		if err := DB.Where("user_id = ?", user.User_id).First(&dosen).Error; err == nil {
			claims.Kode_prodi = dosen.Kode_prodi
			claims.N_ip = dosen.NIP
			// Untuk dosen, kita perlu mengambil nama dari tabel struktural_prodi
			var struktural structs.Struktural_prodi
			if err := DB.Where("n_ip = ?", dosen.NIP).First(&struktural).Error; err == nil {
				claims.Nama = struktural.Nama
			}
		}
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
