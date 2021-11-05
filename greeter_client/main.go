package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

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

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.GetName() == "G1" || in.GetName() == "G2" || in.GetName() == "G3" {
		G_now = in.GetName()
	}
	if in.GetName() == "Ready" {
		ReadyToPlay = "Ready"
	}
	return &pb.HelloReply{Message: "recibido"}, nil
}

func grpcChannel(message string) string {

	conn, err := grpc.Dial("dist214.inf.santiago.usm.cl:50051", grpc.WithInsecure(), grpc.WithBlock())
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

func C_Lider(msg string) string {
	//fmt.Println("Me voy a comunicar con el Lider")
	var message string = ":50052" + " " + G_now + " " + R_Game + " " + id_user + " " + msg
	print(message)
	if msg == "Start" {
		//fmt.Println("Entrando al grpcChanel pa mandarle algo al Lider")
		return grpcChannel("yes")
	}

	if msg == "muerte" {
		message = ":50051" + " " + G_now + " " + R_Game + " " + id_user + " Muertos"
	}

	if msg == "resultado" {
		message = ":50051" + " " + G_now + " " + R_Game + " " + id_user + " Resultados"
	}
	if msg == "ValuePozo" {
		message = ":50051" + " " + G_now + " " + R_Game + " " + id_user + " ValuePozo"
	}

	r := grpcChannel(message)
	return r
}

func main() {

	var choice string
	forever := make(chan bool)
	go ListenInstr()

	fmt.Println("Bienvenido a SquidGame")
	fmt.Println("¿Quieres Jugar Squid game? y/n")
	fmt.Scanf("%s", &choice)

	if choice != "y" {
		return
	}
	fmt.Println("Hablemos Con el Lider entonces...")
	id_user = C_Lider("Start")
	fmt.Println("El Lider fue Avisado")

	for {

		fmt.Println("Esperando ...")
		for {
			if ReadyToPlay == "Ready" {
				break
			}
		}
		ReadyToPlay = ""

		if G_now == "G1" {
			var respuesta string
			var num string
			var ronda int = 0
			var total int = 0

			fmt.Println("Juego Luz verde Luz roja")
			fmt.Println("Si usted elige un numero mayor o igual que el lider, quedara descalificado")

			for ronda < 4 {
				fmt.Println("El lider esta decisidiendo")
				for {
					if ReadyToPlay == "Ready" {
						break
					}
				}

				fmt.Println("Elija un numero")
				fmt.Scanf("%s", &num)
				aux, _ := strconv.Atoi(num)
				total = total + aux

				respuesta = C_Lider(num)
				if respuesta == "muerte" {
					fmt.Println("Ha muerto")
					return
				}
				ReadyToPlay = ""
				ronda = ronda + 1
			}

			if total < 21 {
				fmt.Println("Ha muerto")
				_ = C_Lider("muerte")
				return
			}
			fmt.Println("El valor del pozo actual es: ", C_Lider("ValuePozo"))
			ReadyToPlay = ""
			G_now = ""
		}

		if G_now == "G2" {

			for {
				if ReadyToPlay == "Ready" {
					respuesta := C_Lider("RandomDeath")
					if respuesta == "muerte" {
						fmt.Println("A perido X(")
						return
					}
					break
				}
			}
			ReadyToPlay = ""

			fmt.Println("Bienvenido al segundo juego: Tirar la cuerda")
			fmt.Println("Debera elegir un numero entre el 1 y el 4")
			fmt.Println("Si usted elige un numero de diferente paridad que el lider, quedara descalificado")
			fmt.Println("Esperando la choiceion del Lider")
			for {
				if ReadyToPlay == "Ready" {
					break
				}
			}
			ReadyToPlay = ""

			var num string
			fmt.Println("Que numero deseas elegir_")
			fmt.Scanf("%s", &num)
			respuesta := C_Lider(num)

			fmt.Println("Esperando al resultado ...")

			if respuesta == "wait" {
				for {
					if ReadyToPlay == "Ready" {
						respuesta := C_Lider("resultado")
						if respuesta == "muerte" {
							fmt.Println("Ha muerto")
							return
						}
						fmt.Println("Gano en Tirar la cuerda")
						break
					}
				}
			}
			fmt.Println("El valor del pozo actual es: ", C_Lider("ValuePozo"))
			ReadyToPlay = ""
			G_now = ""
		}

		if G_now == "G3" {

			for {
				if ReadyToPlay == "Ready" {
					respuesta := C_Lider("RandomDeath")
					if respuesta == "muerte" {
						fmt.Println("Ha muerto")
						return
					}
					break
				}
			}
			ReadyToPlay = ""

			fmt.Println("Bienvenido al tercer juego: Todo o nada")
			fmt.Println("Debera elegir un numero entre el 1 y el 10")
			fmt.Println("Si usted elige un numero muy lejano al del lider, quedara descalificado")
			fmt.Println("Esperando al Lider")
			for {
				if ReadyToPlay == "Ready" {
					break
				}
			}

			ReadyToPlay = ""
			var num string
			fmt.Println("Elija un numero")
			fmt.Scanf("%s", &num)
			respuesta := C_Lider(num)

			fmt.Println("Esperando al resultado ...")

			if respuesta == "wait" {
				for {
					if ReadyToPlay == "Ready" {
						respuesta := C_Lider("resultado")
						if respuesta == "muerte" {
							fmt.Println("Ha muerto")
							return
						}
						fmt.Println("Gano en Todo o Nada")
						break
					}
				}
			}
			fmt.Println("El valor del pozo actual es: ", C_Lider("ValuePozo"))
			ReadyToPlay = ""
			G_now = ""
		}

	}

	<-forever
}
