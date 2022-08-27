package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
)

func getClient(db int) *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: db,
	})
	pong, err := client.Ping().Result()
	if err != nil{
		fmt.Println(pong, err)
		return nil
	}
	return client
}


func GetChatNumber(token string) int{
	db ,_ := strconv.Atoi(os.Getenv("redisDB"))
	client := getClient(db)
	chatNumber := -1
	value, err := client.Incr(token).Result()
	if err != nil {
		fmt.Println(err)
		return chatNumber
	}
	chatNumber = int(value)
	return chatNumber
}

func GetMessageNumber(token string, chatNumber int) int{
	db ,_ := strconv.Atoi(os.Getenv("redisDB"))
	client := getClient(db)
	messageNumber := -1
	key := token + "_" + strconv.Itoa(chatNumber)
	value, err := client.Incr(key).Result()
	if err != nil {
		fmt.Println(err)
		return messageNumber
	}
	messageNumber = int(value)
	return messageNumber
}


