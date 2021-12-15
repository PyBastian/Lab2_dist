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
	var comando string

	comando = msg + " " + n_planeta + " " + n_ciudad + " " + n_valor + " " + "1"

	if msg == "AddCity" {
		return grpcChannel(comando)
	}
	if msg == "UpdateName" {
		return grpcChannel(comando)
	}
	if msg == "UpdateNumber" {
		return grpcChannel(comando)
	}
	if msg == "DeleteCity" {
		return grpcChannel(comando)
	}

	r := grpcChannel(msg)
	return r
}

func Menu() {
	fmt.Println("AddCity {N_planeta} {N_ciudad} [nuevo valor]")
	fmt.Println("UpdateName {N_planeta} {N_ciudad} [nuevo valor]")
	fmt.Println("UpdateNumber {N_planeta} {N_ciudad} [nuevo valor]")
	fmt.Println("DeleteCity {N_planeta} {N_ciudad}")
	fmt.Println("")
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

func updateMaquina(t time.Time) {

	fmt.Println("Preguntando si hay que hacer merge...")
	var comandos string

	comandos = grpcChannel("update 2")

	fmt.Printf("Comandos recibidos: %s", comandos)

}

func main() {

	go ListenInstr()

	fmt.Println("Bienvenide Ahsoka, estos seran tus comandos:")

	Menu()

	for {
		go doEvery(120*time.Second, updateMaquina)

		var choice, N_planeta, N_ciudad, N_valor string
		var respuesta_host string = ""

		fmt.Println("Ingresa tus comandos")
		fmt.Scanf("%s %s %s %s", &choice, &N_planeta, &N_ciudad, &N_valor)

		if choice == "close" {
			break
		}

		fmt.Println("Hablemos Con el Broker Mos Eisley entonces...")
		if choice == "AddCity" {
			respuesta_host = C_Lider(choice, N_planeta, N_ciudad, N_valor)
			fmt.Println("El Broker fue Avisado")
		}

		if choice == "UpdateName" {
			respuesta_host = C_Lider(choice, N_planeta, N_ciudad, N_valor)
			fmt.Println("El Broker fue Avisado, la info se va a la maquina")
		}

		if choice == "UpdateNumber" {
			respuesta_host = C_Lider(choice, N_planeta, N_ciudad, N_valor)
			fmt.Println("El Broker fue Avisado, la info se va a la maquina")
		}

		if choice == "DeleteCity" {
			respuesta_host = C_Lider(choice, N_planeta, N_ciudad, N_valor)
			fmt.Println("El Broker fue Avisado, la info se va a la maquina")
		}

		if respuesta_host == "216" {
			usecomando(choice, N_planeta, N_ciudad, N_valor)
		}

		fmt.Println("Finalizado, puedes ingresar nuevo comando")
	}
}
