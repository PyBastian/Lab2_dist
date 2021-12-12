package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"io"
	"os"
	//"strconv"
	"time"
	//"strings"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)
func isError(err error) bool {
    if err != nil {
        fmt.Println(err.Error())
    }

    return (err != nil)
}

var (
	name = flag.String("name", defaultName, "Name to greet")
)

type server struct{ pb.UnimplementedGreeterServer }

var G_now string = ""
var id_user string = ""
var ReadyToPlay string = ""
var R_Game = ""
func ListenInstr() {
	lis, err := net.Listen("tcp", "dist216.inf.santiago.usm.cl:50071")
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

func C_Lider(msg string, n_planeta string, n_ciudad string, n_valor string) string {
	//fmt.Println("Me voy a comunicar con el Lider")
	var comando string;

	comando = msg + " " + n_planeta + " " + n_ciudad + " " + n_valor

	fmt.Printf("Comando Final \n")
	fmt.Printf(comando)

	if msg == "AddCity" {
		//fmt.Println("Entrando al grpcChanel pa mandarle algo al Lider")
		return grpcChannel(comando)
	}
	if msg == "UpdateName" {
		//fmt.Println("Entrando al grpcChanel pa mandarle algo al Lider")
		return grpcChannel(comando)
	}
	if msg == "UpdateNumber" {
		//fmt.Println("Entrando al grpcChanel pa mandarle algo al Lider")
		return grpcChannel(comando)
	}
	if msg == "DeleteCity" {
		//fmt.Println("Entrando al grpcChanel pa mandarle algo al Lider")
		return grpcChannel(comando)
	}

	r := grpcChannel(msg)
	return r
}

func Menu() {
	fmt.Println("AddCity {N_planeta} {N_ciudad} [nuevo valor]")
	fmt.Println("UpdateName {N_planeta} {N_ciudad} [nuevo valor]")
	fmt.Println("UpdateNumber {N_planeta} {N_ciudad} [nuevo valor]")
	fmt.Println("DeleteCity {N_planeta} {N_ciudad} \n")

}
func createFile(path string){
	// check if file exists
	var _, err = os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
			var file, err = os.Create(path)
			if isError(err) {
					return
			}
			defer file.Close()
	}
	fmt.Println("Se ha agregado un nuevo registro", path)
}

func usecomando(choice string, N_planeta string ,N_ciudad string ,N_valor string){
	var path string
	path = N_planeta + ".txt"
	if choice == "AddCity"{
		createFile(path)
		var file, err = os.OpenFile(path, os.O_RDWR, 0644)
		if isError(err) {
				return
		}
		defer file.Close()
		_, err = file.WriteString(N_ciudad + " " + N_valor + "\n")
		if isError(err) {
				return
		}
		// Read file, line by line
		var text = make([]byte, 1024)
		for {
				_, err = file.Read(text)
				// Break if finally arrived at end of file
				if err == io.EOF {
						break
				}
				// Break if error occured
				if err != nil && err != io.EOF {
						isError(err)
						break
				}
		}
		fmt.Println(string(text))
	}
	if choice == "UpdateName"{

	}
	if choice == "UpdateNumber"{

	}
	if choice == "DeleteCity "{
    var err = os.Remove(path)
    if isError(err) {
        return
    }
    fmt.Println("")
	}

}

func main() {

	var choice, N_planeta, N_ciudad, N_valor string
	var respuesta_host string

	forever := make(chan bool)
	go ListenInstr()

	fmt.Println("Bienvenide Informante Ahsoka Tano, estos seran tus comandos:\n")

	Menu()

	fmt.Scanf("%s %s %s %s", &choice, &N_planeta, &N_ciudad, &N_valor)
  //fmt.Printf("captured: %s %s %s %s\n", choice, N_planeta, N_ciudad, N_valor)

	fmt.Println("Hablemos Con el Broker Mos Eisley entonces...")

	for {
		if choice == "AddCity" {
			fmt.Println("Okey agregemos")
			respuesta_host = C_Lider(choice, N_planeta,N_ciudad,N_valor)
			fmt.Println("El Lider fue Avisado")
			fmt.Println(respuesta_host)
			//return
		}
		if choice == "UpdateName" {
			fmt.Println("Okey uName")
			respuesta_host = C_Lider(choice, N_planeta,N_ciudad,N_valor)
			fmt.Println("El Broker fue Avisado, la info se va a la maquina")
			fmt.Println(respuesta_host)
			//return
		}
		if choice == "UpdateNumber" {
			fmt.Println("Okey uNumber")
			respuesta_host = C_Lider(choice, N_planeta,N_ciudad,N_valor)
			fmt.Println("El Broker fue Avisado, la info se va a la maquina")
			fmt.Println(respuesta_host)
			//return
		}
		if choice == "DeleteCity" {
			fmt.Println("Okey dCity")
			respuesta_host = C_Lider(choice, N_planeta,N_ciudad,N_valor)
			fmt.Println("El Broker fue Avisado, la info se va a la maquina")
			fmt.Println(respuesta_host)
			//return
		}
		if respuesta_host == "213"{
			fmt.Printf("Vamos a guardar la wea en dist 213")
		}
		if respuesta_host == "215"{
			fmt.Printf("Vamos a guardar la wea en dist 215")
		}
		if respuesta_host == "216"{
			usecomando(choice, N_planeta,N_ciudad,N_valor)
			fmt.Printf("Vamos a proceder a guardar aqui nomas ch en 216")

		}
		fmt.Println("Comenzando nueva iteración ...")

	<-forever
	}
}
