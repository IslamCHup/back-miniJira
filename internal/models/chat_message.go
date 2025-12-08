package models

type ChatMessage struct {
	Base
	UserID uint   `json:"user_id"`
	Text   string `json:"text"`

	ChatableID   uint   `json:"chatable_id"`
	ChatableType string `json:"chatable_type"`
}
