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

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port                = ":50051"
	addressPozoGRPC     = ":50053"
	addressNameNodeGRPC = ":50054"
)

type server struct{ pb.UnimplementedGreeterServer }

//VARIABLES GENERALES
var MaxPlayers int = 1
var NumberOfPlayers int = 0
var NumberOfPlayersReady int = 0
var ListOfLivePlayers = [3]string{"y", "y", "y"}

//VARIABLES JUEGO 1
var numberG1 int = 0

//VARIABLES JUEGO 2
var RPlayerEliminated string = "-"
var TotalG2T1 int = 0
var TotalG2T2 int = 0
var choiceG2 int = 0
var TeamPlayers = [3]string{"-", "-", "-"}
var LoseTeam string = "0"

//VARIABLE JUEGO3
var numberG3 int = 0
var PairPlayers = [3]int{-1, -1, -1}
var AnswerPlayers = [3]int{-1, -1, -1}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//CONEXIONES
func grpcChannel(ipAdress string, message string) string {
	fmt.Println("Entramos al grpChannel_1")
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

func rabbitmqChannel(message string) {
	conn, err := amqp.Dial("amqp://guest:guest@:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	body := message + " 1 100"
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(body)})
	failOnError(err, "Failed to publish a message")
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

//Enviamos mensaje al cliente
func SMPlayer(msgLider string, IDplayer int) {
	fmt.Println("Mandando info al cliente")
	var message string
	var UserToEliminated int = IDplayer

	if msgLider == "Round" {
		message = "Ready"
	}
	if msgLider == "1" || msgLider == "2" || msgLider == "3" {
		fmt.Println("Entramos al if")
		message = "G" + msgLider
	}
	if msgLider == "D" || msgLider == "DT" {
		NumberOfPlayers = NumberOfPlayers - 1
		if msgLider == "D" {
			UserToEliminated = A_IDplayer()
		}
		_ = SendMessageToPozo("", strconv.FormatInt(int64(UserToEliminated), 10))
		message = "death " + strconv.FormatInt(int64(UserToEliminated), 10)
	}
	fmt.Println("Estamos a punto de entrar al grpcChannel")
	_ = grpcChannel("dist216.inf.santiago.usm.cl:50071", message)

}

//MANDAR MENSAJES AL POZO
func SendMessageToPozo(msg string, player string) string {
	if msg == "val" {
		return grpcChannel(addressPozoGRPC, msg)
	}
	rabbitmqChannel(player)
	return ""
}

func SendMessageToNameNode(msg string) string {
	return grpcChannel(addressNameNodeGRPC, msg)
}

// ESCUCHAR MENSAJES DE LOS JUGADORES
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Entramos al SayHello")
	//INGRESAR AL JUEGO
	if in.GetName() == "yes" {
		NumberOfPlayers = NumberOfPlayers + 1
		fmt.Println("Esperando a los Jugadores, llevamos = ", NumberOfPlayers, " de ", MaxPlayers)
		return &pb.HelloReply{Message: strconv.FormatInt(int64(NumberOfPlayers), 10)}, nil
	}

	//SEPARACION DEL MENSAJE
	text := strings.Split(in.GetName(), " ")
	IDplayer, _ := strconv.Atoi(text[3])
	Jugada, _ := strconv.Atoi(text[4])

	if text[4] != "R" && text[4] != "RD" && text[4] != "VP" {
		_ = SendMessageToNameNode(text[1] + " " + text[3] + " " + text[4])
	}

	if text[4] == "R" {
		if text[1] == "Game2" && TeamPlayers[IDplayer-1] == LoseTeam {
			_ = SendMessageToPozo("", text[3])
			fmt.Println("Jugador " + text[3] + " ha muerto")

			ListOfLivePlayers[IDplayer-1] = "n"
			NumberOfPlayers = NumberOfPlayers - 1
			fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)

			return &pb.HelloReply{Message: "death"}, nil
		}
		if text[1] == "Game3" {
			var aux int = FindPair(PairPlayers[IDplayer-1], IDplayer)
			if AnswerPlayers[IDplayer-1] > AnswerPlayers[aux] {
				_ = SendMessageToPozo("", text[3])
				fmt.Println("Jugador " + text[3] + " ha muerto")

				ListOfLivePlayers[IDplayer-1] = "n"
				NumberOfPlayers = NumberOfPlayers - 1
				fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)

				return &pb.HelloReply{Message: "death"}, nil
			}
		}

		NumberOfPlayersReady = NumberOfPlayersReady + 1
		fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)
		return &pb.HelloReply{Message: "live"}, nil
	}

	if text[4] == "RD" {
		if RPlayerEliminated == text[3] {
			_ = SendMessageToPozo("", text[3])
			fmt.Println("Jugador " + text[3] + " ha muerto")

			ListOfLivePlayers[IDplayer-1] = "n"
			NumberOfPlayers = NumberOfPlayers - 1
			fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)

			return &pb.HelloReply{Message: "death"}, nil
		}

		NumberOfPlayersReady = NumberOfPlayersReady + 1
		fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)
		return &pb.HelloReply{Message: "live"}, nil
	}

	if text[4] == "VP" {
		return &pb.HelloReply{Message: SendMessageToPozo("val", "")}, nil
	}

	if text[1] == "G1" {
		if Jugada >= numberG1 || text[4] == "death" {
			_ = SendMessageToPozo("", text[3])
			fmt.Println("Jugador " + text[3] + " ha muerto")

			ListOfLivePlayers[IDplayer-1] = "n"
			NumberOfPlayers = NumberOfPlayers - 1
			fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)
			return &pb.HelloReply{Message: "death"}, nil
		}

		NumberOfPlayersReady = NumberOfPlayersReady + 1
		fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)
		return &pb.HelloReply{Message: "live"}, nil
	}

	if text[1] == "G2" {
		if TeamPlayers[IDplayer-1] == "1" {
			TotalG2T1 = TotalG2T1 + Jugada
		}
		if TeamPlayers[IDplayer-1] == "2" {
			TotalG2T2 = TotalG2T2 + Jugada
		}
		NumberOfPlayersReady = NumberOfPlayersReady + 1
		fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)
		return &pb.HelloReply{Message: "wait"}, nil
	}

	if text[1] == "G3" {
		AnswerPlayers[IDplayer-1] = Jugada
		NumberOfPlayersReady = NumberOfPlayersReady + 1
		fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)
		return &pb.HelloReply{Message: "wait"}, nil
	}

	return nil, nil
}

