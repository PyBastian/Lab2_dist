package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type server struct{ pb.UnimplementedGreeterServer }

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func WriteInTXT(dir string, msg string) {
	var file, err = os.OpenFile(dir, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		file, err = os.Create(dir)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
	defer file.Close()
	_, err = file.WriteString(msg + "\n")
	if err != nil {
		log.Fatalf("%s", err)
	}
	err = file.Sync()
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Println("Archivo actualizado existosamente.")
}

func grpcCh() {
	lis, err := net.Listen("tcp", ":50081")
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
	text := strings.Split(in.GetName(), " ")
	var dir string = text[1] + "__" + text[0] + ".txt"
	WriteInTXT(dir, text[2])
	return &pb.HelloReply{Message: "recibido"}, nil
}

func main() {
	forever := make(chan bool)

	go grpcCh()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
