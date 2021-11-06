package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type server struct{ pb.UnimplementedGreeterServer }

var dir = "pozo.txt"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func StringManage(msg string, t int) string {
	file, _ := os.OpenFile(dir, os.O_RDONLY, 0644)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	if t == 1 {
		if len(text) == 0 {
			return "0"
		}
		maxValue := strings.Split(text[len(text)-1], " ")[2]
		return maxValue
	}

	var aux = strings.Split(msg, " ")
	if len(text) == 0 {
		return msg
	}
	aux1, _ := strconv.Atoi(strings.Split(text[len(text)-1], " ")[2])
	aux2, _ := strconv.Atoi(aux[2])
	aux[2] = strconv.FormatInt(int64(aux1+aux2), 10)
	return strings.Join(aux[:], " ")
}

func rabbitmqCh() {

	conn, err := amqp.Dial("amqp://guest:guest@dist214.inf.santiago.usm.cl:50051/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		var file, err = os.OpenFile(dir, os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			file, err = os.Create(dir)
			if err != nil {
				log.Fatalf("%s", err)
			}
		}
		defer file.Close()
		var msgF = StringManage(string(d.Body), 0)
		_, err = file.WriteString(msgF + "\n")
		if err != nil {
			log.Fatalf("%s", err)
		}
		err = file.Sync()
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Println("Archivo actualizado existosamente.")
	}
}

func grpcCh() {
	lis, err := net.Listen("tcp", "dist215.inf.santiago.usm.cl:50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if string(in.GetName()) == "val" {
		return &pb.HelloReply{Message: StringManage("", 1)}, nil
	}
	return nil, nil
}

func main() {
	forever := make(chan bool)

	go rabbitmqCh()
	go grpcCh()

	//log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
