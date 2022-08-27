package controllers

import (
	"chatSystemGoAPIs/models"
	"chatSystemGoAPIs/repositories/chat"
	"chatSystemGoAPIs/repositories/message"
	"chatSystemGoAPIs/services/rabbitMQ"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

func elasticClient()(client *elasticsearch.Client){
	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("elasticSearchHost"),
		},
		Username: "",
		Password: "",
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		return nil
	}

	_, err = client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil
	}
	return client
}

func mysqlConnection()(db * sql.DB){

	db,err := sql.Open("mysql",os.Getenv("mySqlDataStoreName"))

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}


func CreateChat(w http.ResponseWriter ,r * http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	token := vars["token"]


	conn := mysqlConnection()
	chatRepo := chat.NewSQLChatRepo(conn)

	var tmpChat models.Chat

	json.NewDecoder(r.Body).Decode(&tmpChat)
	if tmpChat.ChatName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json,_ := json.Marshal(tmpChat)
		w.Write(json)
	}

	appId := chatRepo.ValidApplication(token)

	if appId != -1 {
		tmpChat.ApplicationId = appId
		tmpChat.Number = 6 // to be concurent calculate
		inserted := chatRepo.Insert(tmpChat)
		if inserted {
			w.WriteHeader(http.StatusOK)
			json,_ := json.Marshal(tmpChat)
			w.Write(json)
		}else {
			w.WriteHeader(http.StatusBadRequest)
			json,_ := json.Marshal(tmpChat)
			w.Write(json)
		}
	}else {
		w.WriteHeader(http.StatusBadRequest)
		json,_ := json.Marshal(tmpChat)
		w.Write(json)
	}


}

func CreateMessage(w http.ResponseWriter ,r * http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	token := vars["token"]

	chatNumberStr := vars["chat_number"]
	chatNumber := 0
	chatNumber, _ = strconv.Atoi(chatNumberStr)

	mysqlConn := mysqlConnection()
	mysqlMessageRepo := message.NewSQLMessageRepo(mysqlConn)

	var tmpMessage models.Message

	json.NewDecoder(r.Body).Decode(&tmpMessage)

	tmpMessage.ChatNumber = chatNumber
	tmpMessage.Token = token
	tmpMessage.Number = 5 // to be concurent calculate

	rabbitMQ.SendMessage("hey")

	elasticClient := elasticClient()
	elasticMessageRepo := message.NewElasticMessageRepo(elasticClient)

	elasticMessageRepo.Insert(tmpMessage)


	if tmpMessage.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		json,_ := json.Marshal(tmpMessage)
		w.Write(json)
	}

	chatId := mysqlMessageRepo.ValidChat(token, chatNumber)

	if chatId != -1 {
		tmpMessage.ChatId = chatId
		tmpMessage.Number = 5 // to be concurent calculate
		inserted := mysqlMessageRepo.Insert(tmpMessage)
		if inserted {
			elasticMessageRepo.Insert(tmpMessage)
			w.WriteHeader(http.StatusOK)
			json,_ := json.Marshal(tmpMessage)
			w.Write(json)
		}else {
			w.WriteHeader(http.StatusBadRequest)
			json,_ := json.Marshal(tmpMessage)
			w.Write(json)
		}
	}else{
		w.WriteHeader(http.StatusBadRequest)
		json,_ := json.Marshal(tmpMessage)
		w.Write(json)
	}


}
