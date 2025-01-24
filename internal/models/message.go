package models

type Message struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Text string `json:"text"`
}
