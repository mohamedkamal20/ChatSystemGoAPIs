package chat

import (
	"chatSystemGoAPIs/models"
	cRepo "chatSystemGoAPIs/repositories"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)


func NewSQLChatRepo(Conn *sql.DB) cRepo.ChatRepo {
	return &mysqlChatRepo{
		Conn: Conn,
	}
}

type mysqlChatRepo struct {
	Conn *sql.DB
}

func (m mysqlChatRepo) Insert(chat models.Chat) bool {
	db := m.Conn
	defer db.Close()


	today := time.Now()
	chat.CreatedAt = today.Format("2006-01-02 15:04:05")
	chat.UpdatedAt = today.Format("2006-01-02 15:04:05")

	insert,err := db.Prepare("insert into chats(number, application_id, chat_name, created_at, updated_at) values(?,?,?,?,?)")
	if err != nil{
		fmt.Println(err)
		return false
	}

	insert.Exec(chat.Number,chat.ApplicationId,chat.ChatName, chat.CreatedAt, chat.UpdatedAt)

	defer insert.Close()

	return true
}

func (m mysqlChatRepo) ValidApplication(token string) int {
	db := m.Conn
	defer db.Close()

	appId := -1
	get,err := db.Query("select id from applications where token = ?",token)
	if err != nil{
		fmt.Println(err)
		return appId
	}

	for get.Next(){
		var id int
		err := get.Scan(&id)
		if err != nil{
			fmt.Println(err)
			return appId
		}
		appId = id
	}

	defer get.Close()

	return appId
}




