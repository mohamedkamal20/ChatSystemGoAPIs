package rabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)


//func jsonStruct(message models.Message) string {
//	// Marshal the struct to JSON and check for errors
//	b, err := json.Marshal(message)
//	if err != nil {
//		return ""
//	}
//	return string(b)
//}

func SendMessage(message string)  bool{
	conn, err := amqp.Dial(os.Getenv("rabbitMQHost"))
	if err != nil{
		return false
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil{
		return false
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		os.Getenv("rabbitMQQueue"), // name
		false,  				    // durable
		false,   				// delete when unused
		false,   				// exclusive
		false,   				// no-wait
		nil,     					// arguments
	)
	if err != nil{
		return false
	}

	body := message
	err = ch.Publish(
		"",     // exchange
		q.Name,          // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil{
		return false
	}
	return true
}

func ReceiveMessage() {
	conn, _ := amqp.Dial(os.Getenv("rabbitMQHost"))
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	q, _ := ch.QueueDeclare(
		os.Getenv("rabbitMQQueue"), // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s \n", d.Body)
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
