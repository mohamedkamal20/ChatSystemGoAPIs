package main

import (
	api "chatSystemGoAPIs/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func setEnvVariables()  {
	os.Setenv("rabbitMQHost", "amqp://guest:guest@localhost:5672/")
	os.Setenv("rabbitMQChatsQueue", "chats_development")
	os.Setenv("rabbitMQMessagesQueue", "messages_development")
	os.Setenv("mySqlDataStoreName", "root:root@tcp(127.0.0.1:3306)/ChatSystem_development")
	os.Setenv("redisDB", "1")
}

func main() {
	// set env variables
	setEnvVariables()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/applications/{token}/chats",api.CreateChat).Methods("POST")
	r.HandleFunc("/api/v1/applications/{token}/chats/{chat_number}/messages", api.CreateMessage).Methods("POST")

	log.Fatal(http.ListenAndServe(":8085", r))
}