package message

import (
	"chatSystemGoAPIs/models"
	mRepo "chatSystemGoAPIs/repository"
	"github.com/elastic/go-elasticsearch/v7"
)


func NewElasticMessageRepo(Conn *elasticsearch.Client) mRepo.MessageRepo {
	return &ElasticMessageRepo{
		Client: Conn,
	}
}

type ElasticMessageRepo struct {
	Client *elasticsearch.Client
}

func (m ElasticMessageRepo) Insert(message models.Message) bool {
	////client := m.Client
	////collection := client.Database("Test").Collection("Student")
	////insertResult,err := collection.InsertOne(context.TODO(),message)
	//if err != nil {
	//	fmt.Println(err)
	//	return false
	//}
	//if insertResult == nil {
	//	return false
	//}
	return true
}

func (m ElasticMessageRepo) ValidChat(token string, chatNumber int) int {
	chatId := -1
	return chatId
}

