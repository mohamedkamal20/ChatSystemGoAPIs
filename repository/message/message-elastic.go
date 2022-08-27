package message

import (
	"chatSystemGoAPIs/models"
	mRepo "chatSystemGoAPIs/repository"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"strings"
)


type ElasticDocs struct {
	Message string `json:"message,omitempty"`
	Number int `json:"number,omitempty"`
	Token string `json:"token,omitempty"`
	ChatNumber int `json:"chat_number,omitempty"`
}

func jsonStruct(doc models.Message) string {

	// Create struct instance of the Elasticsearch fields struct object
	docStruct := &ElasticDocs{
		Message: doc.Message,
		Number: doc.Number,
		Token: doc.Token,
		ChatNumber: doc.ChatNumber,
	}

	// Marshal the struct to JSON and check for errors
	b, err := json.Marshal(docStruct)
	if err != nil {
		return ""
	}
	return string(b)
}

func NewElasticMessageRepo(Conn *elasticsearch.Client) mRepo.MessageRepo {
	return &ElasticMessageRepo{
		Client: Conn,
	}
}

type ElasticMessageRepo struct {
	Client *elasticsearch.Client
}

func (m ElasticMessageRepo) Insert(message models.Message) bool {
	jsonStruct := jsonStruct(message)
	es := m.Client
	res, err := es.Index(
		"messages_development_20220825215747432",   // Index name
		strings.NewReader(jsonStruct),   // Document body
		es.Index.WithRefresh("true"),  // Refresh
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
		return false
	}
	defer res.Body.Close()
	return true
}

func (m ElasticMessageRepo) ValidChat(token string, chatNumber int) int {
	chatId := -1
	return chatId
}



