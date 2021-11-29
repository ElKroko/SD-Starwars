package main

import ("fmt"
		"os"
		"log"
        "bufio"
        "strings"
        "time"
)

var index int64 = 0

func (s *server) PlayerDie(ctx context.Context, in *pb.JoinRequest) (*pb.JoinReply, error) {
	log.Printf("Recibido: %s", in.GetMessage())

    f, err := os.Open("JugadoresEliminados.txt")
    val := mensaje
    data := []byte(val)
    _, err := f.WriteAt(data, index)
    index = int64(len(data))+index
    if err != nil {
        log.Fatal(err)
    }

	return &pb.JoinReply{IdJugador: int32(totalPlayers-1), Alive: true, Round: 0}, nil
}

func ActualAmmount() string {
    f, err := os.Open("JugadoresEliminados.txt")
    var actualAmmount string
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        actualAmmount = strings.Split(scanner.Text(), " ")[2]
    }
    return actualAmmount
}

func main() {

    go func() {
		listner, err := net.Listen("tcp", ":8083")
		//conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

		if err != nil {
			panic("cannot connect with server " + err.Error())
		}

		serv := grpc.NewServer()
		pb.RegisterSquidGameServiceServer(serv, &server{})
		if err = serv.Serve(listner); err != nil {
			panic("cannot initialize the server" + err.Error())

		}
	}()

    f, err := os.Create("JugadoresEliminados.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    log.Printf("El Pozo ha comenzado...")
	go func() {
		listner, err := net.Listen("tcp", ":8080")

		if err != nil {
			panic("cannot create tcp connection " + err.Error())
		}

		serv := grpc.NewServer()
		pb.RegisterPozoServer(serv, &server{})
		if err = serv.Serve(listner); err != nil {
			panic("cannot initialize the server" + err.Error())
		}
	}()

    for !finish{
        time.Sleep(2* time.Second)
    }

}