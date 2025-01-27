package models

import "time"

type Task struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	TaskName        string    `json:"task_name" gorm:"not null"`
	TaskDescription string    `json:"task_description"`
	IsDone          bool      `json:"is_done" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	StatusUpdatedAt time.Time `json:"status_updated_at" gorm:"autoUpdateTime"`
}