func LivePlayers() {
	for i := 0; i < len(ListOfLivePlayers); i++ {
		if ListOfLivePlayers[i] == "y" {
			fmt.Println("Jugador ", i+1, "ha sobrevivido")
		}
	}
}

func Menu() {

	fmt.Println("Elija 0 para ver el valor del pozo")
	fmt.Println("Elija 1 para comenzar el juego Luz Roja Luz Verde")
	fmt.Println("Elija 2 para comenzar el juego Tirar la cuerda")
	fmt.Println("Elija 3 para comenzar el juego Todo o nada")

}

func A_IDplayer() int {
	IDAleatorio := rand.Intn(NumberOfPlayers) + 1
	for i := 0; i < len(ListOfLivePlayers); i++ {
		if ListOfLivePlayers[i] == "y" {
			IDAleatorio = IDAleatorio - 1
		}
		if IDAleatorio == 0 && ListOfLivePlayers[i] == "y" {
			return i + 1
		}
	}
	return 0
}

func FindPair(t int, p int) int {
	var aux int = 0
	for i := 0; i < len(ListOfLivePlayers); i++ {
		if PairPlayers[i] == t && (p-1) != i {
			aux = i
		}
	}
	return aux
}

func DefineTeamsG2() {
	var aux int = 0
	for i := 0; i < len(ListOfLivePlayers); i++ {
		if aux == 0 && ListOfLivePlayers[i] == "y" {
			TeamPlayers[i] = "1"
		}
		if aux == 1 && ListOfLivePlayers[i] == "y" {
			TeamPlayers[i] = "2"
		}
		aux = (aux + 1) % 2
	}
}

func DefineTeamsG3() {
	var aux int = 0
	for i := 0; i < len(ListOfLivePlayers); i++ {
		if ListOfLivePlayers[i] == "y" {
			PairPlayers[i] = aux
		}
		aux = (aux + 1) % (NumberOfPlayers / 2)
	}
}

