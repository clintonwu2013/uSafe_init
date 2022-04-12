package models

type SettingServer struct {
	Port             int    `json:"Port"`
	CertFilePath     string `json:"CertFilePath"`
	CertFilePassword string `json:"CertFilePassword"`
}
