package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	//"os"
	//"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port       = ":50051"
	addrs_pozo = "dist215.inf.santiago.usm.cl:50053"
	addrs_node = ":50054"
)

type server struct{ pb.UnimplementedGreeterServer }

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func grpcChannel(ipAdress string, message string) string {
	fmt.Println("Nos conectamos al Cliente")
	fmt.Println(ipAdress)
	conn, err := grpc.Dial(ipAdress, grpc.WithInsecure())
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

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("Recibimos Comando \n")
	//214 no porque es la direcciones del Broker
	dir_maquinas := []string{"213", "215", "216"}

	fmt.Printf(in.GetName())

	randomIndex := rand.Intn(3)
	selected_value := dir_maquinas[randomIndex]

	text := strings.Split(in.GetName(), " ")
	//fmt.Printf(text)
	if text[0] == "GetNumberRebelds" {
		return &pb.HelloReply{Message: "Ligerito te entregamos respsuesta"}, nil
	}
	fmt.Printf(text[0])

	return &pb.HelloReply{Message: selected_value}, nil
}

func ListenMessage() {
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

func MenuBroker() {
	fmt.Println("1.- AddCity {N_planeta} {N_ciudad} [nuevo valor]")
	fmt.Println("2.- UpdateName {N_planeta} {N_ciudad} [nuevo valor]")
	fmt.Println("3.- UpdateNumber {N_planeta} {N_ciudad} [nuevo valor]")
	fmt.Println("4.- DeleteCity {N_planeta} {N_ciudad}")
	fmt.Println("")
}

func main() {

	go ListenMessage()

	forever := make(chan bool)
	var choice string

	fmt.Println("Esperando solicitudes")

	for {
		//Menu() del broker
		MenuBroker()

		fmt.Scanf("%s", &choice)
		<-forever
	}
}
