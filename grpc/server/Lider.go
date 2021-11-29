package main

import (
	"context"
	"fmt"
	pb "lab/game/proto"
	"log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)




var total int = 2 									// Aqui se define el total de jugadores maximos!



type server struct {
	pb.UnimplementedCalamardoGameServer
}

type PlayerStruct struct {
	id    int32
	alive bool
	score int32
	jugada int32
	etapa int32
}


type Jugadas struct {
	id    int32
	ronda int
	jugada string
}

// Hacer una funcion que busque en el array de jugadas x id
// Luego, segun la etapa, append a la lista o inicializar valor de jugada 

// Variables
var totalPlayers int								// Jugadores actuales
var players_online []PlayerStruct
var started bool

var playsBoss []int
var jugadas1 []Jugadas
var jugadas2 []Jugadas
var jugadas3 []Jugadas



var etapaActual int
var total_juego1 int = 0
var total_juego2 int
var total_juego3 int = 0



// Funciones de mensajeria
func (s *server) JoinGame(ctx context.Context, in *pb.JoinRequest) (*pb.JoinReply, error) {
	log.Printf("Recibido: %s", in.GetMessage())
	
	players_online = append(players_online, PlayerStruct{int32(totalPlayers), true, 0, 0, 0})
	totalPlayers += 1

	return &pb.JoinReply{IdJugador: int32(totalPlayers-1), Alive: true, Round: 0}, nil
}

func (s *server) StartGame(ctx context.Context, in *pb.StartRequest) (*pb.StartReply, error) {

	// log.Printf("El jugador %d envio: %s \t en la etapa: %d",in.GetId(), in.GetMessage(), in.GetEtapa())
	etapa := in.GetEtapa()

	switch etapa {
	case 0:
		if len(players_online) == total {					// Cantidad de jugadores = maxio
			return &pb.StartReply{Started: true}, nil
		} else {
			return &pb.StartReply{Started: false}, nil
		}
	case 1:
		if started{
			return &pb.StartReply{Started: true}, nil
		} else {
			return &pb.StartReply{Started: false}, nil
		}
	case 2:
		if started{
			return &pb.StartReply{Started: true}, nil
		} else {
			return &pb.StartReply{Started: false}, nil
		}
	default:
		return &pb.StartReply{Started: false}, nil
	}
}


// Funcion para enviar y recibir jugadas desde el lider, para todos los juegos.

func (s *server) JuegoMsg(ctx context.Context, in *pb.JuegoRequest) (*pb.JuegoReply, error) {
	vivo := true
	etapa_jugador := in.GetEtapa()
	var temp_jugador Jugadas

	if etapaActual == 1 && etapa_jugador == 1{
		ronda := in.GetRound()
		id_jugador := in.GetId()
		playPlayer := in.GetJugada()
		suma := in.GetScore()
		log.Println("El Jugador ", id_jugador, "en la ronda ", ronda)
		log.Println("Ha jugado la carta", playPlayer)

		

		if playPlayer >= int32(playsBoss[ronda]) && playPlayer >= 1 && playPlayer <= 10{
			RemovePlayer(indexOfPlayers(id_jugador))
			vivo = false
			total_juego1 +=1
		} else {
			suma = suma + playPlayer
			if suma >= 21 {
				etapa_jugador += 1
				log.Println("")
				log.Println("El Jugador ", id_jugador, "en la ronda ", ronda)
				log.Println("Paso a la siguiente etapa! con puntaje ", suma)
				ronda = 0														 // Reiniciamos la ronda para la etapa 2
				started = false
				total_juego1 += 1
			} else if ronda == 3 {
				vivo = false
				RemovePlayer(indexOfPlayers(id_jugador))
				total_juego1 += 1
			} else {
				ronda = ronda + 1
			}
		}
		indexJugada := indexOfJugadas(jugadas1, in.GetId())
		if indexJugada == -1 {
			temp_jugador.id = in.GetId()
			temp_jugador.jugada = strconv.Itoa(int(in.GetJugada()))
			temp_jugador.ronda = 1
			jugadas1 = append(jugadas1, temp_jugador)
		} else {
			jugadas1[indexJugada].jugada = jugadas1[indexJugada].jugada + "\n" + strconv.Itoa(int(in.GetJugada()))
		}
		
		return &pb.JuegoReply{Alive: vivo, Round: ronda, Score: suma, Etapa: etapa_jugador}, nil


	} else if etapa_jugador == 2{
		// Aqui va la logica del juego 2

		ronda := in.GetRound()
		suma := in.GetScore()


		fmt.Println("\t Recibi mensaje juego 2")
		log.Println("El Jugador ", in.GetId(), "en la ronda ", in.GetRound())
		log.Println("Ha jugado la carta", in.GetJugada())
		// Recibir el mensaje, y devolver todo igual... 
		players_online[indexOfPlayers(in.Id)].jugada = in.GetJugada()
		players_online[indexOfPlayers(in.Id)].etapa = in.GetEtapa()
		// Para decidir si el jugador vivio o no, se usara una funcion alive() que verificara si estoy vivo al final... 
		// Preguntando si el id de el jugador esta en la lista.
		total_juego2 += 1
		etapa_jugador += 1
		started = false
		temp_jugador.id = in.GetId()
		temp_jugador.jugada = strconv.Itoa(int(in.GetJugada()))
		temp_jugador.ronda = 2
		jugadas2 = append(jugadas2, temp_jugador)
		return &pb.JuegoReply{Alive: true, Round: ronda, Score: suma, Etapa: etapa_jugador}, nil


	} else if etapa_jugador == 3{
		// Aqui va la logica del juego 3

		ronda := in.GetRound()
		suma := in.GetScore()

		fmt.Println("\tRecibi mensaje juego 3")
		log.Println("El Jugador ", in.GetId(), "en la ronda ", in.GetRound())
		log.Println("Ha jugado la carta", in.GetJugada())
		total_juego3 += 1
		players_online[indexOfPlayers(in.Id)].jugada = in.GetJugada()
		players_online[indexOfPlayers(in.Id)].etapa = in.GetEtapa()
		etapa_jugador += 1
		started = false
		temp_jugador.id = in.GetId()
		temp_jugador.jugada = strconv.Itoa(int(in.GetJugada()))
		temp_jugador.ronda = 3
		jugadas3 = append(jugadas3, temp_jugador)
		return &pb.JuegoReply{Alive: true, Round: ronda, Score: suma, Etapa: etapa_jugador}, nil


	} else{
		fmt.Println("Algo extraño paso!")
		return &pb.JuegoReply{Alive: vivo, Round: 8, Score: -2, Etapa: etapa_jugador}, nil
	}
	
}

// Funcion que recibe el id de un player, lo busca en la lista players_online y si no lo encuentra,
// retorna alive = false
func (s *server) Muerte (ctx context.Context, in *pb.MuerteRequest) (*pb.MuerteReply, error) {
	fmt.Printf("El jugador %d esta preguntando si esta muerto", in.GetId())
	id := in.GetId()
	vivo:= true

	if indexOfPlayers(id) == -1 {
		vivo = false
	} else {
		vivo = true
	}

	return &pb.MuerteReply{Id:id, Alive:vivo}, nil
}

// Funciones del Juego


// Funciones Auxiliares
func indexOfPlayers(element int32) (int) {
	fmt.Println("")
	fmt.Println("Entre a Index of Players")
	fmt.Println("Estoy buscando la posicion para el id:", element)
	for k, v := range players_online {
		fmt.Println("Indice array: ", k)
		fmt.Println("Valor array: ", v)

		if element == v.id {
			fmt.Println("Lo encontre! con valor ", k)
			fmt.Println("")
			return k
		}
	}
	fmt.Println("No encontre nada para ", element)
	return -1    //not found.
}
func indexOfJugadas(array []Jugadas, element int32) (int) {
	for k, v := range array {
		if element == v.id {
			fmt.Println("Lo encontre! con valor ", k)
			fmt.Println("")
			return k
		}
	}
	fmt.Println("No encontre nada IndexOfJugadas para", element)
	return -1    //not found.
}

func indexOf(array []PlayerStruct, element int32) (int) {
	for k, v := range array {
		if element == v.id {
			return k
		}
	}
	fmt.Println("No encontre nada IndexOf para", element)
	return -1    //not found.
}

func RemovePlayer(i int){
	value := players_online[i].id 
    players_online[i] = players_online[len(players_online)-1]
	players_online = players_online[:len(players_online)-1]
	fmt.Printf("El jugador_%d fue eliminado\n",value)
	
	
	//avisar a Pozo 
}

func RemoveElemArray(array []PlayerStruct, i int) []PlayerStruct {
    array[i] = array[len(array)-1]
    return array[:len(array)-1]
}

func print_id_player (players []PlayerStruct) {
	//log.Println("")
	//log.Println("==========================")
	//log.Println("        Ganadores         ")
	//log.Println("==========================")
	log.Println("")
	log.Println("Lista: ")
	log.Print("\t id: ")
	for _, v := range players {
		log.Print(v.id, " ,")
	}
	log.Println("")
}

func print_id_player_ganadores (players []PlayerStruct) {
	log.Println("")
	log.Println("==========================")
	log.Println("        Ganadores         ")
	log.Println("==========================")
	log.Println("")
	log.Println("Lista: ")
	log.Print("\t id: ")
	for _, v := range players {
		log.Print(v.id, " ,")
	}
	log.Println("")
}


func mandar_jugadas (jugadas []Jugadas ) {
	
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure()) // Conectamos al IP de 10.6.43.109:8080, el lider.

	if err != nil {
		panic("cannot connect with server " + err.Error())
	}

	serviceClient := pb.NewCalamardoGameClient(conn)
	
	for _, value := range jugadas {
		id := value.id
		ronda := int32(value.ronda)
		jugada := value.jugada

		// Usar InsertPlays
		res, err := serviceClient.InsertPlays(context.Background(), &pb.InsertPlaysRequest{Id: id, Ronda: ronda, Jugada: jugada})

		if err != nil {
			panic("No se pudo anadir la jugada en el namenode! " + err.Error())
		}

		exito := res.GetExito()

		if exito {
			log.Println("Se ha logrado enviar la jugada con exito! para el id: ", id)
		}
	}
}


