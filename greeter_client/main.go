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
	"strconv"
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
var reloj [3]int

func ListenInstr() {
	lis, err := net.Listen("tcp", "dist216.inf.santiago.usm.cl:50051")
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
	//text := strings.Split(in.GetName(), " ")
	fmt.Print(in.GetName())
	//usecomando(text[0], text[1], text[2], text[3])
	return &pb.HelloReply{Message: "recibid chuhc escrito"}, nil
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

func C_Lider(msg string, n_planeta string, n_ciudad string, n_valor string) string {
	//fmt.Println("Me voy a comunicar con el Lider")
	var comando string

	comando = msg + " " + n_planeta + " " + n_ciudad + " " + n_valor + " " + "1"

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
	var path_log string = "greeter_client/" + "historial" + ".txt"
	fmt.Println(path)

	if choice == "AddCity" {
		createFile(path)
		var file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0644)
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

	createFile(path_log)
	var file, err = os.OpenFile(path_log, os.O_RDWR|os.O_APPEND, 0644)
	if isError(err) {
		return
	}
	reloj[2] = reloj[2] + 1
	defer file.Close()
	_, err = file.WriteString(N_ciudad + " " + N_valor + " " + "Ahsoka Tano" + strconv.Itoa(reloj[0]) + strconv.Itoa(reloj[1]) + strconv.Itoa(reloj[2]))
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func helloworld(t time.Time) {

	var comandos string
	comandos = grpcChannel("update 1")
	fmt.Println("Envie el update")
	fmt.Println(comandos)

}

func main() {

	go ListenInstr()

	fmt.Println("Bienvenide Ahsoka, estos seran tus comandos:")

	Menu()

	for {

		var choice, N_planeta, N_ciudad, N_valor, instruccion string
		var respuesta_host string = ""
		var comando_input string

		go doEvery(30*time.Second, helloworld)

		fmt.Println("Ingresa tus comandos")
		fmt.Scanf("%s %s %s %s", &choice, &N_planeta, &N_ciudad, &N_valor)
		comando_input = choice + " " + N_planeta + " " + N_ciudad + " " + N_valor
		fmt.Println(len(instruccion))

		if choice == "close" {
			break
		}

		fmt.Println("Hablemos Con el Broker Mos Eisley entonces...")
		if choice == "AddCity" {
			fmt.Println("Okey agregemos")
			respuesta_host = C_Lider(choice, N_planeta, N_ciudad, N_valor)
			fmt.Println("El Lider fue Avisado")
			fmt.Println(respuesta_host)
			//return
		}

		if choice == "UpdateName" {
			fmt.Println("Okey uName")
			respuesta_host = C_Lider(choice, N_planeta, N_ciudad, N_valor)
			fmt.Println("El Broker fue Avisado, la info se va a la maquina")
			fmt.Println(respuesta_host)
			//return
		}

		if choice == "UpdateNumber" {
			fmt.Println("Okey uNumber")
			respuesta_host = C_Lider(choice, N_planeta, N_ciudad, N_valor)
			fmt.Println("El Broker fue Avisado, la info se va a la maquina")
			fmt.Println(respuesta_host)
			//return
		}

		if choice == "DeleteCity" {
			fmt.Println("Okey dCity")
			respuesta_host = C_Lider(choice, N_planeta, N_ciudad, N_valor)
			fmt.Println("El Broker fue Avisado, la info se va a la maquina")
			fmt.Println(respuesta_host)
			//return
		}

		if respuesta_host == "213" {
			grpcChannel213(comando_input)
			//usecomando(choice, N_planeta, N_ciudad, N_valor)
			//usecomando(choice, N_planeta, N_ciudad, N_valor)
			fmt.Printf("Vamos a guardar la wea en dist 213")
		}

		if respuesta_host == "215" {
			grpcChannel215(comando_input)
			//usecomando(choice, N_planeta, N_ciudad, N_valor)
			fmt.Printf("Vamos a guardar la wea en dist 215")
		}

		if respuesta_host == "216" {
			usecomando(choice, N_planeta, N_ciudad, N_valor)
			fmt.Printf("Vamos a proceder a guardar aqui nomas ch en 216")
			// var respuesta = grpcChannel216(comando_input)
			// fmt.Print(respuesta + " eesta wea se mando")
		}

		if respuesta_host == "Update" {
			fmt.Println("Debemos actualizar")
			// updateMaquina(comando)
		}

		fmt.Println("Finalizado, puedes ingresar nuevo comando")
	}
}
