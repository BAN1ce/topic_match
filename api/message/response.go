package message

import "time"

type GetResponseData struct {
	Message []Message `json:"message,omitempty"`
	Total   int       `json:"total,omitempty"`
}

type Message struct {
	ID          string    `json:"id"`
	CreatedTime time.Time `json:"created_time"`
	ExpiredTime time.Time `json:"expired_time"`
	SenderID    string    `json:"sender_id"`
}
