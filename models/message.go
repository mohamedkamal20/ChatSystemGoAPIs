package models

type Message struct {
	Number     int    `json:"number,omitempty" bson:"number,omitempty"`
	Chat_id    int    `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	Message    string `json:"message,omitempty" bson:"message,omitempty"`
	Created_at string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated_at string `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

