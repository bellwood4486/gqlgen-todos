package model

type Todo struct {
	ID     int32  `json:"id"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	UserID int32  `json:"user"`
}
