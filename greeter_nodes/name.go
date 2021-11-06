package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type data struct{ pb.UnimplementedGreeterServer }

const (
	port = ":50054"
	dir  = "namenode.txt"
)

var ListOfDataNode = [3]string{"localhost:50081", "localhost:50082", "localhost:50083"}
var aux int = 0

func grpcChannel(message string, ipAddress string) string {
	conn, err := grpc.Dial(ipAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return r.GetMessage()
}

func TXT(msg string) {
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
	fmt.Println("Se actualizo el archivo de texto.")
}

func ListenInstr() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func Exist(IDPlayer string, game string) int {
	file, _ := os.OpenFile(dir, os.O_RDONLY, 0644)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	for i := 0; i < len(text); i++ {
		aux := strings.Split(text[i], " ")
		if aux[0] == IDPlayer && aux[1] == game {
			return 1
		}
	}
	return 0
}

func SearchIP(IDPlayer string, game string) string {
	file, _ := os.OpenFile(dir, os.O_RDONLY, 0644)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	for i := 0; i < len(text); i++ {
		aux := strings.Split(text[i], " ")
		if aux[0] == IDPlayer && aux[1] == game {
			return aux[2]
		}
	}
	return ""
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	text := strings.Split(in.GetName(), " ")
	if Exist(text[1], text[0]) != 1 {
		WriteInTXT(text[1]+" "+text[0]+" "+ListOfDataNode[aux], "Esta es una prueba del namenode")
	}
	aux = (aux)%len(ListOfDataNode) + 1
	_ = grpcChannel(in.GetName(), SearchIP(text[1], text[0]))
	return &pb.HelloReply{Message: "recibido"}, nil
}

func main() {
	forever := make(chan bool)

	go ListenInstr()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}
