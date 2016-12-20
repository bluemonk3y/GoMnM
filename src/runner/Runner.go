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
	"io"
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

				var parts = strings.SplitAfterN(cmd, " ", 3)
				var args = strings.Fields(parts[2])

				log.Println("Docker Run: PP:", args)

				cmd := exec.Command("docker", args...)
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr
				//err = cmd.Run()
				err = cmd.Start()

				if err != nil {
					log.Println("Boom - bad", stderr.String())
					log.Fatal(err)
					// respond on error channel
				}


				log.Println(stderr.String())
				log.Println(stdout.String())
				var svc_name = parts[2]
				log.Println("Adding to map:" + svc_name)
				dockerProcessMap[svc_name] = cmd

				var stdin,_ = cmd.StdinPipe()
				//stdin.
				io.WriteString(stdin, " 1111\n")
				//stdin.

			}
			// "command: docker svc-input --name svc-key-id"
			if (strings.HasPrefix(cmd, "command: svc-input --name")){

				// docker input -i svc-key-name data
				var parts = strings.Split(cmd, " ")
				var svcKeyId = parts[3]
				var cmd = dockerProcessMap[svcKeyId]
				log.Println("Going to send to stdin key,", svcKeyId, " processmap-size", len(dockerProcessMap))
				// write to std in
				var stdin,_ = cmd.StdinPipe()
//				stdin.Write("Helloooo! cmd 11111:\n")
				log.Println(stdin)
				io.WriteString(stdin, " 22222\n")


				//args.Stdou
			}
			if (strings.HasPrefix(cmd, "command: kill --name svc-key-id")){
				var parts = strings.Split(cmd, " ")
				var svcKeyId = parts[3]
				var cmd = dockerProcessMap[svcKeyId]
				log.Println("Killing", cmd)
				var err = cmd.Process.Kill()
				if (err != nil) {
					log.Println("Failed to kill:", cmd)
				}

			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func main() {
	runnerListen()


}