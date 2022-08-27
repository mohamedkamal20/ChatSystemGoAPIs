package message


import (
"chatSystemGoAPIs/models"
mRepo "chatSystemGoAPIs/repository"
"database/sql"
"fmt"
_ "github.com/go-sql-driver/mysql"
	"time"
)


func NewSQLMessageRepo(Conn *sql.DB) mRepo.MessageRepo {
	return &mysqlMessageRepo{
		Conn: Conn,
	}
}

type mysqlMessageRepo struct {
	Conn *sql.DB
}

func (m mysqlMessageRepo) Insert(message models.Message) bool {
	db := m.Conn
	defer db.Close()

	today := time.Now()
	message.CreatedAt = today.Format("2006-01-02 15:04:05")
	message.UpdatedAt = today.Format("2006-01-02 15:04:05")

	insert,err := db.Prepare("insert into messages(number, chat_id, message, created_at, updated_at) values(?,?,?,?,?)")
	if err != nil{
		fmt.Println(err)
		return false
	}

	insert.Exec(message.Number,message.ChatId,message.Message, message.CreatedAt, message.UpdatedAt)

	defer insert.Close()

	return true
}

func (m mysqlMessageRepo) ValidChat(token string, chatNumber int) int {
	db := m.Conn
	defer db.Close()

	chatId := -1
	get,err := db.Query("select c.id from chats as c join applications as a on c.application_id = a.id where a.token = ? and c.number = ?",token, chatNumber)
	if err != nil{
		fmt.Println(err)
		return chatId
	}

	for get.Next(){
		var id int
		err := get.Scan(&id)
		if err != nil{
			fmt.Println(err)
			return chatId
		}
		chatId = id
	}

	defer get.Close()

	return chatId
}



