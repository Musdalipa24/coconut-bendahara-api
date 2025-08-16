package dto

type IuranRequest struct {
	Periode      string `json:"periode"`
	MingguKe     int    `json:"minggu_ke"`
	TanggalBayar string `json:"tanggal_bayar"`
	Status       string `json:"status"`
	JumlahBayar  int64  `json:"jumlah_bayar"`
}

type IuranResponse struct {
	IdIuran      string `json:"id_iuran"`
	Periode      string `json:"periode"`
	MingguKe     int    `json:"minggu_ke"`
	Status       string `json:"status"`
	JumlahBayar  int64  `json:"jumlah_bayar"`
	TanggalBayar string `json:"tanggal_bayar"`
}

type MemberRequest struct {
	NRA    string `json:"nra"`
	Nama   string `json:"nama"`
	Status string `json:"status"`
}

type MemberResponse struct {
	IdMember  string          `json:"id_member"`
	NRA       string          `json:"nra"`
	Nama      string          `json:"nama"`
	Status    string          `json:"status"`
	Iuran     []IuranResponse `json:"iuran"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}
