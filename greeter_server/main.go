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

var TotalPlayer int = 1
var N_play int = 0
var N_Playerr int = 0
var L_player = [3]string{"l", "l", "l"}

var numberG1 int = 0

var RPlayerEliminated string = "-"
var T_T1 int = 0
var TotalG2T2 int = 0
var choiceG2 int = 0
var TeamPlayers = [3]string{"-", "-", "-"}
var LoseTeam string = "0"

var numberG3 int = 0
var PairPlayers = [3]int{-1, -1, -1}
var AnswerPlayers = [3]int{-1, -1, -1}

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
	selected_value := "216"

	text := strings.Split(in.GetName(), " ")
	//fmt.Printf(text)
	if text[0] == "GetNumberRebelds" {
		return &pb.HelloReply{Message: "Ligerito te entregamos respsuesta"}, nil
	}
	fmt.Printf("\n" + text[0])

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

func main() {

	go ListenMessage()

	forever := make(chan bool)
	var choice string

	fmt.Println("Esperando solicitudes")
	DurationOfTime := time.Duration(3) * time.Second

		// Defining function parameter of
		// AfterFunc() method
		f := func() {


				fmt.Println("Function called by "+
						"AfterFunc() after 3 seconssds")
		}

		// Calling AfterFunc() method with its
		// parameter
		Timer1 := time.AfterFunc(DurationOfTime, f)

		defer Timer1.Stop()
		// Calling sleep method
		time.Sleep(10 * time.Second)

	for {
		//Menu()
		fmt.Scanf("%s", &choice)

		if choice == "0" {
			fmt.Println("Valor Actual del pozo: ")
		}

		<-forever
	}
}
