package main



import (
	"testing"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

func sendMessage(msg string) {

	var port = "32788" // 5672
	var host = "192.168.99.100" //"localhost"
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		CMD_TOPIC, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare a queue")


	//body := "command: docker run -tid monkey/GoMnM_1.0"
	body := msg
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

func TestRunner_Splitit(t *testing.T) {
	var cmd = "command: docker run -t -i -d busybox:latest sh"
	var parts = strings.SplitAfterN(cmd, " ", 3)
	var args = parts[2]
	log.Println("[", args, "]")
}
func TestRunner_RunsContainerCalledRunMe(t *testing.T) {
	sendMessage("command: docker run --name runner-me -tid busybox:latest sh")
}
