package models

import "time"

type Message struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" validate:"required,min=1,max=255" gorm:"not null"`
	IsDone    bool      `json:"is_done" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
