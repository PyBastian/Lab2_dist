package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	//"os"
	"strconv"
	"strings"
	"time"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port       = ":50051"
	addrs_pozo = "dist215.inf.santiago.usm.cl:50053"
	addrs_node = ":50054"
)

type server struct{ pb.UnimplementedGreeterServer }

var TotalPlayer int = 1
var N_play int = 0
var N_Playerr int = 0
var L_player = [3]string{"l", "l", "l"}

var numberG1 int = 0

var RPlayerEliminated string = "-"
var T_T1 int = 0
var TotalG2T2 int = 0
var choiceG2 int = 0
var TeamPlayers = [3]string{"-", "-", "-"}
var LoseTeam string = "0"

var numberG3 int = 0
var PairPlayers = [3]int{-1, -1, -1}
var AnswerPlayers = [3]int{-1, -1, -1}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func grpcChannel(ipAdress string, message string) string {
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
//Enviamos mensaje al cliente
func SMPlayer(msgLider string, id_user int) {
	fmt.Println("Mandando info al cliente")
	var message string
	var Eliminar_user int = id_user

	if msgLider == "Round" {
		message = "Ready"
	}
	if msgLider == "1" || msgLider == "2" || msgLider == "3" {
		fmt.Println("Entramos al if")
		message = "E" + msgLider
	}
	if msgLider == "D" || msgLider == "DT" {
		N_play = N_play - 1
		if msgLider == "D" {
			Eliminar_user = A_id_user()
		}
		_ = SendMessageToPozo("", strconv.FormatInt(int64(Eliminar_user), 10))
		message = "death " + strconv.FormatInt(int64(Eliminar_user), 10)
	}
	_ = grpcChannel("dist216.inf.santiago.usm.cl:50071", message)
}

func SendMessageToPozo(msg string, player string) string {
	if msg == "val" {
		return grpcChannel(addrs_pozo, msg)
	}
	return ""
}

func MsgToNode(msg string) string {
	return grpcChannel(addrs_node, msg)
}


func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("Comando Final server \n")
	//214 no porque es la direcciones del Broker
	dir_maquinas := []string{"213", "215", "216"}

	fmt.Printf(in.GetName())

	randomIndex := rand.Intn(3)
	selected_value := dir_maquinas[randomIndex]

	text := strings.Split(in.GetName(), " ")
	//fmt.Printf(text)
	if text[0] == "GetNumberRebelds"{
			return &pb.HelloReply{Message: "Ligerito te entregamos respsuesta"}, nil
	}
	fmt.Printf(text[0])

	return &pb.HelloReply{Message: selected_value}, nil
}


func A_id_user() int {
	IDAleatorio := rand.Intn(N_play) + 1
	for i := 0; i < len(L_player); i++ {
		if L_player[i] == "l" {
			IDAleatorio = IDAleatorio - 1
		}
		if IDAleatorio == 0 && L_player[i] == "l" {
			return i + 1
		}
	}
	return 0
}

func FindPair(t int, p int) int {
	var aux int = 0
	for i := 0; i < len(L_player); i++ {
		if PairPlayers[i] == t && (p-1) != i {
			aux = i
		}
	}
	return aux
}

func TeamsG2() {
	var aux int = 0
	for i := 0; i < len(L_player); i++ {
		if aux == 0 && L_player[i] == "l" {
			TeamPlayers[i] = "1"
		}
		if aux == 1 && L_player[i] == "l" {
			TeamPlayers[i] = "2"
		}
		aux = (aux + 1) % 2
	}
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

func main() {

	go ListenMessage()

	forever := make(chan bool)
	var choice string

	fmt.Println("Esperando solicitudes")

	for {
		//Menu()
		fmt.Scanf("%s", &choice)
		SMPlayer(choice, 0)

		if choice == "0" {
			aux := SendMessageToPozo("val", "")
			fmt.Println("Valor Actual del pozo: ", aux)
		}
		//Casos de Juegos {1,2,3}
		if choice == "1" {
			SMPlayer("Round", 0)
			fmt.Println("Primer juego")
			fmt.Println("Deberas elegir 4 numeros entre 6 y 10")
			var n_rondas = 4
			for round := 0; round < n_rondas; round++ {
				fmt.Println("Elija un numero")
				fmt.Scanf("%d", &numberG1)

				SMPlayer("Round", 0)

				fmt.Println("Esperando jugadores", N_Playerr, "de", N_play)
				for {
					if N_Playerr == N_play {
						break
					}
				}

				if N_play == 0 {
					fmt.Println("Todos los jugadores murieron")
					break
				}
				N_Playerr = 0
			}

			fmt.Println("Juego finalizado")
			fmt.Println("Jugadores sobrevivientes ", N_Playerr)
			LivePlayers()
		}
		//Juego 2
	<-forever
}
