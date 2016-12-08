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

/**
  * Run commands like
     - docker run -id blu3monk3y/simple-ms:v1
     - pipe into docker using
     - docker [start|stop] <any-name>
     - docker rm <any-name>
     How do I pipe to container stdin after docker run?
     docker attach --detach-keys=ctrl-a c4ca4f19d4cd
 */
func runnerListen() {


	var port = "32788" // 5672
	var host = "192.168.99.100" //"localhost"
	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")
	defer conn.Close()

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

	go func() {
		//make(map[int]PidMap)
		var serviceMap = make(map[string]*exec.Cmd)
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var cmd = string(d.Body)
			// docker run --name SVC-ID -tid blu3monk3y/simple-ms:v1
			if (strings.HasPrefix(cmd, "command: docker run")){

				var parts = strings.SplitAfterN(cmd, " ", 3)
				var args = strings.Fields(parts[2])

				log.Println("Docker CCC Run: PP:", args)


				cmd := exec.Command("docker", args...)
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr
				err = cmd.Run()

				if err != nil {
					log.Println("Boom - bad", stderr.String())
					log.Fatal(err)
					// respond on error channel
				}


				log.Println(stderr.String())
				log.Println(stdout.String())
				var svc_name = parts[2]
				log.Println("Adding to map:" + svc_name)
				serviceMap[svc_name] = cmd

			}
			if (strings.HasPrefix(cmd, "command: docker svc-input --name svc-key-id")){
				// docker input -i svc-key-name data
				var parts = strings.SplitAfterN(cmd, " ", 3)
				var args = strings.Fields(parts[2])
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func main() {
	runnerListen()


}