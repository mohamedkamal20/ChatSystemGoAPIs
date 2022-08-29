package redis

import (
	"chatSystemGoAPIs/repositories/chat"
	"chatSystemGoAPIs/repositories/message"
	"context"
	"database/sql"
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"os"
	"strconv"
)


func getClient(ctx context.Context,db int) *goredislib.Client{
	client := goredislib.NewClient(&goredislib.Options{
		Addr: os.Getenv("redisHost"),
		Password: "",
		DB: db,
	})
	pong, err := client.Ping(ctx).Result()
	if err != nil{
		fmt.Println(pong, err)
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


func GetChatNumber(token string, appId int) int{
	chatNumber := -1
	key := token
	db ,_ := strconv.Atoi(os.Getenv("redisDB"))
	ctx := context.Background()
	client := getClient(ctx, db)
	keyFound := false

	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	mutex := rs.NewMutex("redSync-chat")

	if err := mutex.LockContext(ctx); err != nil {
		panic(err)
	}

	if ! IsKeyFound(client, ctx, key) {
		conn := mysqlConnection()
		chatRepo := chat.NewSQLChatRepo(conn)
		count := chatRepo.GetApplicationChatsCount(appId)
		if count != -1{
			err := client.Set(ctx, key, count, 0).Err()
			if err == nil{
				keyFound = true
			}
		}else{
			keyFound = false
		}

	}else{
		keyFound = true
	}

	if _, err := mutex.UnlockContext(ctx); err != nil {
		panic(err)
	}

	if keyFound{
		value, err := client.Incr(ctx, key).Result()
		if err != nil {
			return chatNumber
		}
		chatNumber = int(value)
		return chatNumber
	}else{
		return chatNumber
	}
}

func GetMessageNumber(token string, chatNumber int, chatId int) int{

	messageNumber := -1
	key := token + "_" + strconv.Itoa(chatNumber)
	db ,_ := strconv.Atoi(os.Getenv("redisDB"))
	ctx := context.Background()
	client := getClient(ctx, db)
	keyFound := false

	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	mutex := rs.NewMutex("redSync-message")

	if err := mutex.LockContext(ctx); err != nil {
		panic(err)
	}

	if ! IsKeyFound(client, ctx, key) {
		conn := mysqlConnection()
		messageRepo := message.NewSQLMessageRepo(conn)
		count := messageRepo.GetChatMessagesCount(chatId)
		if count != -1{
			err := client.Set(ctx, key, count, 0).Err()
			if err == nil{
				keyFound = true
			}
		}else{
			keyFound = false
		}

	}else{
		keyFound = true
	}

	if _, err := mutex.UnlockContext(ctx); err != nil {
		panic(err)
	}

	if keyFound{
		value, err := client.Incr(ctx, key).Result()
		if err != nil {
			return messageNumber
		}
		messageNumber = int(value)
		return messageNumber
	}else{
		return messageNumber
	}
}

func IsKeyFound(client *goredislib.Client, ctx context.Context , key string) bool {
	value, err := client.Get(ctx, key).Result()
	if err != nil || value == ""{
		return false
	}
	return true
}


