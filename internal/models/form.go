package models

import "time"

type Form struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	OwnerID     uint       `gorm:"not null" json:"owner_id"`
	Owner       User       `gorm:"foreignKey:OwnerID" json:"-"`
	Questions   []Question `gorm:"foreignKey:FormID" json:"questions,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Question struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	FormID  uint   `gorm:"not null" json:"form_id"`
	Label   string `gorm:"not null" json:"label"`
	Type    string `gorm:"not null;default:'text'" json:"type"` // text, radio, checkbox, dropdown
	Options string `json:"options"`                             // JSON-encoded options for radio/checkbox/dropdown
	Order   int    `gorm:"not null;default:0" json:"order"`
}
