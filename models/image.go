package models

type Image struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}
