package models

type Chat struct {
	Number         int    `json:"number,omitempty" bson:"number,omitempty"`
	Application_id int    `json:"application_id,omitempty" bson:"application_id,omitempty"`
	Chat_name      string `json:"chat_name,omitempty" bson:"chat_name,omitempty"`
	Created_at     string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated_at     string `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
