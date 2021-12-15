package main

import (
	"context"
	"fmt"
	pb "lab3/proto"
	"log"
	"net"

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

func crear_Planeta(nombre_planeta string, ultimo_reloj []int32, ultimo_servidor string, nombre_ciudad string, cant_soldados int32) {
	ciudad := crear_Ciudad(nombre_ciudad, cant_soldados)

	lista_ciudades := [1]Ciudad{ciudad}
	planeta := Planeta{nombre_planeta, lista_ciudades[:], ultimo_reloj[:], ultimo_servidor}

	planetas = append(planetas, planeta)
}

// Para usar cuando se hace un get y el planeta ya existe y MONOLYTIC READS TRUE
func update_Planeta(planeta Planeta, ultimo_reloj []int32, ultimo_servidor string, nombre_ciudad string, cant_soldados int32) {
	esta_ciudad := buscar_ciudad(planeta.lista_ciudades, nombre_ciudad)
	if esta_ciudad > -1 { // Que pasa si la ciudad existe
		planeta.lista_ciudades[esta_ciudad].cant_soldados = cant_soldados
	} else { // que pasa si la ciudad no existe
		ciudad := crear_Ciudad(nombre_ciudad, cant_soldados)
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
//
//

//
//		Monolytic Reads if false hay que pedir merge
//
func monotonic_reads(planeta string, reloj_server []int32) bool {
	var num_servidor int32
	num_planeta := buscar_Planeta(planeta)
	if num_planeta == -1 {
		return true
	}

	servidor := planetas[num_planeta].ultimo_servidor

	if servidor == "10.6.43.110" {
		num_servidor = 0
	} else if servidor == "10.6.43.111" {
		num_servidor = 1
	} else {
		num_servidor = 2
	}
	log.Println("\n"+"reloj planeta: ", planetas[num_planeta].ultimo_reloj[num_servidor], "reloj servidor: ", reloj_server[num_servidor])
	if planetas[num_planeta].ultimo_reloj[num_servidor] > reloj_server[num_servidor] {
		return false
	}
	return true
}

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.To4())
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

//
//		Main Game
//

func main() {

	activo := true
	var accion string
	var planeta string
	var ciudad string
	var cant_soldados int32
	var reloj []int32
	var servidor string

	ip := GetIP()
	if ip == "10.6.43.110" {
		log.Println("Bienvenida Leia Organa, iniciando servicios...")
	} else {
		log.Panicln("Por favor, ejecuta este archivo en la maquina dist@122")
	}

	conn, err := grpc.Dial("10.6.43.109:8080", grpc.WithInsecure()) // Conectamos al IP de 10.6.43.109:8080, el lider.
	if err != nil {
		panic("cannot connect with server " + err.Error())
	}
	serviceClient := pb.NewStarwarsGameClient(conn)

	fmt.Println("Bienvenida princesa Leia")

	for activo {
		fmt.Println("")
		fmt.Print("-> ")
		fmt.Println("Que desea hacer?")
		fmt.Println("1) Preguntar el numero de Rebeldes en una ciudad")
		fmt.Println("2) Cerrar la terminal")
		fmt.Scanln(&accion)
		if accion == "1" {
			fmt.Println("En que planeta queda la ciudad?")
			fmt.Scanln(&planeta)
			fmt.Println("Que ciudad desea buscar?")
			fmt.Scanln(&ciudad)
			res, err := serviceClient.GetCantSoldadosBroker(context.Background(), &pb.GetBrokerRequest{Planeta: planeta, Ciudad: ciudad})
			if err != nil {
				panic("No se pudo hacer el GET  " + err.Error())
			}
			cant_soldados = res.GetRebeldes()
			reloj = res.GetReloj()
			servidor = res.GetServidor()

			log.Println("[PreMono] Reloj: ", reloj)

			monotonic := monotonic_reads(planeta, reloj)
			for monotonic == false {
				log.Println("Entro a Monotonic")

				res, err := serviceClient.MergeLeiaBroker(context.Background(), &pb.MergeLeiaRequest{Planeta: planeta, Ciudad: ciudad})
				if err != nil {
					panic("No se pudo hacer el GET  " + err.Error())
				}
				cant_soldados = res.GetRebeldes()
				reloj = res.GetReloj()
				monotonic = monotonic_reads(planeta, reloj)
				log.Println("[PostMono] Reloj: ", reloj)
			}

			fmt.Println("Cantidad soldados: ", cant_soldados)

			if cant_soldados > -1 {
				fmt.Println("Cantidad soldados: ", cant_soldados)
				if buscar_Planeta(planeta) == -1 {
					crear_Planeta(planeta, reloj, servidor, ciudad, cant_soldados)
					fmt.Println("cree planeta! \n")
				} else {
					update_Planeta(planetas[buscar_Planeta(planeta)], reloj, servidor, ciudad, cant_soldados)
				}
			} else {
				fmt.Println("Planeta no existe")
			}
			fmt.Println()
			fmt.Println(planetas)

		} else if accion == "2" {
			fmt.Println("Adios")
			activo = false
		} else {
			fmt.Println("Escriba una opcion valida")
		}
	}

}
