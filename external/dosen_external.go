package external

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Dosen struct {
	IDPegawai   string `json:"id_pegawai"`
	NamaPegawai string `json:"nama_pegawai"`
}

type DosenResponse struct {
	Success bool    `json:"success"`
	Data    []Dosen `json:"data"`
}

// IsDosenValid memeriksa apakah kode dosen tersedia di sistem dosen kelompok 2
func IsDosenValid(kodeDosen string) (bool, error) {
	url := "https://ti054c02.agussbn.my.id/api/data/pegawai/dosen"

	resp, err := http.Get(url)
	if err != nil {
		return false, errors.New("gagal mengakses API dosen")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.New("gagal membaca response dari API dosen")
	}

	var response DosenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return false, errors.New("gagal parsing JSON dari API dosen")
	}

	if !response.Success || len(response.Data) == 0 {
		return false, nil
	}

	for _, dosen := range response.Data {
		if dosen.IDPegawai == kodeDosen {
			return true, nil
		}
	}

	return false, nil
}
