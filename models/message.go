package models

type Message struct {
	Number     int    `json:"number,omitempty" bson:"number,omitempty"`
	ChatId    int    `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	Message    string `json:"message,omitempty" bson:"message,omitempty"`
	ChatNumber    int `json:"chat_number,omitempty" bson:"chat_number,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
	CreatedAt string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (t Message) MessageResponse() map[string]interface{} {
	return map[string]interface{}{
		"message":      t.Message,
		"number":    t.Number,
	}
}

func (t Message) MessageErrorResponse(error string) map[string]interface{} {
	return map[string]interface{}{
		"errorMessage":  error,
	}
}