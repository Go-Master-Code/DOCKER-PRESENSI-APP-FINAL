package model

type Role struct {
	ID        int    `json:"id"  gorm:"primaryKey"`
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
	User      []User //reverse relation: 1 to many (1 role bisa milik beberapa user)
}

func (Role) TableName() string {
	return "role"
}
