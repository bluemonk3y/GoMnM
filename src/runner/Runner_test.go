package main



import (
	"testing"
	"github.com/streadway/amqp"
	"log"
)

func TestRunner_SendsMsg(t *testing.T) {
	var port = "32788" // 5672
	var host = "192.168.99.100" //"localhost"
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"cmd", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "command: docker run monkey/GoMnM_1.0"
	err = ch.Publish(
		"", // exchange
		q.Name, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")

}