//MAIN
func main() {

	go ListenMessage()

	forever := make(chan bool)
	var choice string
	fmt.Println("Los jugadores estan Entrando!! ", NumberOfPlayers, "/", MaxPlayers)
	for {
		if NumberOfPlayers == MaxPlayers {
			break
		}
	}

	for {
		Menu()
		fmt.Scanf("%s", &choice)
		SMPlayer(choice, 0)

		if choice == "0" {
			aux := SendMessageToPozo("val", "")
			fmt.Println("Valor Actual del pozo: ", aux)
		}
		//Casos de Juegos
		if choice == "1" {
			SMPlayer("Round", 0)
			fmt.Println("Primer juego")
			fmt.Println("Debe elegir 4 numeros entre 6 y 10")
			for round := 0; round < 4; round++ {
				fmt.Println("Elija un numero")
				fmt.Scanf("%d", &numberG1)

				SMPlayer("Round", 0)

				fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)
				for {
					if NumberOfPlayersReady == NumberOfPlayers {
						break
					}
				}

				if NumberOfPlayers == 0 {
					fmt.Println("Todos los jugadores murieron")
					break
				}
				NumberOfPlayersReady = 0
			}

			fmt.Println("Juego finalizado")
			fmt.Println("Jugadores sobrevivientes ", NumberOfPlayersReady)
			LivePlayers()
		}

		if choice == "2" {

			NumberOfPlayersReady = 0
			SMPlayer("R", 0)
			if NumberOfPlayers%2 == 1 && NumberOfPlayers != 1 {
				RPlayerEliminated = strconv.FormatInt(int64(A_IDplayer()), 10)
				SMPlayer("Round", 0)
				for {
					if NumberOfPlayersReady == NumberOfPlayers {
						break
					}
				}
			}
			NumberOfPlayersReady = 0

			DefineTeamsG2()
			fmt.Println("Segundo juego")
			fmt.Println("Debe elegir un numero entre 1 y 4")

			fmt.Scanf("%d", &choiceG2)
			choiceG2 = choiceG2 % 2

			SMPlayer("Round", 0)

			for {
				if NumberOfPlayersReady == NumberOfPlayers {
					break
				}
			}
			NumberOfPlayersReady = 0

			if TotalG2T1%2 == TotalG2T2%2 && TotalG2T1 != choiceG2 {
				aux := rand.Intn(2)
				if aux == 0 {
					LoseTeam = "1"
				}
				if aux == 1 {
					LoseTeam = "2"
				}
			}

			if TotalG2T1%2 != choiceG2 {
				LoseTeam = "1"
			}
			if TotalG2T2%2 != choiceG2 {
				LoseTeam = "2"
			}

			fmt.Println("Equipo perdedor ", LoseTeam)
			SMPlayer("Round", 0)
			for {
				if NumberOfPlayersReady == NumberOfPlayers {
					break
				}
			}

			fmt.Println("Juego finalizado")
			fmt.Println("Jugadores sobrevivientes ", NumberOfPlayersReady)
			LivePlayers()
			NumberOfPlayersReady = 0
			//SendMessageToPlayers("R", 0)
		}

		if choice == "3" {

			NumberOfPlayersReady = 0
			SMPlayer("Round", 0)
			if NumberOfPlayers%2 == 1 && NumberOfPlayers != 1 {
				RPlayerEliminated = strconv.FormatInt(int64(A_IDplayer()), 10)
				SMPlayer("R", 0)
				for {
					if NumberOfPlayersReady == NumberOfPlayers {
						break
					}
				}
			}
			NumberOfPlayersReady = 0

			DefineTeamsG3()

			fmt.Println("Tercer juego")
			fmt.Println("Debe elegir un numero entre 1 y 10")

			fmt.Println("Elija un numero")
			fmt.Scanf("%d", &numberG3)

			SMPlayer("Round", 0)
			fmt.Println("Esperando jugadores", NumberOfPlayersReady, "/", NumberOfPlayers)
			for {
				if NumberOfPlayersReady == NumberOfPlayers {
					break
				}
			}
			NumberOfPlayersReady = 0

			for i := 0; i < len(ListOfLivePlayers); i++ {
				if ListOfLivePlayers[i] == "y" {
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
				if NumberOfPlayersReady == NumberOfPlayers {
					break
				}
			}

			fmt.Println("Juego finalizado")
			fmt.Println("Los jugadores ganadores son ", NumberOfPlayersReady)
			LivePlayers()
			NumberOfPlayersReady = 0
		}
	}
	<-forever
}
