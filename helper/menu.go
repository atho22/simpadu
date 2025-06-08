package helper

const (
	MenuDashboard     = "dashboard"
	MenuKelas         = "kelas"
	MenuTahunAkademik = "tahun_akademik"
	MenuProfil        = "profil"
	MenuAdminPanel    = "admin_panel"
)

var RoleMenuAccess = map[string][]string{
	"admin_akademik": {MenuDashboard, MenuKelas, MenuTahunAkademik, MenuProfil, MenuAdminPanel},
	"admin_prodi":    {MenuDashboard, MenuKelas, MenuTahunAkademik, MenuProfil, MenuAdminPanel},
	"dosen":          {MenuDashboard, MenuKelas, MenuProfil},
	"mahasiswa":      {MenuDashboard, MenuProfil},
}

// Fungsi cek apakah role boleh akses menu tertentu
func CanAccessMenu(role string, menu string) bool {
	allowedMenus, exists := RoleMenuAccess[role]
	if !exists {
		return false
	}
	for _, m := range allowedMenus {
		if m == menu {
			return true
		}
	}
	return false
}
