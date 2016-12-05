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
     - docker run -i -t <any-name> bash
     - pipe into docker using
     - docker [start|stop] <any-name>
     - docker rm <any-name>
     How do I pipe to container stdin after docker run?
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
	//docker run -tid busybox sh
	//feb49e8e62a7e1c8b81b83be618963852403f220294c5b398bb127f49e29e344
	//root|~/development/gocode/src/github.com/docker/docker > echo hi | docker exec -i feb49e8e62a7e1c8b81b83be618963852403f220294c5b398bb127f49e29e344 cat
	//hi
	//root|~/development/gocode/src/github.com/docker/docker > echo hi | docker exec -i feb49e8e62a7e1c8b81b83be618963852403f220294c5b398bb127f49e29e344 wc -l
	//1
	//

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var cmd = string(d.Body)
			if (strings.HasPrefix(cmd, "command: docker run")){

				var parts = strings.SplitAfterN(cmd, " ", 3)
				var args = strings.Fields(parts[2])
				//var parts = strings.FieldsFunc(cmd, func)

				log.Println("Docker CCC Run: PP:", args)

				//@FOR /f "tokens=*" %i IN ('docker-machine env default') DO @%i
				//cmd := exec.Command("docker", "run", "-tid", "busybox", "sh")
				// command: docker run -tid busybox sh

				cmd := exec.Command("docker", args...)//"run", "-tid", parts[4], parts[5])
				//cmd.Stdin = strings.NewReader("some input")
				//inPipe, _ := cmd.StdinPipe()
				//outPipe,err := cmd.StdoutPipe()
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr
				err = cmd.Run()

				if err != nil {
					log.Println("Boom - bad", stderr.String())
					log.Fatal(err)
				}
				//cmdReader := bufio.NewReader(outPipe)

				//for {
				log.Println(stdout.String())

					//time.Sleep(1)
					//var line, err2 = cmdReader.ReadString('\n')
					//if err2 == io.EOF {
					//	//err := fmt.Errorf("EOF received: %q", line)
					//	panic(err)
					//}
					////
					//// read stdout - need to get the docker image instance
					//log.Printf("STDOUT %s", line)
				//}
				// os. run docker container
			}
			if (strings.HasPrefix(cmd, "command: docker exec -i")){
				// os. run docker container
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func main() {
	runnerListen()


}