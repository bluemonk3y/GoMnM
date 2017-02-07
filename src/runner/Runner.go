package main

import (
	"log"
	_ "os"
	"github.com/streadway/amqp"
	"strings"
	"os/exec"
	_ "bufio"
	_ "fmt"
	_ "io"
	_ "time"
	//"bytes"
	//"io"
	"io"
	//"bytes"
	"bytes"
)

const (
	CMD_TOPIC = "mnm-cmds"
	MyDB2 = "introspector-1"
)


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func checkErr1(err error) {
	if err != nil {
		log.Println("Boom - cannot RUN:", err)
		log.Fatal(err)
	}
}


/**
  * Run commands like
     - docker run -id blu3monk3y/simple-ms:v1
     - pipe into docker using
     - docker [start|stop] <any-name>
     - docker rm <any-name>
     How do I pipe to container stdin after docker run?
     docker attach --detach-keys=ctrl-a c4ca4f19d4cd
 */
func handleRunCmd(cmdString string) {

	var parts = strings.SplitAfterN(cmdString, " ", 3)
	var args = strings.Fields(parts[2])

	log.Println("Docker Run: ARGS:", args)

	cmd := exec.Command("docker", args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Println("Boom - cannot RUN:", stderr.String())
		log.Fatal(err)
	}

	//var svc_name = args[2]
	//dockerProcessMap[svc_name] = cmd

}
func handleStdIn(cmdString string) {
	log.Println("Docker StdIn: ARGS:", cmdString)

	var parts = strings.SplitAfterN(cmdString, " ", 4)
	var svcKeyId = parts[2]

	cmd := exec.Command("docker", "attach", strings.Trim(svcKeyId, " "))

	//log.Println("StdIn SVC:'", strings.Trim(svcKeyId, " "), "' msg:", parts[3])
	var stdin,err = cmd.StdinPipe()
	checkErr1(err)
	defer stdin.Close()
	cmd.Start()

	_, err = io.WriteString(stdin, parts[3])
	// need newline to flush
	_, err = io.WriteString(stdin, "\n")
	checkErr1(err)

}
func runnerListen() {


	var port = "32769" // 5672
	var host = "localhost"
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")
	defer conn.Close()

	failOnError(err, "Failed to open a connect")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		CMD_TOPIC, // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	var dockerProcessMap = make(map[string]*exec.Cmd)

	go func() {

		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var cmd = string(d.Body)
			// docker run --name SVC-ID -tid blu3monk3y/simple-ms:v1
			if (strings.HasPrefix(cmd, "command: docker run")){
				handleRunCmd(cmd)
			}
			// "command: stdin SVC-ID \"some text""
			if (strings.HasPrefix(cmd, "command: stdin")){
				handleStdIn(cmd)
			}
			if (strings.HasPrefix(cmd, "command: remove")){
				var parts = strings.Fields(cmd)
				var svcKeyId = parts[2]
				var cmd = dockerProcessMap[svcKeyId]
				log.Println("Killing", svcKeyId)
				//var err = cmd.Process.Kill()
				if (err != nil) {
					log.Println("Failed to kill:", cmd)
				}
				exec.Command("docker", "stop", svcKeyId).Run()
				exec.Command("docker", "rm", svcKeyId).Run()
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func main() {
	runnerListen()
	//handleRunCmd("command: docker run --name SVC-TEST -id blu3monk3y/simple-ms:v1")
	//handleStdIn("command: stdin SVC-TEST q111\n")
	//handleStdIn("command: stdin SVC-TEST q222\n")
	//handleStdIn("command: stdin SVC-TEST \"qhello 3\"\n")
	//handleStdIn("command: stdin SVC-TEST q111")


}