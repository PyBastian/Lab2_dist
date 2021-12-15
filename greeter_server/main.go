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
	port       = "dist214.inf.santiago.usm.cl:50051"
	addrs_pozo = "dist215.inf.santiago.usm.cl:50053"
	addrs_node = ":50054"
)

type server struct{ pb.UnimplementedGreeterServer }

var TotalPlayer int = 1
var N_play int = 0
var N_Playerr int = 0
var L_player = [3]string{"l", "l", "l"}

var dir_maquinas = []string{"213", "215", "216"}

var comandos_enviados []string
var dir_Inf1 = ""
var dir_Inf2 = ""
var forever = make(chan bool)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func grpcChannel(ipAdress string, message string) string {
	fmt.Println("Nos conectamos al Cliente")
	fmt.Println(ipAdress)
	conn, err := grpc.Dial(ipAdress, grpc.WithInsecure(), grpc.WithBlock())
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

	fmt.Printf(in.GetName())
	fmt.Printf("\n")

	text := strings.Split(in.GetName(), " ")
	comandos_enviados = append(comandos_enviados, in.GetName())
	fmt.Printf("%v", comandos_enviados)
	if text[len(text)-1] == "1" {
		if dir_Inf1 == "" {
			randomIndex := rand.Intn(len(dir_maquinas))
			dir_Inf1 = dir_maquinas[randomIndex]
			dir_maquinas = RemoveIndex(dir_maquinas, randomIndex)
			fmt.Print("Maquina 1")
			fmt.Print(dir_Inf1)
			return &pb.HelloReply{Message: dir_Inf1}, nil
		}
		fmt.Print(dir_Inf1)
		return &pb.HelloReply{Message: dir_Inf1}, nil
	}
	if text[len(text)-1] == "2" {
		if dir_Inf1 == "" {
			randomIndex := rand.Intn(len(dir_maquinas))
			dir_Inf2 = dir_maquinas[randomIndex]
			dir_maquinas = RemoveIndex(dir_maquinas, randomIndex)
			fmt.Print("Maquina 2")
			return &pb.HelloReply{Message: dir_Inf2}, nil
		}
		return &pb.HelloReply{Message: dir_Inf2}, nil
	}
	//fmt.Printf(text)
	if text[0] == "GetNumberRebelds" {
		var respuesta = "No se encontro";
		for _, s := range comandos_enviados {
				nuevo := strings.Split(s, " ")
		    if nuevo[1] == text[1] && nuevo[2] == text[2]{
					respuesta = nuevo[2]
				}
		}
		return &pb.HelloReply{Message: respuesta}, nil
	}
	if text[0] == "update" {
		output := strings.Join(comandos_enviados, " ")
		return &pb.HelloReply{Message: output}, nil
	}
	return &pb.HelloReply{Message: "216"}, nil
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

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func helloworld(t time.Time) {

	fmt.Printf("Procederemos al merge\n")

	dir_maquinas = []string{"213", "215", "216"}
	dir_Inf1 = ""
	dir_Inf2 = ""

}

func main() {

	forever := make(chan bool)
	go ListenMessage()

	//var choice string
	fmt.Println("Esperando solicitudes")
	// DurationOfTime := time.Duration(3) * time.Second
	// f := func() {
	// 		fmt.Println("Function called by "+
	// 				"AfterFunc() after 3 seconssds")
	// }
	for {
		//doEvery(120*time.Second, helloworld)
		//defer Timer1.Stop()
		// Calling sleep method

		//Menu()
		//fmt.Scanf("%s", &choice)

		<-forever
	}
}
