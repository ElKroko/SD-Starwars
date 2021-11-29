package main

import (
	"bufio"
	"context"
	"fmt"
	pb "lab/game/proto"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalamardoGameServer
}

var index int64 = 0

// Pregunto por un jugador, y me envia todas las jugadas que tiene este jugador.
// Va a buscar las jugadas al datanode que le corresponde, y retorna las jugada hacia el Lider
// un jugador pupede tener las jugadas en diferentes datanodes, hay que revisarlos todos
// El namenode recibe la request del lider, y luego hace una request a los datanodes sin todavia enviar al Lider el resultado.

// Cuando recibe el resultado del Datanode, y encuentra TODAS LAS JUGADAS, retorna al lider.

func preguntar_datanode(id int32, ronda int32, ip string) (jugada string) {
	
	fmt.Println("entre al datanode")
	
	fmt.Println("El Jugador: ", id, "quiere saber sus jugadas")
	ip_y_puerto := ip + ":8082"
	conn, err := grpc.Dial(ip_y_puerto, grpc.WithInsecure())

	if err != nil {
		panic("cannot connect with server " + err.Error())
	}

	serviceClient := pb.NewCalamardoGameClient(conn)

	res, err := serviceClient.SendPlays(context.Background(), &pb.SendPlaysRequest{Id: id, Ronda: ronda})

	if err != nil {
		panic("No se pudo anadir la jugada en el namenode! " + err.Error())
	}

	jugada = res.GetJugada()
	fmt.Println("Encontre las jugadas!")
	fmt.Println("")
	return jugada
}

func (s *server) ReturnPlays(ctx context.Context, in *pb.ReturnPlaysRequest) (*pb.ReturnPlaysReply, error) {
	fmt.Println("")
	log.Printf("ReturnPlays: \tEl jugador_%d esta enviando sus datos", in.GetId())
	log.Println("Para saber de sus jugadas...")

	f, err := os.Open("lugarJugadores.txt")
	var splitText []string
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var ronda int32
	var rondaString string
	var ip string
	var jugadas string
	returnString := ""
	for scanner.Scan() {
		fmt.Println("entre al scanner!")
		splitText = strings.Split(scanner.Text(), " ")
		fmt.Println("split text es: ",splitText[0])
		fmt.Println("id text es: ",splitText[0][len(splitText[0])-1])

		idText, _ := strconv.Atoi(string(splitText[0][len(splitText[0])-1]))
		fmt.Println("idText en int es: " ,idText)
		if idText ==  int(in.GetId()) {

			// Conectar con el Datanode
			ronda = int32(splitText[1][len(splitText[1])-1])
			rondaString = string(ronda)
			ip = splitText[2]
			jugadas = preguntar_datanode(in.GetId(), ronda, ip)
			returnString = returnString + "Juego " + rondaString + "\n" + jugadas
		}
	}

	fmt.Println("Jugadas: ", jugadas)
	return &pb.ReturnPlaysReply{Jugada: returnString}, nil
}

// Funcion utilizada por el lider para enviar las jugadas hacia el Namenode despues de cada juego.

func enviar_datanode(ip string, id int32, ronda int32, jugada string) (exito bool) {
	ip_y_puerto := ip + ":8082"
	conn, err := grpc.Dial(ip_y_puerto, grpc.WithInsecure())

	if err != nil {
		panic("cannot connect with server " + err.Error())
	}

	serviceClient := pb.NewCalamardoGameClient(conn)

	res, err := serviceClient.SavePlays(context.Background(), &pb.SavePlaysRequest{Id: id, Ronda: ronda, Jugada: jugada})

	exito = res.GetExito()

	return exito

}

func (s *server) InsertPlays(ctx context.Context, in *pb.InsertPlaysRequest) (*pb.InsertPlaysReply, error) {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("")
	log.Printf("InsertPlays: \tEl jugador_%d esta enviando sus datos", in.GetId())

	f, err := os.OpenFile("lugarJugadores.txt", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	id := in.GetId()
	ronda := in.GetRonda()
	jugada := in.GetJugada()
	var ip string
	

	selectDataNode := rand.Intn(3)
	if selectDataNode == 0 {
		ip = "10.6.43.109"
	} else if selectDataNode == 1 {
		ip = "10.6.43.110"
	} else {
		ip = "10.6.43.111"
	}
	val := "Jugador_" + strconv.Itoa(int(id)) + " Ronda_" + strconv.Itoa(int(ronda)) + " " + ip + "\n"
	data := []byte(val)
	_, err = f.WriteAt(data, index)
	index = int64(len(data)) + index
	if err != nil {
		log.Fatal(err)
	}
	//mandar jugada a la ip

	// Tengo que enviar: jugador_id, ronda_id (para el nombre del archivo)
	// y las jugadas
	exito := enviar_datanode(ip, id, ronda, jugada)
	if exito {
		log.Println("Operacion lograda con exito!")
	} else {
		log.Println("Operacion fallida :ccccc")
	}

	fmt.Println("")
	return &pb.InsertPlaysReply{Exito: exito}, nil
}

func main() {

	// Instalar el namenode en el ip 112
	log.Printf("Bienvenido al NameNode, iniciando servicios...")
	go func() {
		listner, err := net.Listen("tcp", ":8081")

		if err != nil {
			panic("cannot create tcp connection " + err.Error())
		}

		serv := grpc.NewServer()
		pb.RegisterCalamardoGameServer(serv, &server{})
		if err = serv.Serve(listner); err != nil {
			panic("cannot initialize the server" + err.Error())
		}
	}()

	f, err := os.Create("lugarJugadores.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var terminar string = ""

	for terminar != "T" {
		fmt.Println("Si desea terminar el programa, escriba 'T'")
		fmt.Scanln(&terminar)
	}

}
