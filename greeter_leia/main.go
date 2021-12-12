package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	//"os"
	//"strconv"
	"time"
	//"strings"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)

var (
	name = flag.String("name", defaultName, "Name to greet")

)

type server struct{ pb.UnimplementedGreeterServer }

var G_now string = ""
var id_user string = ""
var ReadyToPlay string = ""
var R_Game = ""
func ListenInstr() {
	lis, err := net.Listen("tcp", "dist215.inf.santiago.usm.cl:50071")
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
	if in.GetName() == "G1" || in.GetName() == "G2" || in.GetName() == "G3" {
		G_now = in.GetName()
	}
	if in.GetName() == "Ready" {
		ReadyToPlay = "Ready"
	}
	return &pb.HelloReply{Message: "recibido"}, nil

}

//Esto envia automaticaamente ingo a 214 (Server)
func grpcChannel(message string) string {
	fmt.Println("")
	conn, err := grpc.Dial("dist214.inf.santiago.usm.cl:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Error de conecc'on con host: %v", err)
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

func C_Lider(msg string, n_planeta string, n_ciudad string) string {
	//fmt.Println("Me voy a comunicar con el Lider")
	var comando string;

	comando = msg + " " + n_planeta + " " + n_ciudad
	fmt.Printf("Comando Final \n")
	fmt.Printf(comando)

	if msg == "GetNumberRebelds" {
		//fmt.Println("Entrando al grpcChanel pa mandarle algo al Lider")
		return grpcChannel(comando)
	}
	r := grpcChannel(msg)
	return r
}

func Menu() {
	fmt.Println("GetNumberRebelds {N_planeta} {N_ciudad}")
}

func main() {

	var choice, N_planeta, N_ciudad string
	var respuesta_host string

	forever := make(chan bool)
	go ListenInstr()

	fmt.Println("Bienvenide Leia Organa, asi seran tus comandos:\n")

	Menu()


  //fmt.Printf("captured: %s %s %s %s\n", choice, N_planeta, N_ciudad, N_valor)


	for {

		fmt.Scanf("%s %s %s", &choice, &N_planeta, &N_ciudad)
		fmt.Println("Hablemos Con el Broker Mos Eisley entonces...")

		if choice == "GetNumberRebelds" {
			fmt.Println("Okey Preguntando")
			respuesta_host = C_Lider(choice, N_planeta,N_ciudad)
			fmt.Println("Respuesta Mos\n")
			fmt.Println(respuesta_host)
			//return
		}
		if respuesta_host == "213"{
			fmt.Printf("Vamos a guardar la wea en dist 213")
		}
		if respuesta_host == "215"{
			fmt.Printf("Vamos a proceder a guardar aqui nomas ch en 215")
		}
		if respuesta_host == "216"{
			fmt.Printf("Vamos a guardar la wea en dist 216")
		}
		fmt.Println("Comenzando nueva iteración ...")

	<-forever
	}
}
