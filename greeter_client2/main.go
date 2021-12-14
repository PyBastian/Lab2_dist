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

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func ListenInstr() {
	lis, err := net.Listen("tcp", "dist213.inf.santiago.usm.cl:50071")
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

	//text := strings.Split(in.GetName(), " ")
	//fmt.Printf(text)
	selected_value := "Te entregsmos desde 213"

	// if text[0] == "GetNumberRebelds " {
	// 	return &pb.HelloReply{Message: "Ligerito te entregamos respsuesta"}, nil
	// }
	// fmt.Printf(text[0], text[1], text[2], text[3])
	// usecomando(text[0], text[1], text[2], text[3])

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

// func grpcChannel216(message string) string {
// 	fmt.Println("Enviando mensaje a 216")
// 	fmt.Println(message)
// 	conn, err := grpc.Dial("dist216.inf.santiago.usm.cl:50071", grpc.WithInsecure(), grpc.WithBlock())
// 	if err != nil {
// 		log.Fatalf("Error de conecc'on con 213: %v", err)
// 	}
// 	defer conn.Close()
// 	c := pb.NewGreeterClient(conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()
// 	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: message})
// 	if err != nil {
// 		log.Fatalf("could not greet: %v", err)
// 	}
// 	return r.GetMessage()
// }

func grpcChannel215(choice string, N_planeta string, N_ciudad string, N_valor string) string {
	fmt.Println("Enviando mensaje a 215")
	fmt.Println(choice + " " + N_planeta + " " + N_ciudad + " " + N_valor)
	var message string = "devuelve la wea qlo"
	conn, err := grpc.Dial("dist215.inf.santiago.usm.cl:50051", grpc.WithInsecure(), grpc.WithBlock())
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

func grpcChannel216(choice string, N_planeta string, N_ciudad string, N_valor string) string {
	fmt.Println("Enviando mensaje a 216")
	var message string = "devuelve la wea qlo"
	conn, err := grpc.Dial("dist216.inf.santiago.usm.cl:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Error de conecc'on con 216: %v", err)
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

	comando = msg + " " + n_planeta + " " + n_ciudad + " " + n_valor + " " + "2"

	fmt.Printf("Comando Final \n")

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
	fmt.Println("DeleteCity {N_planeta} {N_ciudad}")
	fmt.Println("")
}

func usecomando(choice string, N_planeta string, N_ciudad string, N_valor string) {
	var path string
	path = N_planeta + ".txt"
	fmt.Println(path)

	if choice == "AddCity" {
		AddCity(path, N_ciudad, N_valor)
	}

	if choice == "UpdateName" {
		UpdateName(path, N_ciudad, N_valor)
	}

	if choice == "UpdateNumber" {
		UpdateNumber(path, N_ciudad, N_valor)
	}
	if choice == "DeleteCity " {
		DeleteCity(path)
	}

}

func updateMaquina(comandos []string) {
	//esta funcion deberia ser llamada cada 2 minutos para ejecutar los comandos que se puedan haber usado en otra maquina

	for i := 0; i < len(comandos); i++ {
		comando := strings.Split(string(comandos[i]), " ")
		path := comando[1] + ".txt"
		N_ciudad := comando[2]
		N_valor := "0"
		if len(comando) == 4 {
			N_valor = comando[3]
		}

		switch comando[0] {
		case "AddCity":
			AddCity(path, N_ciudad, N_valor)
		case "UpdateName":
			UpdateName(path, N_ciudad, N_valor)
		case "UpdateNumber":
			UpdateNumber(path, N_ciudad, N_valor)
		case "DeleteCity":
			DeleteCity(path)
		}
	}
}

func AddCity(path string, N_ciudad string, N_valor string) {
	if N_valor == "" {
		N_valor = "0"
	}

	if !isPlanetFileCreated(path) {
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
	} else {
		var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if isError(err) {
			return
		}

		defer file.Close()

		if N_valor == "" || N_valor == " " {
			N_valor = "0"
		}

		_, err = file.WriteString(N_ciudad + " " + N_valor + "\n")
		if isError(err) {
			return
		}
		fmt.Printf("\nSe ha agregado la ciudad %s al registro", N_ciudad)

	}

}
func isPlanetFileCreated(path string) bool {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func UpdateName(path string, N_ciudad string, N_valor string) {
	var didchange int = 0
	input, err := ioutil.ReadFile(path)
	lines := strings.Split(string(input), "\n")
	if isError(err) {
		return
	}
	for i, line := range lines {
		if strings.Contains(line, N_ciudad) {
			lines[i] = N_ciudad + " " + N_valor
			fmt.Printf("\nEl nombre de la ciudad se actualizo correctamente")
			didchange = 1
			break
		}
	}
	if didchange == 0 {
		fmt.Printf("\nNo se encontro el nombre de la ciudad")
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	didchange = 0
}

func UpdateNumber(path string, N_ciudad string, N_valor string) {
	var didchange int = 0
	input, err := ioutil.ReadFile(path)
	lines := strings.Split(string(input), "\n")
	if isError(err) {
		return
	}
	for i, line := range lines {
		if strings.Contains(line, N_ciudad) {
			lines[i] = N_ciudad + " " + N_valor
			fmt.Print("\nEl numero de rebeldes fue actualizado correctamente")
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

func DeleteCity(path string) {
	var err = os.Remove(path)
	if isError(err) {
		return
	}
	fmt.Println("Deleted")
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
	fmt.Println("\nSe ha agregado un nuevo registro", path)
}

func serverResponse(choice string, N_planeta string, N_ciudad string, N_valor string) {
	var respuesta_host string

	respuesta_host = C_Lider(choice, N_planeta, N_ciudad, N_valor)

	// respuesta_host = "216"

	if respuesta_host == "213" {
		usecomando(choice, N_planeta, N_ciudad, N_valor)
		fmt.Printf("Vamos a guardar la wea en dist 213")
	}

	if respuesta_host == "215" {
		fmt.Printf("Vamos a guardar la wea en dist 215")
		grpcChannel215(choice, N_planeta, N_ciudad, N_valor)
		//usecomando(choice, N_planeta, N_ciudad, N_valor)
	}

	if respuesta_host == "216" {
		fmt.Printf("Vamos a proceder a guardar aqui nomas ch en 216")
		grpcChannel216(choice, N_planeta, N_ciudad, N_valor)
		// var respuesta = grpcChannel216(comando_input)
		// fmt.Print(respuesta + " eesta wea se mando")
	}

}

func main() {

	fmt.Println("Bienvenide Almirante Thrawn, estos seran tus comandos:")

	for {

		Menu()

		go ListenInstr()

		var choice, N_planeta, N_ciudad, N_valor string

		fmt.Println("Ingresa tus comandos")
		fmt.Scanf("%s %s %s %s", &choice, &N_planeta, &N_ciudad, &N_valor)

		if choice == "close" {
			break
		}

		fmt.Println("Hablemos Con el Broker Mos Eisley entonces...")

		serverResponse(choice, N_planeta, N_ciudad, N_valor)

		// if respuesta_host == "update" {
		// 	fmt.Println("Debemos actualizar")
		// 	// updateMaquina(comando)
		// }

		fmt.Println("Finalizado, puedes ingresar nuevo comando")
	}
}
