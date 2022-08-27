package models

type Chat struct {
	Number         int    `json:"number,omitempty" bson:"number,omitempty"`
	ApplicationId int    `json:"application_id,omitempty" bson:"application_id,omitempty"`
	ChatName      string `json:"chat_name,omitempty" bson:"chat_name,omitempty"`
	CreatedAt     string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
