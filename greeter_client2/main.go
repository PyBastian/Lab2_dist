package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address           = "localhost:50051"
	addressReal       = "dist214.inf.santiago.usm.cl:50051"
	nombre_informante = "Almirante Thrawn"
)

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func sendComandoToBroker(comando string, nombre_planeta string, nombre_ciudad string, nuevo_valor string) string {
	grpcChannel(comando)
	r := grpcChannel(comando)
	return r
}

func grpcChannel(message string) string {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
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

func grpcChannelServerFulctrum(comando string, nombre_planeta string, nombre_ciudad string, nuevo_valor string, direccion_maquina string) string {
	var message = comando + " " + nombre_planeta + " " + nombre_ciudad + " " + nuevo_valor
	conn, err := grpc.Dial(direccion_maquina, grpc.WithInsecure(), grpc.WithBlock())
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

func menuInformante() {
	fmt.Println("AddCity {nombre_planeta} {nombre_ciudad} [numero de rebeldes]")
	fmt.Println("UpdateName {nombre_planeta} {nombre_ciudad} [nuevo valor nombre de ciudad]")
	fmt.Println("UpdateNumber {nombre_planeta} {nombre_ciudad} [nuevo valor nro de rebeldes]")
	fmt.Println("DeleteCity {nombre_planeta} {nombre_ciudad}")
	fmt.Println("")
}

func isValidComando(comando string) bool {
	if comando != "AddCity" && comando != "UpdateName" && comando != "DeleteCity" && comando != "UpdateNumber" {
		return false
	}
	return true
}

func createComando(comando string, nombre_planeta string, nombre_ciudad string, nuevo_valor string) {
	if nuevo_valor == "" || nuevo_valor == " " {
		nuevo_valor = "0"
	}

	if isValidComando(comando) {
		var comando_completo = comando + " " + nombre_planeta + " " + nombre_ciudad + " " + nuevo_valor
		var brokerResponse = sendComandoToBroker(comando, nombre_planeta, nombre_ciudad, nuevo_valor)
		if brokerResponse == "213" {
			var direccion_servidor = "localhost:50053"
			var response = comandosInformante(comando, nombre_planeta, nombre_ciudad, nuevo_valor, direccion_servidor)
			fmt.Printf("%s", response)
			//agregar if para la respuesta del servidor y de ahi dejar el registro
			registroAcciones(nombre_informante, comando_completo, "213")
		}
		if brokerResponse == "215" {
			var direccion_servidor = "localhost:50055"
			comandosInformante(comando, nombre_planeta, nombre_ciudad, nuevo_valor, direccion_servidor)
			//agregar if para la respuesta del servidor y de ahi dejar el registro
			registroAcciones(nombre_informante, comando_completo, "215")
		}
		if brokerResponse == "216" {
			var direccion_servidor = "localhost:50056"
			comandosInformante(comando, nombre_planeta, nombre_ciudad, nuevo_valor, direccion_servidor)
			//agregar if para la respuesta del servidor y de ahi dejar el registro
			registroAcciones(nombre_informante, comando_completo, "216")
		}
	} else {
		fmt.Println("Comando no existe")
	}
}

func comandosInformante(comando string, nombre_planeta string, nombre_ciudad string, nuevo_valor string, direccion_servidor string) string {
	r := grpcChannelServerFulctrum(comando, nombre_planeta, nombre_ciudad, nuevo_valor, direccion_servidor)
	return r
}

//crea el archivo donde se guardan las acciones que ha enviado el informante
func registroAcciones(nombre_informante string, comando string, maquina string) {
	if isCreatedActionsFile(nombre_informante) {
		addRegisterToActionsFile(nombre_informante, comando, maquina)
	} else {
		createActionsFile(nombre_informante)
		addRegisterToActionsFile(nombre_informante, comando, maquina)
	}
}

func isCreatedActionsFile(nombre_informante string) bool {
	var path = nombre_informante + ".txt"
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func createActionsFile(nombre_informante string) {
	var path = nombre_informante + ".txt"

	var file, err = os.Create(path)
	if isError(err) {
		return
	}
	defer file.Close()

	fmt.Printf("\nSe ha creado un nuevo archivo de registros para el informante %s", nombre_informante)
}

func addRegisterToActionsFile(nombre_informante string, comando string, maquina string) {
	var path = nombre_informante + ".txt"

	var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if isError(err) {
		return
	}

	defer file.Close()

	_, err = file.WriteString(comando + " -> Maquina " + maquina + "\n")
	if isError(err) {
		return
	}
	fmt.Printf("\nSe ha guardado el comando ejecutado al archivo del informante %s", nombre_informante)
}

func main() {
	var comando, nombre_planeta, nombre_ciudad, nuevo_valor string

	fmt.Printf("Bienvenide Informante %s, estos seran tus comandos: \n", nombre_informante)

	menuInformante()

	forever := make(chan bool)

	for {
		fmt.Println("Ingresa el comando: ")
		fmt.Scanf("%s %s %s %s", &comando, &nombre_planeta, &nombre_ciudad, &nuevo_valor)

		fmt.Println("\n Hablemos Con el Broker Mos Eisley entonces...")

		createComando(comando, nombre_planeta, nombre_ciudad, nuevo_valor)
		<-forever
	}
}
