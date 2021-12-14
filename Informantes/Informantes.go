package main

import (
	"context"
	"fmt"
	pb "lab3/proto"
	"log"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

//
//		Definicion de
//			Variables
//

type PlayerStruct struct {
	id    int32
	alive bool
	round int32
	score int32
	etapa int32
}

type Ciudad struct {
	nombre        string
	cant_soldados int32
}

type Planeta struct {
	planeta         string
	lista_ciudades  []Ciudad
	ultimo_reloj    []int32
	ultimo_servidor string
}

var planetas []Planeta

var minombre string = "Informante"

//
//		CRUD
//		Structs
//

func leer_ciudades(lista_ciudades []Ciudad) {
	for i := 0; i < len(lista_ciudades); i++ {
		ciudad := lista_ciudades[i]
		fmt.Println(i, " Ciudad: \t", ciudad.nombre)
		fmt.Println("Cantidad_soldados: ", ciudad.cant_soldados)
	}
}

func leer_struct(planeta Planeta) {
	fmt.Println("")
	fmt.Println("Struct del planeta:")
	fmt.Println("nombre: \t", planeta.planeta)
	fmt.Println("ultimo_reloj: \t", planeta.ultimo_reloj)
	fmt.Println("ultimo_servidor \t", planeta.ultimo_servidor)
	fmt.Println("")
	leer_ciudades(planeta.lista_ciudades)
	fmt.Println("")
}

/// revisar aca!
func crear_Ciudad(nombre_ciudad string, cant_soldados int32) Ciudad {
	ciudad := Ciudad{nombre_ciudad, cant_soldados}
	return ciudad
}

func crear_Planeta(nombre_planeta string, ultimo_reloj []int32, ultimo_servidor string, nombre_ciudad string, cant_soldados int32) Planeta {
	ciudad := crear_Ciudad(nombre_ciudad, cant_soldados)

	lista_ciudades := [1]Ciudad{ciudad}
	planeta := Planeta{nombre_planeta, lista_ciudades[:], ultimo_reloj[:], ultimo_servidor}

	return planeta
}

// Para usar cuando se hace un get y el planeta ya existe y MONOLYTIC READS TRUE
func update_Planeta(planeta Planeta, ultimo_reloj []int32, ultimo_servidor string, nombre_ciudad string, nuevo_nombre_ciudad string, cant_soldados int32) {
	esta_ciudad := buscar_ciudad(planeta.lista_ciudades, nombre_ciudad)
	if esta_ciudad > -1 { // Que pasa si la ciudad existe
		if cant_soldados > -1 {
			planeta.lista_ciudades[esta_ciudad].cant_soldados = cant_soldados
		}
		planeta.lista_ciudades[esta_ciudad].nombre = nuevo_nombre_ciudad
	} else { // que pasa si la ciudad no existe
		if cant_soldados == -1 {
			cant_soldados = planeta.lista_ciudades[esta_ciudad].cant_soldados
		}
		ciudad := crear_Ciudad(nuevo_nombre_ciudad, cant_soldados)
		planeta.lista_ciudades = append(planeta.lista_ciudades, ciudad)
	}

	planeta.ultimo_reloj = ultimo_reloj
	planeta.ultimo_servidor = ultimo_servidor

}

//
//		Busqueda
//		structs
//

func buscar_Planeta(nombre_buscado string) int32 {
	var planeta Planeta
	for i := 0; i < len(planetas); i++ {
		planeta = planetas[i]
		if nombre_buscado == planeta.planeta {
			return int32(i)
		}
	}
	return -1
}

func buscar_ciudad(lista_ciudades []Ciudad, nombre_buscado string) int32 {
	var ciudad Ciudad
	for i := 0; i < len(lista_ciudades); i++ {
		ciudad = lista_ciudades[i]
		if nombre_buscado == ciudad.nombre {
			return int32(i)
		}
	}
	return -1
}

//
//		Read you Writes if false hay que pedir merge
//
func read_your_write(planeta string, reloj_server []int32) bool {
	var num_servidor int32
	num_planeta := buscar_Planeta(planeta)
	servidor := planetas[num_planeta].ultimo_servidor

	if servidor == "10.6.43.110" {
		num_servidor = 0
	} else if servidor == "10.6.43.111" {
		num_servidor = 1
	} else {
		num_servidor = 2
	}
	if planetas[num_planeta].ultimo_reloj[num_servidor] > reloj_server[num_servidor] {
		return false
	}
	return true
}

//
//		gRPC
//

func ConectarServidores(comando string) (reloj []int32, servidor string) {
	conn, err := grpc.Dial("10.6.43.109:8080", grpc.WithInsecure()) // Conexion con el Broker
	if err != nil {
		panic("cannot connect with server " + err.Error())
	}
	serviceClient := pb.NewStarwarsGameClient(conn)

	res, err := serviceClient.AskForServers(context.Background(), &pb.AskForServersRequest{Comando: comando})

	servidor_a_conectar := res.GetServidor()
	log.Println("El Broker me indico conectar con:", servidor_a_conectar)

	conn2, err := grpc.Dial(servidor_a_conectar+":8081", grpc.WithInsecure()) // Conexion con el Servidor, cambiar ip por respuesta de servidor
	if err != nil {
		panic("cannot connect with server " + err.Error())
	}
	serviceClient2 := pb.NewStarwarsGameClient(conn2)

	res2, err := serviceClient2.AskedServer(context.Background(), &pb.AskedServerRequest{Informante: minombre, Comando: comando})

	ultimo_reloj := res2.GetReloj()

	// Revisar el comando y hacer cambios necesarios
	splitted_comando := strings.Split(comando, " ")

	// read := read_your_write(planeta, reloj)
	// if read == false {
	// 	res, err := serviceClient.MargeInformante(context.Background(), &pb.GetRequest{planeta: planeta, ciudad: ciudad, cant_soldados: cant_soldados})
	// 	if err != nil {
	// 		panic("No se pudo hacer el GET  " + err.Error())
	// 	}
	// 	reloj = res.GetReloj()
	// 	servidor = res.GetServidor()
	// }

	if splitted_comando[0] == "AddCity" {
		log.Println("[AddCity]")
		nombre_planeta := splitted_comando[1]
		nombre_ciudad := splitted_comando[2]
		string_soldados := splitted_comando[3]
		int_soldados, _ := strconv.Atoi(string_soldados)
		cant_soldados := int32(int_soldados)

		num_planeta := buscar_Planeta(nombre_planeta)
		if num_planeta == -1 {
			planeta := crear_Planeta(nombre_planeta, ultimo_reloj, servidor_a_conectar, nombre_ciudad, cant_soldados)
			planetas = append(planetas, planeta)
		} else {
			ciudad := crear_Ciudad(nombre_ciudad, cant_soldados)
			planetas[num_planeta].lista_ciudades = append(planetas[num_planeta].lista_ciudades, ciudad)
			planetas[num_planeta].ultimo_reloj = ultimo_reloj
			planetas[num_planeta].ultimo_servidor = servidor_a_conectar
		}

	} else if splitted_comando[0] == "UpdateName" {
		log.Println("[UpdateName]")
		nombre_planeta := splitted_comando[1]
		nombre_ciudad := splitted_comando[2]
		nuevo_nombre_ciudad := splitted_comando[3]

		num_planeta := buscar_Planeta(nombre_planeta)
		planeta := planetas[num_planeta]
		update_Planeta(planeta, ultimo_reloj, servidor_a_conectar, nombre_ciudad, nuevo_nombre_ciudad, -1)

	} else if splitted_comando[0] == "UpdateNumber" {
		log.Println("[UpdateNumber]")
		nombre_planeta := splitted_comando[1]
		nombre_ciudad := splitted_comando[2]
		string_soldados := splitted_comando[3]
		int_soldados, _ := strconv.Atoi(string_soldados)
		cant_soldados := int32(int_soldados)

		num_planeta := buscar_Planeta(nombre_planeta)
		planeta := planetas[num_planeta]
		update_Planeta(planeta, ultimo_reloj, servidor_a_conectar, nombre_ciudad, nombre_ciudad, cant_soldados)

	} else if splitted_comando[0] == "DeleteCity" {
		log.Println("[DeleteCity]")
		nombre_planeta := splitted_comando[1]
		nombre_ciudad := splitted_comando[2]

		num_planeta := buscar_Planeta(nombre_planeta)
		num_ciudad := buscar_ciudad(planetas[num_planeta].lista_ciudades, nombre_ciudad)
		planetas[num_planeta].lista_ciudades = append(planetas[num_planeta].lista_ciudades[:num_ciudad], planetas[num_planeta].lista_ciudades[num_ciudad+1:]...)

	}

	// aqui

	log.Println("El servidor me envio el reloj ", ultimo_reloj)

	return ultimo_reloj, servidor_a_conectar

}

//
//		Main Game
//

func main() {

	activo := true
	var accion string
	var planeta string
	var ciudad string
	var nueva_ciudad string
	var cant_soldados string

	var comando string

	fmt.Println("Bienvenida Informante")
	fmt.Print("-> ")
	for activo {
		fmt.Println("¿Que desea hacer?")
		fmt.Println("1) Agregar una ciudad")
		fmt.Println("2) Actualizar el nombre de una ciudad")
		fmt.Println("3) Actualizar el numero de rebeldes de una ciudad")
		fmt.Println("4) Eliminar una ciudad")
		fmt.Println("5) Cerrar la terminal")
		fmt.Print("> ")
		fmt.Scanln(&accion)
		fmt.Println("")
		if accion == "1" {
			fmt.Println("¿En que planeta queda la ciudad?")
			fmt.Print("> ")
			fmt.Scanln(&planeta)
			fmt.Println("¿Cual es el nombre de la ciudad que desea agregar?")
			fmt.Print("> ")
			fmt.Scanln(&ciudad)
			fmt.Println("¿Cuantos soldados tiene esta ciudad?")
			fmt.Print("> ")
			fmt.Scanln(&cant_soldados)

			if cant_soldados == "" {
				cant_soldados = "0"
			}

			comando = "AddCity " + planeta + " " + ciudad + " " + cant_soldados

			reloj, servidor := ConectarServidores(comando)
			fmt.Println("")
			fmt.Println("Reloj: \t", reloj, " Servidor: \t", servidor)

		} else if accion == "2" {
			fmt.Println("¿En que planeta queda la ciudad?")
			fmt.Print("> ")
			fmt.Scanln(&planeta)
			fmt.Println("¿Cual es el nombre actual de la ciudad?")
			fmt.Print("> ")
			fmt.Scanln(&ciudad)
			fmt.Println("¿Cual es el nuevo nombre de la ciudad?")
			fmt.Print("> ")
			fmt.Scanln(&nueva_ciudad)

			comando = "UpdateName " + planeta + " " + ciudad + " " + nueva_ciudad

			reloj, servidor := ConectarServidores(comando)
			fmt.Println("")
			fmt.Println("Reloj: \t", reloj, " Servidor: \t", servidor)

		} else if accion == "3" {
			fmt.Println("¿En que planeta queda la ciudad?")
			fmt.Print("> ")
			fmt.Scanln(&planeta)
			fmt.Println("¿Cual es el nombre de la ciudad?")
			fmt.Print("> ")
			fmt.Scanln(&ciudad)
			fmt.Println("¿Cual es el nuevo numero de rebeldes?")
			fmt.Print("> ")
			fmt.Scanln(&cant_soldados)

			if cant_soldados == "" {
				cant_soldados = "0"
			}

			comando = "UpdateNumber " + planeta + " " + ciudad + " " + cant_soldados

			reloj, servidor := ConectarServidores(comando)
			fmt.Println("")
			fmt.Println("Reloj: \t", reloj, " Servidor: \t", servidor)

		} else if accion == "4" {
			fmt.Println("¿En que planeta queda la ciudad?")
			fmt.Print("> ")
			fmt.Scanln(&planeta)
			fmt.Println("¿Cual es el nombre de la ciudad?")
			fmt.Print("> ")
			fmt.Scanln(&ciudad)

			comando = "DeleteCity " + planeta + " " + ciudad

			reloj, servidor := ConectarServidores(comando)
			fmt.Println("")
			fmt.Println("Reloj: \t", reloj, " Servidor: \t", servidor)

		} else if accion == "5" {
			fmt.Println("Adios")
			activo = false
		} else {
			fmt.Println("Escriba una opción valida")
		}
		fmt.Println("")
	}

}
