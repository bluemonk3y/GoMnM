package main



import (
	"testing"
	"github.com/streadway/amqp"
	"log"
	"strings"
	"time"
	"os/exec"
	//"bytes"
	"io"
	"bytes"
)

func connect() (*amqp.Channel, string) {

	log.Println("Connecting")

	var port = "32769" // 5672
	var host = "localhost"
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")

	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	//defer ch.Close()

	q, err := ch.QueueDeclare(
		CMD_TOPIC, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare a queue")




	return ch, q.Name
}
func sendMessage(ch *amqp.Channel, qname string, msg string) {
	body := msg
	err := ch.Publish(
		"", // exchange
		qname, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func checkErr(err error) {
	if err != nil {
		log.Println("Boom - cannot RUN:", err)
		log.Fatal(err)
	}
}
func TestRunner_RunsServicesRemotely(t *testing.T) {
	ch, qname := connect()
	sendMessage(ch, qname, "command: docker run --name SVC-TEST -id blu3monk3y/simple-ms:v1")
	sendMessage(ch, qname, "command: stdin SVC-TEST q111\n")
	sendMessage(ch, qname, "command: stdin SVC-TEST q222\n")
	sendMessage(ch, qname, "command: stdin SVC-TEST \"qhello 3\"\n")

	sendMessage(ch, qname, "command: remove SVC-TEST")
}

func TestRunner_RunsProcessToReadStdIn(t *testing.T) {
	//sendMessage("command: docker run --name SVC-TEST -ti blu3monk3y/simple-ms:v1")
	log.Println("Docker Running")

	var run1 = "run --name SVC-TEST -i blu3monk3y/simple-ms:v1"
	cmd := exec.Command("docker", strings.Fields(run1)...)

	log.Println("Docker Running 111")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	stdin, err := cmd.StdinPipe()
	checkErr(err)
	defer stdin.Close()

	err = cmd.Start()
	checkErr(err)

	log.Println("Docker Running 222")

	time.Sleep(2 * time.Second)
	log.Println("Started simple-ms")

	num, err := io.WriteString(stdin, " 1111\n")
	checkErr(err)
	io.WriteString(stdin, " 1111\n")
	io.WriteString(stdin, " 2222\n")
	io.WriteString(stdin, " 3333\n")
	io.WriteString(stdin, " 444\n")


	log.Println("Wrote:", num)

	log.Println("StdOut", stdout.String())
	log.Println("StdErr", stderr.String())


	exec.Command("docker", "stop", "SVC-TEST").Run()
	exec.Command("docker", "rm", "SVC-TEST").Run()

}
func TestRunner_RunsDockerThenAttachesToWriteToStdIn(t *testing.T) {
	//sendMessage("command: docker run --name SVC-TEST -ti blu3monk3y/simple-ms:v1")
	//log.Println("Docker Running")
	//
	//var run1 = "run --name SVC-TEST -id blu3monk3y/simple-ms:v1"
	//cmd := exec.Command("docker", strings.Fields(run1)...)
	//
	//err := cmd.Run()
	//checkErr(err)
	//
	//time.Sleep(2 * time.Second)
	log.Println("Started simple-ms")

	cmd2 := exec.Command("docker", "attach", "SVC-TEST")

	//var stdout bytes.Buffer
	//var stderr bytes.Buffer
	//cmd2.Stdout = &stdout
	//cmd2.Stderr = &stderr


	stdin, err := cmd2.StdinPipe()
	checkErr(err)
	defer stdin.Close()

	err = cmd2.Start()
	checkErr(err)

	num, err := io.WriteString(stdin, " 1111\n")
	checkErr(err)
	io.WriteString(stdin, " 1111\n")
	io.WriteString(stdin, " 2222\n")
	io.WriteString(stdin, " 3333\n")
	io.WriteString(stdin, " 444\n")


	log.Println("Wrote:", num)

	//cmd2.Wait()
	//log.Println("StdOut", stdout.String())
	//log.Println("StdErr", stderr.String())
	//
	//time.Sleep(1 * time.Second)
	//
	//exec.Command("docker", "stop", "SVC-TEST").Run()
	//exec.Command("docker", "rm", "SVC-TEST").Run()
}

func TestRunner_Splitit(t *testing.T) {
	var cmd = "command: docker run -id busybox:latest sh"
	var parts = strings.SplitAfterN(cmd, " ", 3)
	var args = parts[2]
	log.Println("[", args, "]")
}
func TestRunner_RunsContainerCalledRunMe(t *testing.T) {
	//sendMessage("command: docker run --name SVC-TEST -tid busybox:latest sh")
	//
	//// send remaining messages
	//time.Sleep(100 * time.Millisecond)
	//sendMessage("command: svc-input --name SVC-TEST BOOOM")


}
