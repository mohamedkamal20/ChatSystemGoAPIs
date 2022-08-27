package controllers

import (
	"chatSystemGoAPIs/models"
	"chatSystemGoAPIs/repositories/chat"
	"chatSystemGoAPIs/repositories/message"
	"chatSystemGoAPIs/services/rabbitMQ"
	"chatSystemGoAPIs/services/redis"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"time"
)

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
		responseJson,_ := json.Marshal(tmpChat.ChatErrorResponse("chat_name can not be empty"))
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
			chatNumber := redis.GetChatNumber(token)
			if chatNumber != -1{
				tmpChat.Number = chatNumber
				rabbitMQJson,_ := json.Marshal(tmpChat)
				if rabbitMQ.SendMessage(string(rabbitMQJson), os.Getenv("rabbitMQChatsQueue")){
					w.WriteHeader(http.StatusOK)
					responseJson,_ := json.Marshal(tmpChat.ChatResponse())
					w.Write(responseJson)
				}else{
					w.WriteHeader(http.StatusBadRequest)
					responseJson,_ := json.Marshal(tmpChat.ChatErrorResponse("Something went wrong please try again later"))
					w.Write(responseJson)
				}
			}else{
				w.WriteHeader(http.StatusBadRequest)
				responseJson,_ := json.Marshal(tmpChat.ChatErrorResponse("Something went wrong please try again later"))
				w.Write(responseJson)
			}

		}else {
			w.WriteHeader(http.StatusBadRequest)
			responseJson,_ := json.Marshal(tmpChat.ChatErrorResponse("Application token does not exist"))
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
		responseJson,_ := json.Marshal(tmpMessage.MessageErrorResponse("Message can not be empty"))
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
			messageNumber := redis.GetMessageNumber(token, chatNumber)
			if messageNumber != -1{
				tmpMessage.Number = messageNumber

				rabbitMQJson,_ := json.Marshal(tmpMessage)
				if rabbitMQ.SendMessage(string(rabbitMQJson), os.Getenv("rabbitMQMessagesQueue")){
					w.WriteHeader(http.StatusOK)
					responseJson,_ := json.Marshal(tmpMessage.MessageResponse())
					w.Write(responseJson)
				}else{
					w.WriteHeader(http.StatusBadRequest)
					responseJson,_ := json.Marshal(tmpMessage.MessageErrorResponse("Something went wrong please try again later"))
					w.Write(responseJson)
				}
			}else{
				w.WriteHeader(http.StatusBadRequest)
				responseJson,_ := json.Marshal(tmpMessage.MessageErrorResponse("Something went wrong please try again later"))
				w.Write(responseJson)
			}
		}else{
			w.WriteHeader(http.StatusBadRequest)
			responseJson,_ := json.Marshal(tmpMessage.MessageErrorResponse("Application token or chat number does not exist"))
			w.Write(responseJson)
		}
	}
}
