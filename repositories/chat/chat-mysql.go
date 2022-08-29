package chat

import (
	"chatSystemGoAPIs/models"
	cRepo "chatSystemGoAPIs/repositories"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

func (m mysqlChatRepo) GetApplicationChatsCount(appId int) int{
	db := m.Conn
	defer db.Close()

	count := 0
	get,err := db.Query("select COUNT(*) as count from chats as c where c.application_id = ?",appId)
	if err != nil{
		fmt.Println(err)
		count = -1
		return count
	}

	for get.Next(){
		var c int
		err := get.Scan(&c)
		if err != nil{
			fmt.Println(err)
			count = -1
			return count
		}
		count = c
	}

	defer get.Close()

	return count
}





