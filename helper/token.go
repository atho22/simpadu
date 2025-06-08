package helper

import (
	"fmt"
	"time"
)

// GenerateKodePertemuan generates a unique meeting code
func GenerateKodePertemuan(kodeMatakuliah string, pertemuanKe int) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d-%d", kodeMatakuliah, pertemuanKe, timestamp)
}
