package models

import "time"

type Form struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Status      string     `gorm:"not null;default:'open'" json:"status"` // "open" or "closed"
	OwnerID     uint       `gorm:"not null" json:"owner_id"`
	Owner       User       `gorm:"foreignKey:OwnerID" json:"-"`
	Questions   []Question `gorm:"foreignKey:FormID" json:"questions,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Question struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	FormID   uint   `gorm:"not null" json:"form_id"`
	Label    string `gorm:"not null" json:"label"`
	Type     string `gorm:"not null;default:'short_answer'" json:"type"` // short_answer, radio, checkbox, dropdown
	Options  string `json:"options"`                                     // JSON-encoded options for radio/checkbox/dropdown
	Required bool   `gorm:"not null;default:true" json:"required"`
	Order    int    `gorm:"not null;default:0" json:"order"`
}

// ValidQuestionTypes contains all allowed question types.
var ValidQuestionTypes = map[string]bool{
	"short_answer": true,
	"radio":        true,
	"checkbox":     true,
	"dropdown":     true,
}

type Response struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FormID    uint      `gorm:"not null" json:"form_id"`
	Form      Form      `gorm:"foreignKey:FormID" json:"-"`
	Answers   []Answer  `gorm:"foreignKey:ResponseID" json:"answers,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Answer struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ResponseID uint   `gorm:"not null" json:"response_id"`
	QuestionID uint   `gorm:"not null" json:"question_id"`
	Value      string `gorm:"not null" json:"value"` // plain text; for checkbox use JSON array string
}
