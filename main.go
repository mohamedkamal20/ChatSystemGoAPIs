package main

import (
	api "chatSystemGoAPIs/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/applications/{token}/chats",api.CreateChat).Methods("POST")
	r.HandleFunc("/api/v1/applications/{token}/chats/{chat_number}/messages", api.CreateMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8085", r))
}