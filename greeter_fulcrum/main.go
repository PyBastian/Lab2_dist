package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	//"os"
	//"strconv"

	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	//revisar puertos de maquinas virtuales
	port = ":50053"
)

type server struct{ pb.UnimplementedGreeterServer }

func grpcChannelServerFulctrum(ipAdress string, message string) string {
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
	fmt.Printf("\n Recibimos Conexión ...")

	text := strings.Split(in.GetName(), " ")

	var comando, nombre_planeta, nombre_ciudad, nuevo_valor string

	comando = text[0]
	nombre_planeta = text[1]
	nombre_ciudad = text[2]
	nuevo_valor = text[3]

	if comando == "AddCity" {
		addCity(nombre_planeta, nombre_ciudad, nuevo_valor)
	}
	if comando == "UpdateName" {
		updateNameToCity(nombre_planeta, nombre_ciudad, nuevo_valor)
	}
	if comando == "UpdateNumber" {
		updateNumberToCity(nombre_planeta, nombre_ciudad, nuevo_valor)
	}
	if comando == "DeleteCity" {
		deleteCity(nombre_planeta, nombre_ciudad)
	}

	return &pb.HelloReply{Message: "vector qlo loco"}, nil
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

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

//funcionalidades a partir de los comandos
//funcionalidad addCity
func addCity(nombre_planeta string, nombre_ciudad string, nuevo_valor string) {
	if isPlanetFileCreated(nombre_planeta) {
		addRegisterToPlanet(nombre_planeta, nombre_ciudad, nuevo_valor)
	} else {
		createPlanet(nombre_planeta)
		addRegisterToPlanet(nombre_planeta, nombre_ciudad, nuevo_valor)
	}
}

func isPlanetFileCreated(nombre_planeta string) bool {
	var path = nombre_planeta + ".txt"
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func createPlanet(nombre_planeta string) {
	var path = nombre_planeta + ".txt"

	var file, err = os.Create(path)
	if isError(err) {
		return
	}
	defer file.Close()

	fmt.Printf("\nSe ha creado un nuevo archivo para el planeta %s", nombre_planeta)
}

//agregar registro a archivo de planetas, esto va en el servidor
func addRegisterToPlanet(nombre_planeta string, nombre_ciudad string, nuevo_valor string) {
	var path = nombre_planeta + ".txt"

	var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if isError(err) {
		return
	}

	defer file.Close()

	if nuevo_valor == "" || nuevo_valor == " " {
		nuevo_valor = "0"
	}

	_, err = file.WriteString(nombre_ciudad + " " + nuevo_valor + "\n")
	if isError(err) {
		return
	}
	fmt.Printf("\nSe ha agregado la ciudad %s al archivo del planeta %s", nombre_ciudad, nombre_planeta)
}

//funcionalidad actualizar ciudad, esto va en el servidor
func updateNameToCity(nombre_planeta string, nombre_ciudad string, nuevo_valor string) {
	if !isPlanetFileCreated(nombre_planeta) {
		fmt.Printf("\n El planeta solicitado no existe")
		return
	}

	var path = nombre_planeta + ".txt"
	var flag = false

	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, nombre_ciudad) {
			var number = strings.Split(line, " ")
			lines[i] = nuevo_valor + " " + number[1]
			flag = true
		}
	}

	if !flag {
		fmt.Printf("\n La ciudad solicitada no existe")
		return
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("\n La ciudad %s fue cambiada por %s del planeta %s", nombre_ciudad, nuevo_valor, nombre_planeta)
}

//funcionalidad de actualizar numero de rebeldes, esto va en el servidor
func updateNumberToCity(nombre_planeta string, nombre_ciudad string, nuevo_valor string) {
	if !isPlanetFileCreated(nombre_planeta) {
		fmt.Printf("\n El planeta solicitado no existe")
		return
	}

	var path = nombre_planeta + ".txt"
	var flag = false

	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, nombre_ciudad) {
			lines[i] = nombre_ciudad + " " + nuevo_valor
			flag = true
		}
	}

	if !flag {
		fmt.Printf("\n La ciudad solicitada no existe")
		return
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("\n La ciudad %s del planeta %s fue actualizada a %s rebeldes", nombre_ciudad, nombre_planeta, nuevo_valor)
}

//funcionalidad de eliminar una ciudad, esto va en el servidor
func deleteCity(nombre_planeta string, nombre_ciudad string) {
	if !isPlanetFileCreated(nombre_planeta) {
		fmt.Printf("\n El planeta solicitado no existe")
		return
	}

	var path = nombre_planeta + ".txt"
	var flag = false

	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, nombre_ciudad) {
			lines[i] = ""
			flag = true
		}
	}

	newFile := []string{}
	for i := range lines {
		if lines[i] != "" {
			newFile = append(newFile, lines[i])
		}
	}

	if !flag {
		fmt.Printf("\n La ciudad solicitada no existe")
		return
	}

	output := strings.Join(newFile, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("\n La ciudad %s del planeta %s fue eliminada", nombre_ciudad, nombre_planeta)
}

func main() {

	go ListenMessage()

	forever := make(chan bool)
	var choice string

	fmt.Println("Esperando...")

	for {

		fmt.Scanf("%s", &choice)
		<-forever
	}
}
