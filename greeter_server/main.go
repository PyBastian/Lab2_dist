package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
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

	if in.GetName() == "yes" {
		N_play = N_play + 1
		fmt.Println("Esperando a los Jugadoress, llevamos = ", N_play, " de ", TotalPlayer)
		return &pb.HelloReply{Message: "Wena hermano aqui server, te mando un saludo camarada"}, nil
	}

	text := strings.Split(in.GetName(), " ")
	id_user, _ := strconv.Atoi(text[3])
	Jugada, _ := strconv.Atoi(text[4])

	if text[4] != "R" && text[4] != "RD" && text[4] != "VP" {
		_ = MsgToNode(text[1] + " " + text[3] + " " + text[4])
	}

	if text[4] == "R" {
		if text[1] == "Game2" && TeamPlayers[id_user-1] == LoseTeam {
			_ = SendMessageToPozo("", text[3])
			fmt.Println("Jugador " + text[3] + " ha muerto")

			L_player[id_user-1] = "n"
			N_play = N_play - 1
			fmt.Println("Esperando jugadores", N_Playerr, "/", N_play)

			return &pb.HelloReply{Message: "death"}, nil
		}
		if text[1] == "Game3" {
			var aux int = FindPair(PairPlayers[id_user-1], id_user)
			if AnswerPlayers[id_user-1] > AnswerPlayers[aux] {
				_ = SendMessageToPozo("", text[3])
				fmt.Println("Jugador " + text[3] + " ha muerto")

				L_player[id_user-1] = "n"
				N_play = N_play - 1
				fmt.Println("Esperando jugadores, llevamos", N_Playerr, "de", N_play)

				return &pb.HelloReply{Message: "death"}, nil
			}
		}

		N_Playerr = N_Playerr + 1
		fmt.Println("Esperando jugadores", N_Playerr, "/", N_play)
		return &pb.HelloReply{Message: "live"}, nil
	}

	if text[4] == "RD" {
		if RPlayerEliminated == text[3] {
			_ = SendMessageToPozo("", text[3])
			fmt.Println("Jugador " + text[3] + " ha muerto")

			L_player[id_user-1] = "n"
			N_play = N_play - 1
			fmt.Println("Esperando jugadores", N_Playerr, " de ", N_play)

			return &pb.HelloReply{Message: "death"}, nil
		}

		N_Playerr = N_Playerr + 1
		fmt.Println("Esperando jugadores", N_Playerr, "/", N_play)
		return &pb.HelloReply{Message: "live"}, nil
	}

	if text[4] == "VP" {
		return &pb.HelloReply{Message: SendMessageToPozo("val", "")}, nil
	}

	if text[1] == "E1" {
		if Jugada >= numberG1 || text[4] == "death" {
			_ = SendMessageToPozo("", text[3])
			fmt.Println("Jugador " + text[3] + " ha muerto")

			L_player[id_user-1] = "n"
			N_play = N_play - 1
			fmt.Println("Esperando jugadores", N_Playerr, "/", N_play)
			return &pb.HelloReply{Message: "death"}, nil
		}

		N_Playerr = N_Playerr + 1
		fmt.Println("Esperando jugadores", N_Playerr, "/", N_play)
		return &pb.HelloReply{Message: "live"}, nil
	}

	if text[1] == "E2" {
		if TeamPlayers[id_user-1] == "1" {
			T_T1 = T_T1 + Jugada
		}
		if TeamPlayers[id_user-1] == "2" {
			TotalG2T2 = TotalG2T2 + Jugada
		}
		N_Playerr = N_Playerr + 1
		fmt.Println("Esperando jugadores", N_Playerr, "/", N_play)
		return &pb.HelloReply{Message: "wait"}, nil
	}

	if text[1] == "E3" {
		AnswerPlayers[id_user-1] = Jugada
		N_Playerr = N_Playerr + 1
		fmt.Println("Esperando jugadores", N_Playerr, "/", N_play)
		return &pb.HelloReply{Message: "wait"}, nil
	}

	return nil, nil
}

func LivePlayers() {
	for i := 0; i < len(L_player); i++ {
		if L_player[i] == "l" {
			fmt.Println("Jugador ", i+1, "ha sobrevivido")
		}
	}
}

func Menu() {
	fmt.Println("Presione 0 para ver el valor del pozo")
	fmt.Println("Presiones 1 para comenzar el juego Luz Roja Luz Verde")
	fmt.Println("Presiones 2 para comenzar el juego Tirar la cuerda")
	fmt.Println("Presiones 3 para comenzar el juego Todo o nada")
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
		if choice == "2" {

			N_Playerr = 0
			SMPlayer("R", 0)
			if N_play%2 == 1 && N_play != 1 {
				RPlayerEliminated = strconv.FormatInt(int64(A_id_user()), 10)
				SMPlayer("Round", 0)
				for {
					if N_Playerr == N_play {
						break
					}
				}
			}
			N_Playerr = 0

			TeamsG2()
			fmt.Println("Segundo juego")
			fmt.Println("Elige un numero entre 1 y 4")

			fmt.Scanf("%d", &choiceG2)
			choiceG2 = choiceG2 % 2

			SMPlayer("Round", 0)

			for {
				if N_Playerr == N_play {
					break
				}
			}
			N_Playerr = 0

			if T_T1%2 == TotalG2T2%2 && T_T1 != choiceG2 {
				aux := rand.Intn(2)
				if aux == 0 {
					LoseTeam = "1"
				}
				if aux == 1 {
					LoseTeam = "2"
				}
			}

			if T_T1%2 != choiceG2 {
				LoseTeam = "1"
			}
			if TotalG2T2%2 != choiceG2 {
				LoseTeam = "2"
			}

			fmt.Println("Perdedores ", LoseTeam)
			SMPlayer("Round", 0)
			for {
				if N_Playerr == N_play {
					break
				}
			}

			fmt.Println("Juego finalizado")
			fmt.Println("Jugadores restantes ", N_Playerr)
			LivePlayers()
			N_Playerr = 0
		}

		if choice == "3" {

			N_Playerr = 0
			SMPlayer("Round", 0)
			if N_play%2 == 1 && N_play != 1 {
				RPlayerEliminated = strconv.FormatInt(int64(A_id_user()), 10)
				SMPlayer("R", 0)
				for {
					if N_Playerr == N_play {
						break
					}
				}
			}
			N_Playerr = 0

			fmt.Println("Tercer juego")
			fmt.Println("Elige entre ore 1 al 10")

			fmt.Println("Escoge un numero")
			fmt.Scanf("%d", &numberG3)

			SMPlayer("Round", 0)
			fmt.Println("Esperando jugadores", N_Playerr, "/", N_play)
			for {
				if N_Playerr == N_play {
					break
				}
			}
			N_Playerr = 0

			for i := 0; i < len(L_player); i++ {
				if L_player[i] == "l" {
					if AnswerPlayers[i] >= numberG3 {
						AnswerPlayers[i] = AnswerPlayers[i] - numberG3
					}
					if AnswerPlayers[i] < numberG3 {
						AnswerPlayers[i] = numberG3 - AnswerPlayers[i]
					}
				}
			}

			SMPlayer("Round", 0)
			for {
				if N_Playerr == N_play {
					break
				}
			}

			fmt.Println("SquidGame has stop")
			fmt.Println("Ganaron ", N_Playerr)
			LivePlayers()
			N_Playerr = 0
		}
	}
	<-forever
}
