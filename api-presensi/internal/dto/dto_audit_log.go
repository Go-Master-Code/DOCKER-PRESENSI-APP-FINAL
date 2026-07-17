package dto

type LogRequestAndResponse struct {
	UserID    string `json:"user_id"`
	Method    string `json:"method"`
	Endpoint  string `json:"endpoint"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	CreatedAt string `json:"created_at"`
}
