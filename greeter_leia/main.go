package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

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
	fmt.Printf("Recibimos Comando \n")
	fmt.Printf(in.GetName())

	text := strings.Split(in.GetName(), " ")
	//fmt.Printf(text)
	if text[0] == "GetNumberRebelds " {
		return &pb.HelloReply{Message: "Ligerito te entregamos respsuesta"}, nil
	}
	fmt.Printf(text[0], text[1], text[2], text[3])
	usecomando(text[0], text[1], text[2], text[3])

	return &pb.HelloReply{Message: selected_value}, nil

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
	var comando string

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

func grpcChannel213(message string) string {
	fmt.Println("Enviando mensaje a 213")
	fmt.Println(message)

	conn, err := grpc.Dial("dist213.inf.santiago.usm.cl:50071", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Error de conecc'on con 213: %v", err)
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

func grpcChannel215(message string) string {
	fmt.Println("Enviando mensaje a 215")
	fmt.Println(message)
	conn, err := grpc.Dial("dist215.inf.santiago.usm.cl:50071", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Error de conecc'on con 215: %v", err)
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

func createFile(path string) {
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

func usecomando(choice string, N_planeta string, N_ciudad string, N_valor string) {
	var path string
	var didchange int = 0
	path = "greeter_client/" + N_planeta + ".txt"
	fmt.Println(path)

	if choice == "AddCity" {
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

	if choice == "UpdateName" {
		input, err := ioutil.ReadFile(path)
		lines := strings.Split(string(input), "\n")
		if isError(err) {
			return
		}
		for i, line := range lines {
			if strings.Contains(line, N_ciudad) {
				lines[i] = N_ciudad + " " + N_valor
				fmt.Printf("El nombre de la ciudad se actualizo correctamente")
				didchange = 1
				break
			}
		}
		if didchange == 0 {
			fmt.Printf("No se encontro el nombre de la ciudad")
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(path, []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}
		didchange = 0
	}

	if choice == "UpdateNumber" {
		input, err := ioutil.ReadFile(path)
		lines := strings.Split(string(input), "\n")
		if isError(err) {
			return
		}
		for i, line := range lines {
			if strings.Contains(line, N_ciudad) {
				lines[i] = N_ciudad + " " + N_valor
				fmt.Print("El numero de rebeldes fue actualizado correctamente")
				didchange = 1
				break
			}
		}
		if didchange == 0 {
			fmt.Printf("No se encontro el nombre de la ciudad")
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(path, []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}
		didchange = 0
	}
	if choice == "DeleteCity " {
		var err = os.Remove(path)
		if isError(err) {
			return
		}
		fmt.Println("")
	}

}

func Menu() {
	fmt.Println("GetNumberRebelds {N_planeta} {N_ciudad}")
}

func main() {

	var choice, N_planeta, N_ciudad, N_valor string
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
			respuesta_host = C_Lider(choice, N_planeta, N_ciudad)
			fmt.Println("Respuesta Mos\n")
			fmt.Println(respuesta_host)
			//return
		}

		if respuesta_host == "213" {
			fmt.Printf("Vamos a guardar la wea en dist 213")
		}
		if respuesta_host == "215" {
			fmt.Printf("Vamos a proceder a guardar aqui nomas ch en 215")
			usecomando(choice, N_planeta, N_ciudad, N_valor)
		}
		if respuesta_host == "216" {
			fmt.Printf("Vamos a guardar la wea en dist 216")
		}
		fmt.Println("Comenzando nueva iteración ...")
		<-forever
	}
}
