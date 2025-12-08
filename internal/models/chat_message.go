package models

type ChatMessage struct {
	Base
	UserID uint   `json:"user_id"`
	Text   string `json:"text"`

	ChatableID   uint   `json:"chatable_id"`
	ChatableType string `json:"chatable_type"`
}

type ChatMessageCreateReq struct {
	UserID uint   `json:"user_id" binding:"required"`
	Text   string `json:"text" binding:"required,min=1,max=5000"`

	ChatableID   uint   `json:"chatable_id" binding:"required"`
	ChatableType string `json:"chatable_type" binding:"required,oneof=projects tasks"`
}
