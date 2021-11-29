package main

import (
    "os"
    "log"
    "strconv"
    "google.golang.org/grpc"
    "context"
    pb "lab/game/proto"
    "bufio"
    "net"
    "fmt"
)

type server struct {
	pb.UnimplementedCalamardoGameServer
}

// Funcion que recibe jugadas y las manda al archivo
func (s *server) SavePlays(ctx context.Context, in *pb.SavePlaysRequest) (*pb.SavePlaysReply, error) {

    player := in.GetId()
    round := in.GetRonda()
    plays := in.GetJugada()
    log.Println("")
    log.Println("SavePlays: \tEl jugador: ", player, "esta guardando las jugadas: ", plays)

    fileName := ""
    fileName = "jugador_" + strconv.Itoa(int(player)) + "__ronda_" + strconv.Itoa(int(round)) + ".txt"
    f, err := os.Create(fileName)
    var index int64 = 0
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    data := []byte(plays)
    _, err = f.WriteAt(data, index)
    index = 1 + index
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Guardado con exito...")
    log.Println("")
    // Que pasa cuando no funciona? como retorno false?
	return &pb.SavePlaysReply{Exito: true}, nil
}

func (s *server) SendPlays(ctx context.Context, in *pb.SendPlaysRequest) (*pb.SendPlaysReply, error){
    // necesita recibir round y player
    

    player := int(in.GetId())
    round := in.GetRonda()
    log.Println("")
    log.Println("SendPlays \tEl Jugador: ", player, "quiere saber sus jugadas en la ronda: ", round)

    plays := ""
    f, err := os.Open("jugador_" + strconv.Itoa(player) + "__ronda_" + string(round) + ".txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
		plays = plays + scanner.Text() + "\n"
    }

    log.Println("Encontramos las jugadas: ", plays)
    log.Println("")
    return &pb.SendPlaysReply{Jugada:plays}, nil
    // retorna string jugadas
}

func main() {


    log.Printf("Bienvenido al DataNode, iniciando servicios...")
	go func() {
		listner, err := net.Listen("tcp", ":8082")

		if err != nil {
			panic("cannot create tcp connection " + err.Error())
		}

		serv := grpc.NewServer()
		pb.RegisterCalamardoGameServer(serv, &server{})
		if err = serv.Serve(listner); err != nil {
			panic("cannot initialize the server" + err.Error())
		}
	}()
    
    var terminar string = ""


    for terminar != "T"{
        fmt.Println("Si desea terminar el programa, escriba 'T'")
        fmt.Scanln(&terminar)
    }
    
}