package dto

type MfaData struct {
	Secret    string `json:"secret"`
	Qrcodeurl string `json:"qrcodeurl"`
}
