package models

type FtpImage struct {
	ID   string `gorm:"not null" json:"id"`
	Name string `gorm:"primaryKey" json:"name"`
}
