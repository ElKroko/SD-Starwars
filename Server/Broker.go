package main

import (
	"context"
	"fmt"
	pb "lab3/proto"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
)

var lista_servidores = [3]string{"1", "2", "3"}
var cant_solicitudes = 0

type server struct {
	pb.UnimplementedStarwarsGameServer
}

//
//		Funciones de GRPC
//

func (s *server) GetCantSoldadosBroker(ctx context.Context, in *pb.GetBrokerRequest) (*pb.GetBrokerReply, error) {
	log.Printf("Estan haciendo request!")
	cant_solicitudes += 1
	planeta := in.GetPlaneta()
	ciudad := in.GetCiudad()

	log.Printf("Planeta: %s \t Ciudad: %s", planeta, ciudad)

	servidor := getRandomServer()
	log.Println("Se ha elegido el servidor: ", servidor)
	return &pb.GetBrokerReply{Servidor: servidor}, nil

}

func (s *server) MergeLeia(ctx context.Context, in *pb.MergeLeiaRequest) (*pb.MergeLeiaReply, error) {
	log.Printf("Nos piden un merge!")
	cant_solicitudes += 1
	planeta := in.GetPlaneta()
	ciudad := in.GetCiudad()

	log.Printf("Planeta: %s \t Ciudad: %s", planeta, ciudad)

	var cant_rebeldes = 1
	var reloj []int32
	reloj[0] = 0
	reloj[1] = 1
	reloj[2] = 2
	servidor := getRandomServer()

	return &pb.MergeLeiaReply{Rebeldes: int32(cant_rebeldes), Reloj: reloj, Servidor: servidor}, nil

}

//
//		Funciones Auxiliares
//

func getRandomServer() string {
	numserver := rand.Intn(3)
	servidor := lista_servidores[numserver]
	return servidor
}
func main() {
	fmt.Println()
	log.Printf("Bienvenido a Broker Mos Eisley, iniciando servicios...")
	rand.Seed(time.Now().UnixNano())

	go func() {
		listner, err := net.Listen("tcp", ":8080")

		if err != nil {
			panic("cannot create tcp connection " + err.Error())
		}

		serv := grpc.NewServer()
		pb.RegisterStarwarsGameServer(serv, &server{})
		if err = serv.Serve(listner); err != nil {
			panic("cannot initialize the server" + err.Error())
		}
	}()

	log.Print("Servicios iniciados, escuchando red...")

	flag_opcion := true
	var opcion string
	for flag_opcion {
		fmt.Println("Que deseas hacer?")
		fmt.Println("\t [1] Preguntar cantidad de solicitudes")
		fmt.Println("\t [2] Esperar 10 segundos")
		fmt.Println("\t [3] Cerrar el programa")
		fmt.Print("> ")
		fmt.Scanln(&opcion)
		if opcion == "1" {
			fmt.Println("La cantidad recibida actualmente es de: ", cant_solicitudes, " solicitudes")
			time.Sleep(5 * time.Second)
		} else if opcion == "2" {
			fmt.Println("Esperando 10 segundos...")
			time.Sleep(10 * time.Second)
		} else if opcion == "3" {
			fmt.Println("Cerrando servicios...")
			flag_opcion = false
		} else {
			fmt.Println("Por favor, escribe una opcion correcta!")
		}
		fmt.Println("")

	}

	log.Println("Se ha cerrado el proceso.")
}
