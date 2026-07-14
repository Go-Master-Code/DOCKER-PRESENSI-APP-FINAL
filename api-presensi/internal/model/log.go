package model

import "time"

type Log struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	Method    string    `json:"method"`
	Endpoint  string    `json:"endpoint"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Log) TableName() string {
	return "log"
}
