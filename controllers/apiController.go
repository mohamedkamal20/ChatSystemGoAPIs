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
	"time"
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

	var tmpChat models.Chat

	json.NewDecoder(r.Body).Decode(&tmpChat)

	if tmpChat.ChatName == "" {
		w.WriteHeader(http.StatusBadRequest)
		responseJson,_ := json.Marshal(tmpChat)
		w.Write(responseJson)
	}else {
		conn := mysqlConnection()
		chatRepo := chat.NewSQLChatRepo(conn)

		appId := chatRepo.ValidApplication(token)

		if appId != -1 {
			tmpChat.ApplicationId = appId
			today := time.Now()
			tmpChat.CreatedAt = today.Format("2006-01-02 15:04:05")
			tmpChat.UpdatedAt = today.Format("2006-01-02 15:04:05")
			tmpChat.Number = 6 // to be concurent calculate

			rabbitMQJson,_ := json.Marshal(tmpChat)
			if rabbitMQ.SendMessage(string(rabbitMQJson), os.Getenv("rabbitMQChatsQueue")){
				w.WriteHeader(http.StatusOK)
				responseJson,_ := json.Marshal(tmpChat)
				w.Write(responseJson)
			}else{

			}
			//inserted := chatRepo.Insert(tmpChat)
			//if inserted {
			//	w.WriteHeader(http.StatusOK)
			//	json,_ := json.Marshal(tmpChat)
			//	w.Write(json)
			//}else {
			//	w.WriteHeader(http.StatusBadRequest)
			//	json,_ := json.Marshal(tmpChat)
			//	w.Write(json)
			//}
		}else {
			w.WriteHeader(http.StatusBadRequest)
			responseJson,_ := json.Marshal(tmpChat)
			w.Write(responseJson)
		}
	}
}

func CreateMessage(w http.ResponseWriter ,r * http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	token := vars["token"]

	chatNumberStr := vars["chat_number"]
	chatNumber := 0
	chatNumber, _ = strconv.Atoi(chatNumberStr)

	var tmpMessage models.Message

	json.NewDecoder(r.Body).Decode(&tmpMessage)

	tmpMessage.ChatNumber = chatNumber
	tmpMessage.Token = token

	if tmpMessage.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		responseJson,_ := json.Marshal(tmpMessage)
		w.Write(responseJson)
	}else{
		mysqlConn := mysqlConnection()
		mysqlMessageRepo := message.NewSQLMessageRepo(mysqlConn)

		chatId := mysqlMessageRepo.ValidChat(token, chatNumber)

		if chatId != -1 {
			tmpMessage.ChatId = chatId
			today := time.Now()
			tmpMessage.CreatedAt = today.Format("2006-01-02 15:04:05")
			tmpMessage.UpdatedAt = today.Format("2006-01-02 15:04:05")
			tmpMessage.Number = 5 // to be concurent calculate

			rabbitMQJson,_ := json.Marshal(tmpMessage)
			if rabbitMQ.SendMessage(string(rabbitMQJson), os.Getenv("rabbitMQMessagesQueue")){
				w.WriteHeader(http.StatusOK)
				responseJson,_ := json.Marshal(tmpMessage)
				w.Write(responseJson)
			}else{
				w.WriteHeader(http.StatusBadRequest)
				responseJson,_ := json.Marshal(tmpMessage)
				w.Write(responseJson)
			}

			//inserted := mysqlMessageRepo.Insert(tmpMessage)
			//if inserted {
			//	elasticClient := elasticClient()
			//	elasticMessageRepo := message.NewElasticMessageRepo(elasticClient)
			//	elasticMessageRepo.Insert(tmpMessage)
			//	w.WriteHeader(http.StatusOK)
			//	json,_ := json.Marshal(tmpMessage)
			//	w.Write(json)
			//}else {
			//	w.WriteHeader(http.StatusBadRequest)
			//	json,_ := json.Marshal(tmpMessage)
			//	w.Write(json)
			//}
		}else{
			w.WriteHeader(http.StatusBadRequest)
			responseJson,_ := json.Marshal(tmpMessage)
			w.Write(responseJson)
		}
	}
}