func pedir_jugadas (jugadas []Jugadas, id_preguntar int) (resultado string) {
	// Preguntamos por un solo id, el cual sera buscado en el struct.

	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure()) // Conectamos al IP de 10.6.43.109:8080, el lider.

	if err != nil {
		panic("cannot connect with server " + err.Error())
	}

	serviceClient := pb.NewCalamardoGameClient(conn)

	id := int32(id_preguntar)

	res, err := serviceClient.ReturnPlays(context.Background(), &pb.ReturnPlaysRequest{Id: id})

	if err != nil {
		panic("No se pudo anadir la jugada en el namenode! " + err.Error())
	}

	jugada_string := res.GetJugada()

	resultado = "id: " + strconv.Itoa(id_preguntar) +" jugadas: " + jugada_string

	return resultado
}


func main() {
	log.Printf("Bienvenido al Calamardo, iniciando servicios...")
	go func() {
		listner, err := net.Listen("tcp", ":8080")

		if err != nil {
			panic("cannot create tcp connection " + err.Error())
		}

		serv := grpc.NewServer()
		pb.RegisterCalamardoGameServer(serv, &server{})
		if err = serv.Serve(listner); err != nil {
			panic("cannot initialize the server" + err.Error())
		}
	}()

	

	started = false
	totalPlayers = 0
	etapaActual = 0
	// Calamardo := ""
	// var muere bool

	var empezar string
	
	

	for totalPlayers < total {
		time.Sleep(2 * time.Second)

		if totalPlayers != total {
			fmt.Println("Aun no hay suficientes calamares... Total:", totalPlayers, "\t esperamos:", total)
		}
	}

	fmt.Println("")
	fmt.Println("Todos a bordo!")
	fmt.Println("")
	started = false

	
	// ==============================================================================
	// 									Inicio Juego 1
	// ==============================================================================


	// Logica juego 1
	fmt.Println("--------------------------------")
	fmt.Println("       Inicio del juego 1       ")
	fmt.Println("--------------------------------")
	fmt.Println("")
	fmt.Println("Hola!, el lider eligira sus cartas")

	var numero int
	//gano := 0
	ronda := 0
	etapaActual = 1

	for ronda < 4 {
		fmt.Println("Elija un numero del 6 al 10")
		fmt.Print("-> ")
		fmt.Scanln(&numero)
		playsBoss = append(playsBoss, numero)
		ronda = ronda + 1
	}

	started = true


	fmt.Println("Jugadores: seleccionen su carta... ")
	fmt.Println("Luego de recibir jugadas, presionar ENTER en esta consola.")
	fmt.Scanln(&empezar)

	started = false


	// mandar las jugadas al namenode
	mandar_jugadas(jugadas1)

	// Preguntar si quiere continuar, o preguntar por las jugadas al datanode

	var opcion string
	flag_opcion := true

	for flag_opcion {
		fmt.Println("Que deseas hacer?")
		fmt.Println("\t [1] Preguntar por jugadas a Datanode")
		fmt.Println("\t [2] Seguir el juego")
		fmt.Scanln(&opcion)
		if opcion == "1" {
			// Preguntar x jugadas al datanode
			var id_preguntar string
			fmt.Println("Ingrese el id del jugador que desea preguntar: ")
			fmt.Print("id: ")
			fmt.Scanln(&id_preguntar)
			id_int, err := strconv.Atoi(id_preguntar)

			if err != nil {
				log.Println("No pude convertir str a int " + err.Error())
			}

			jugadas := pedir_jugadas(jugadas1, id_int)

			log.Println("Las jugadas para el id: ", id_int, "son: ")
			log.Println(jugadas)

			flag_opcion = false
		} else if opcion == "2" {
			flag_opcion = false
		} else {
			fmt.Println("Lo siento, ingresaste una opcion incorrecta...")
			fmt.Println("")
		}
	}
	


	// ==============================================================================
	// 									Inicio Juego 2
	// ==============================================================================
	var bossNumber int
	var teamPlayers []PlayerStruct



	if len(players_online) > 1 {
		// Logica Juego 2
		rand.Seed(time.Now().UnixNano())
		fmt.Println("")
		fmt.Println("--------------------------------")
		fmt.Println("       Inicio del juego 2       ")
		fmt.Println("--------------------------------")
		fmt.Println("")
		fmt.Println("Elija un número del 1 al 4")
		fmt.Scanln(&bossNumber)
		bossParity := bossNumber % 2
		fmt.Println(bossParity)

		total_juego2 = 0
		
		for total_juego1 < len(players_online){
			fmt.Println("Aun no han jugado todos...")
			time.Sleep(2* time.Second)
	
		}


		started = true

		// Esperar a que los jugadores manden su jugada

		for total_juego2 < len(players_online){
			fmt.Println("Aun no han jugado todos...")
			time.Sleep(2* time.Second)
		}
		
		started = false
		// Luego de que hagan su jugada, los separamos en equipo
		// Y calculamos la paridad


		var team1 []PlayerStruct
		var team2 []PlayerStruct

		teamPlayers = players_online 
		var newPlayer1 int
		var newPlayer2 int


			// Cambiar team_players = players_online
		for len(teamPlayers) > 0 {
			if len(teamPlayers) == 1 {
				RemovePlayer(indexOfPlayers(teamPlayers[0].id))
				fmt.Println("Eliminamos al jugador %d, habia un jugador solo", teamPlayers[0].id)
				teamPlayers = RemoveElemArray(teamPlayers, indexOfPlayers(teamPlayers[0].id))
			} else {
				//ingresar juegador a team 1
				newPlayer1 = rand.Intn(len(teamPlayers))				// al azar entre todos los players online
				fmt.Println("\t new player: ")
				team1 = append(team1, teamPlayers[newPlayer1])				// 
				teamPlayers = RemoveElemArray(teamPlayers, newPlayer1)
				//ingresar juegador a team 2
				newPlayer2 = rand.Intn(len(teamPlayers))
				team2 = append(team2, teamPlayers[newPlayer2])
				teamPlayers = RemoveElemArray(teamPlayers, newPlayer2)
			}
		}


		// cuando todos los jugadores me hayan enviado su jugada, calcular que pasa con la paridad
		team1Sum := 0
		team2Sum := 0
		number := 0
		for number < len(team1) {
			indexPlayer:= 0
			fmt.Println("\t Busco para team 1")
			indexPlayer = indexOfPlayers(team1[number].id)
			fmt.Println("indice: ", indexPlayer)
			fmt.Println("Jugada", players_online[indexPlayer].jugada)
			fmt.Println("")

			team1Sum += int(players_online[indexPlayer].jugada)
			number ++
		}
		number = 0
		for number < len(team2) {
			indexPlayer:= 0
			fmt.Println("\t Busco para team 2")
			indexPlayer = indexOfPlayers(team2[number].id)
			fmt.Println("indice: ", indexPlayer)
			fmt.Println("Jugada", players_online[indexPlayer].jugada)
			fmt.Println("")
			team2Sum += int(players_online[indexPlayer].jugada)
			number ++
		}

		

		team1Parity := team1Sum % 2
		team2Parity := team2Sum % 2

		fmt.Println("Paridad Boss: ", bossParity)
		fmt.Println("Paridad team 1: ", team1Parity)
		fmt.Println("Paridad team 2: ", team2Parity)

		// liberar el ping. (listo = true) y actualizar alive (decidiendo si murio o no, para cerrar procesos)

		if bossParity == team1Parity && bossParity == team2Parity {
			fmt.Println("Nadie muerio en otra ronda")
		} else if bossParity == team1Parity && bossParity != team2Parity {
			fmt.Println("Gano Team 1")
			for _, player := range team2 {
				RemovePlayer(int(indexOfPlayers(player.id)))
			}
		} else if bossParity != team1Parity && bossParity == team2Parity {
			fmt.Println("Gano Team 2")
			for _, player := range team1 {
				RemovePlayer(int(indexOfPlayers(player.id)))
			}
		} else {
			fmt.Println("Nadie gano, decidiendo a quien nos echamos...")
			if rand.Intn(2) == 0 {
				for _, player := range team1 {
					RemovePlayer(int(indexOfPlayers(player.id)))
				}
			} else {
				for _, player := range team2 {
					RemovePlayer(int(indexOfPlayers(player.id)))
				}
			}
		}
	} else if len(players_online) == 1{
		RemovePlayer(0)
		fmt.Println("Eliminamos, habia solo 1 jugador")
	}
	
	mandar_jugadas(jugadas2)
	// ==============================================================================
	// 									Inicio Juego 3
	// ==============================================================================

	started = true
	total_juego3 = 0
	fmt.Println("")
	fmt.Println("--------------------------------")
	fmt.Println("       Inicio del juego 3       ")
	fmt.Println("--------------------------------")
	fmt.Println("")
	fmt.Println("Los jugadores deberan enviar su jugada primero...")

	for total_juego3 < len(players_online){
		fmt.Println("Aun no han jugado todos...")
		time.Sleep(2* time.Second)

	}
	// Esperar mensajes de vuelta 
	started = false


	mandar_jugadas(jugadas3)

	// --------------------------------
	// 			Comunicacion
	// --------------------------------

	rand.Seed(time.Now().UnixNano())
	
	if len(players_online)%2 == 1 {
		RemovePlayer(rand.Intn(len(players_online))) // Resolver problemas de paridad...
	} 
	if len(players_online) > 0 {		
		fmt.Println("")
		fmt.Println("Elija un número del 1 al 10")
		fmt.Scanln(&bossNumber)

		gamePlayers := players_online
		var player1 int32
		var player2 int32
		var play1 int
		var play2 int
		for i := 0; i < len(players_online)/2; i++ {
			player1 = gamePlayers[rand.Intn(len(gamePlayers))].id
			gamePlayers = RemoveElemArray(gamePlayers, indexOf(gamePlayers, player1))
			player2 = gamePlayers[rand.Intn(len(gamePlayers))].id
			gamePlayers = RemoveElemArray(gamePlayers, indexOf(gamePlayers, player2))

			fmt.Println(player1)
			fmt.Println(player2)
			
			indexPlayer1 := indexOfPlayers(player1)
			indexPlayer2 := indexOfPlayers(player2)
			fmt.Println(indexPlayer1)
			fmt.Println(indexPlayer2)


			play1 = int(players_online[indexPlayer1].jugada)
			play2 = int(players_online[indexPlayer2].jugada)
			if math.Abs(float64(bossNumber - play1)) > math.Abs(float64(bossNumber - play2)) {
				RemovePlayer(indexOfPlayers(player1))
			} else if math.Abs(float64(bossNumber - play1)) < math.Abs(float64(bossNumber - play2)) {
				RemovePlayer(indexOfPlayers(player1))
			}
		}

		fmt.Println("Terminamos el juego 3, y se removieron los malos!")
	}
	
	if len(players_online) > 0 {	
		fmt.Println("Presiona enter para obtener los ganadores: ")
		fmt.Scanln(&empezar)

		print_id_player_ganadores(players_online)

	} else {
		log.Fatalln("Oh no, no queda ningun jugador...")
	}
	// Se liberan los procesos, para que puedan preguntar si siguen vivos...
	started = true

	

	fmt.Println("Presiona enter para finalizar el juego: ")
	fmt.Scanln(&empezar)

	// Se termina el programa, luego de que todos hayan preguntado
	// anunciar el ganador
	// Y ver weas del pozo / mandar plata

}

