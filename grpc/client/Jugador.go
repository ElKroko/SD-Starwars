package main

import (
	"context"
	"fmt"
	pb "lab/game/proto"
	"log"
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type PlayerStruct struct {
	id    int32
	alive bool
	round int32
	score int32
	etapa int32
}


func print_player(jugador_struct PlayerStruct) {
	fmt.Println("")
	fmt.Println("Struct del Jugador: ")
	fmt.Println("Id:", strconv.Itoa(int(jugador_struct.id)))
	fmt.Println("Alive:", strconv.FormatBool(jugador_struct.alive))
	fmt.Println("Round:", strconv.Itoa(int(jugador_struct.round)))
	fmt.Println("Score:", strconv.Itoa(int(jugador_struct.score)))
	fmt.Println("Etapa:", strconv.Itoa(int(jugador_struct.etapa)))
}

func main() {

	// Decidir si sera Bot o Jugador
	var jugador bool

	fmt.Println("Es usted un jugador? [1] Si [0] no")
	fmt.Print("-> ")
	fmt.Scanln(&jugador)


	// Definicion de Variables

	iniciado := false


	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure()) // Conectamos al IP de 10.6.43.109:8080, el lider.

	if err != nil {
		panic("cannot connect with server " + err.Error())
	}

	serviceClient := pb.NewCalamardoGameClient(conn)

	// Inscripcion del jugador

	message := "Hola, soy jugador"

	res, err := serviceClient.JoinGame(context.Background(), &pb.JoinRequest{Message: message})

	if err != nil {
		panic("No se pudo anadir el jugador  " + err.Error())
	}

	jugador_struct := PlayerStruct{res.GetIdJugador(), res.GetAlive(), res.GetRound(), 0, 0}

	iniciado = true

	print_player(jugador_struct)



	if iniciado {
		fmt.Println("Bienvenidx al juego del Calamardo!")
	}

	// Esperar a los jugadores a que se conecten

	flag1 := true

	for flag1 {


		res1, err := serviceClient.StartGame(context.Background(), &pb.StartRequest{Id: jugador_struct.id, Message: "Puedo jugar?", Etapa: jugador_struct.etapa})
		if err != nil {
			panic("No pudimos chequear si esta inciado  " + err.Error())
		}

		if res1.GetStarted(){
			flag1 = false
			fmt.Println("Preparado? empezamos!")
			fmt.Println("")
		} else {
			fmt.Println("Espera, aun no hay suficientes jugadores...")
			time.Sleep(2 * time.Second)
		}
	}
	

	// comienza juego 1

	fmt.Println("Iniciamos el Juego 1")
	rand.Seed(time.Now().UnixNano())

	var numero int
	ronda := jugador_struct.round
	flag1 = true


	for flag1 {
		res1, err := serviceClient.StartGame(context.Background(), &pb.StartRequest{Id: jugador_struct.id, Message: "Puedo jugar?", Etapa: 1})
		if err != nil {
			panic("No pudimos chequear si esta inciado  " + err.Error())
		}

		if res1.GetStarted(){
			flag1 = false
			fmt.Println("Preparado? empezamos!")
			fmt.Println("")
		} else {
			fmt.Println("Espera, aun no inicia el juego 1...")
			time.Sleep(2 * time.Second)
		}
	}

	// ==============================================================================
	// 									Inicio Juego 1
	// ==============================================================================


	jugador_struct.etapa = 1

	fmt.Println("--------------------------------")
	fmt.Println("       Inicio del juego 1       ")
	fmt.Println("--------------------------------")
	fmt.Println("")

	for ronda < 4 && jugador_struct.etapa == 1 {
		// Seleccionar numero para jugar
		if jugador {
			fmt.Println("Juego 1: Elija un numero del 1 al 10")
			fmt.Print("-> ")
			fmt.Scanln(&numero)

		} else {
			numero = rand.Intn(10)

		}
		fmt.Println("La jugada, de la ronda ", ronda , "fue ", numero)

		// Chequear si el numero es igual o mayor al del lider

		res, err := serviceClient.JuegoMsg(context.Background(), &pb.JuegoRequest{
			Id: jugador_struct.id,
			Jugada:int32(numero),  
			Round: jugador_struct.round,
			Score: jugador_struct.score,
			Etapa: jugador_struct.etapa})

		if err != nil {
			panic("No pudimos mandar la jugada del juego 1  " + err.Error())
		}

		jugador_struct.alive = res.GetAlive()
		jugador_struct.round = res.GetRound()
		jugador_struct.score = res.GetScore()
		jugador_struct.etapa = res.GetEtapa()
		
		print_player(jugador_struct)
		
		if !jugador_struct.alive {
			log.Fatalln("Oh no! te han matado en el juego 1")
		}

		if jugador_struct.etapa == 2{
			fmt.Println("")
			log.Println("Pasaste a la siguiente etapa!")
		}
		// Fin

		
		jugador_struct.round = ronda + 1
		ronda = jugador_struct.round
	}


	// Se calcula el numero a mandar
	// Se manda el numero, y se espera respuesta

	if !jugador_struct.alive {
		log.Fatalln("Oh no! te han matado al final del juego 1")
	}
	
	// ==============================================================================
	// 									Inicio Juego 2
	// ==============================================================================

	jugador_struct.etapa = 2
	fmt.Println("")
	fmt.Println("--------------------------------")
	fmt.Println("       Inicio del juego 2       ")
	fmt.Println("--------------------------------")
	fmt.Println("")


	// Esperamos que el juego 2 este listo
	flag1 = true


	for flag1 {
		res1, err := serviceClient.StartGame(context.Background(), &pb.StartRequest{Id: jugador_struct.id, Message: "Puedo jugar?", Etapa: 1})
		if err != nil {
			panic("No pudimos chequear si esta inciado  " + err.Error())
		}

		if res1.GetStarted(){
			flag1 = false
			fmt.Println("Preparado? empezamos!")
			fmt.Println("")
		} else {
			fmt.Println("Espera, aun no inicia el juego 2...")
			time.Sleep(2 * time.Second)
		}
	}

	// Generar un numero random entre 1 y 4
	if jugador {
		fmt.Println("Juego 2: Elija un numero del 1 al 4")
		fmt.Print("-> ")
		fmt.Scanln(&numero)
	} else {
		numero = rand.Intn(4)
	}


	// Enviar el numero al lider
	res1, err := serviceClient.JuegoMsg(context.Background(), &pb.JuegoRequest{
		Id: jugador_struct.id,
		Jugada:int32(numero),  
		Round: jugador_struct.round,
		Score: jugador_struct.score,
		Etapa: jugador_struct.etapa})

	if err != nil {
		panic("No pudimos mandar la jugada del juego 2  " + err.Error())
	}

	jugador_struct.alive = res1.GetAlive()
	jugador_struct.round = res1.GetRound()
	jugador_struct.score = res1.GetScore()
	jugador_struct.etapa = res1.GetEtapa()
	
	print_player(jugador_struct)


	// Esperar en un for la respuesta del lider sobre este numero. cuando me envie true
	// y el estado del aliado nuevo, decidire si seguir o exit()
	flag1 = true
	for flag1 {
		res1, err := serviceClient.StartGame(context.Background(), &pb.StartRequest{Id: jugador_struct.id, Message: "Puedo jugar?", Etapa: 1})
		if err != nil {
			panic("No pudimos chequear si esta inciado  " + err.Error())
		}

		if res1.GetStarted(){
			flag1 = false
			fmt.Println("Preparado? empezamos!")
			fmt.Println("")
		} else {
			fmt.Println("Espera, aun no inicia el juego 3...")
			time.Sleep(2 * time.Second)
		}
	}

	// Revisar si estoy vivo!
	
	resMuerte, err := serviceClient.Muerte(context.Background(), &pb.MuerteRequest{Id: jugador_struct.id})
	if err != nil {
		panic("No pudimos chequear si esta inciado  " + err.Error())
	}
	
	jugador_struct.alive = resMuerte.GetAlive()
	if !jugador_struct.alive {
		log.Fatalln("Oh no! te han matado al final del juego 2")
	}



	// Luego de revisar si estoy vivo, comienza el juego 3


	// ==============================================================================
	// 									Inicio Juego 3
	// ==============================================================================

	fmt.Println("--------------------------------")
	fmt.Println("       Inicio del juego 3       ")
	fmt.Println("--------------------------------")
	fmt.Println("")

	// Generar un numero random entre 1 y 10
	if jugador {
		fmt.Println("Juego 3: Elija un numero del 1 al 10")
		fmt.Print("-> ")
		fmt.Scanln(&numero)
	} else {
		numero = rand.Intn(10)
	}


	// Enviar el numero al lider
	res2, err := serviceClient.JuegoMsg(context.Background(), &pb.JuegoRequest{
		Id: jugador_struct.id,
		Jugada:int32(numero),  
		Round: jugador_struct.round,
		Score: jugador_struct.score,
		Etapa: jugador_struct.etapa})

	if err != nil {
		panic("No pudimos mandar la jugada del juego 2  " + err.Error())
	}

	jugador_struct.alive = res2.GetAlive()
	jugador_struct.round = res2.GetRound()
	jugador_struct.score = res2.GetScore()
	jugador_struct.etapa = res2.GetEtapa()
	
	print_player(jugador_struct)


	// Esperar en un for la respuesta del lider sobre este numero. cuando me envie true
	// y el estado del aliado nuevo, decidire si seguir o exit()
	flag1 = true
	for flag1 {
		res1, err := serviceClient.StartGame(context.Background(), &pb.StartRequest{Id: jugador_struct.id, Message: "Puedo jugar?", Etapa: 1})
		if err != nil {
			panic("No pudimos chequear si esta inciado  " + err.Error())
		}

		if res1.GetStarted(){
			flag1 = false
			fmt.Println("Preparado? empezamos!")
			fmt.Println("")
		} else {
			fmt.Println("Espera, estamos verificando los resultados...")
			time.Sleep(2 * time.Second)
		}
	}

	// Revisar si estoy vivo!
	resMuerte, err = serviceClient.Muerte(context.Background(), &pb.MuerteRequest{Id: jugador_struct.id})
	if err != nil {
		panic("No pudimos chequear si haz muerto... probablemente si!  " + err.Error())
	}
	
	jugador_struct.alive = resMuerte.GetAlive()
	if !jugador_struct.alive {
		log.Fatalln("Oh no! te han matado al final del juego 3")
	}

}	
