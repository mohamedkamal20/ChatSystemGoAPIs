package services

import (
	"github.com/go-redis/redis"
)

func getClient(db int) *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: db,
	})
	return client
}

//
//func UpdateCachedStudent(db int , student models.Student)bool{
//	client := getClient(db)
//
//	value, err := client.Get(fmt.Sprintf("%d",student.Id)).Result()
//	if err != nil {
//		fmt.Println(err)
//		return false
//	}
//	var cachedStudent models.Student
//	err = json.Unmarshal([]byte(value) , &cachedStudent)
//
//	cachedStudent.Name = student.Name
//	cachedStudent.Age = student.Age
//	cachedStudent.Sex = student.Sex
//
//	json ,_ := json.Marshal(cachedStudent)
//	err = client.Set(fmt.Sprintf("%d",cachedStudent.Id),json,time.Minute * 15).Err()
//	if err != nil {
//		fmt.Println(err)
//		return false
//	}
//	return true
//}
//
//func SaveCachedData(db int, students []models.Student){
//	client := getClient(db)
//	err := client.Set( "count",len(students),time.Minute * 15).Err()
//	if err != nil {
//		fmt.Println(err)
//	}
//	for _,student:= range students{
//		json ,_ := json.Marshal(student)
//		err := client.Set(fmt.Sprintf("%d",student.Id),json,time.Minute * 15).Err()
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//}
