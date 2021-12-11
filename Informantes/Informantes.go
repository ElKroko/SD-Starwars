package main

import (
	"context"
	"fmt"
	pb "lab/proto"
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

func crear_Planeta(nombre_planeta string, ultimo_reloj []int32, ultimo_servidor string, nombre_ciudad string, cant_soldados int32) Planeta {
	ciudad := crear_Ciudad(nombre_ciudad, cant_soldados)

	lista_ciudades := [1]Ciudad{ciudad}
	planeta := Planeta{nombre_planeta, lista_ciudades[:], ultimo_reloj[:], ultimo_servidor}

	return planeta
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
//		Main Game
//

func main() {

	activo := true
	var accion string
	var planeta string
	var ciudad string
	var nueva_ciudad string
	var cant_soldados int32
	var reloj []int32
	var servidor string

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure()) // Conectamos al IP de 10.6.43.109:8080, el lider.
	if err != nil {
		panic("cannot connect with server " + err.Error())
	}
	serviceClient := pb.StarwarsGameClient(conn)

	fmt.Println("Bienvenida Informante")
	fmt.Print("-> ")
	for activo {
		fmt.Println("¿Que desea hacer?")
		fmt.Print("1) Agregar una ciudad")
		fmt.Print("2) Actualizar el nombre de una ciudad")
		fmt.Print("3) Actualizar el numero de rebeldes de una ciudad")
		fmt.Print("4) Eliminar una ciudad")
		fmt.Print("5) Cerrar la terminal")
		fmt.Scanln(&accion)
		if int(accion) == "1" {
			fmt.Println("¿Cual es el nombre de la ciudad que desea agregar?")
			fmt.Scanln(&ciudad)
			fmt.Println("¿En que planeta queda la ciudad?")
			fmt.Scanln(&planeta)
			res, err := serviceClient.CreatePlanet(context.Background(), &pb.GetRequest{planeta, ciudad})
			if err != nil {
				panic("No se pudo obtener IP  " + err.Error())
			}
			ip = res.GetIP()
			res, err := serviceClient.CreatePlanet(context.Background(), &pb.GetRequest{planeta, ciudad})
			if err != nil {
				panic("No se pudo hacer el UPDATE  " + err.Error())
			}
			reloj = res.GetReloj()
			servidor = res.GetServidor()

			read := read_your_write(planeta, reloj)
			if read == false {
				res, err := serviceClient.GetMargeInformante(context.Background(), &pb.GetRequest{planeta: planeta, ciudad: ciudad})
				if err != nil {
					panic("No se pudo hacer el GET  " + err.Error())
				}
				reloj = res.GetReloj()
				servidor = res.GetServidor()
			}
			update_Planeta(planeta, reloj, servidor, ciudad, cant_soldados)
		} else if int(accion) == "2" {
			fmt.Println("¿Cual es el nombre actual de la ciudad?")
			fmt.Scanln(&ciudad)
			fmt.Println("¿Cual es el nuevo nombre de la ciudad?")
			fmt.Scanln(&nueva_ciudad)
			fmt.Println("¿En que planeta queda la ciudad?")
			fmt.Scanln(&planeta)
			res, err := serviceClient.UpdatePlanet(context.Background(), &pb.GetRequest{planeta, ciudad, nueva_ciudad})
			if err != nil {
				panic("No se pudo obtener IP  " + err.Error())
			}
			ip = res.GetIP()
			res, err = serviceClient.UpdatePlanet(context.Background(), &pb.GetRequest{planeta, ciudad, nueva_ciudad})
			if err != nil {
				panic("No se pudo obtener actualizar la ciudad  " + err.Error())
			}
			reloj = res.GetReloj()
			servidor = res.GetServidor()

			read := read_your_write(planeta, reloj)
			if read == false {
				res, err := serviceClient.MargeInformante(context.Background(), &pb.GetRequest{planeta: planeta, ciudad: ciudad, nueva_ciudad: nueva_ciudad})
				if err != nil {
					panic("No se pudo hacer el GET  " + err.Error())
				}
				reloj = res.GetReloj()
				servidor = res.GetServidor()
			}
			update_Planeta(planeta, reloj, servidor, ciudad, cant_soldados)
		} else if int(accion) == "3" {
			fmt.Println("¿Cual es el nombre de la ciudad?")
			fmt.Scanln(&ciudad)
			fmt.Println("¿Cual es el nuevo numero de rebeldes?")
			fmt.Scanln(&cant_soldados)
			fmt.Println("¿En que planeta queda la ciudad?")
			fmt.Scanln(&planeta)
			res, err := serviceClient.UpdateSoldiers(context.Background(), &pb.GetRequest{planeta, ciudad, cant_soldados})
			if err != nil {
				panic("No se pudo obtener IP  " + err.Error())
			}
			ip = res.GetIP()
			res, err = serviceClient.UpdateSoldiers(context.Background(), &pb.GetRequest{planeta, ciudad, cant_soldados})
			if err != nil {
				panic("No se pudo obtener actualizar los soldados  " + err.Error())
			}
			reloj = res.GetReloj()
			servidor = res.GetServidor()

			read := read_your_write(planeta, reloj)
			if read == false {
				res, err := serviceClient.MargeInformante(context.Background(), &pb.GetRequest{planeta: planeta, ciudad: ciudad, cant_soldados: cant_soldados})
				if err != nil {
					panic("No se pudo hacer el GET  " + err.Error())
				}
				reloj = res.GetReloj()
				servidor = res.GetServidor()
			}
			update_Planeta(planeta, reloj, servidor, ciudad, cant_soldados)
		} else if int(accion) == "4" {
			fmt.Println("¿Cual es el nombre de la ciudad?")
			fmt.Scanln(&ciudad)
			fmt.Println("¿En que planeta queda la ciudad?")
			fmt.Scanln(&planeta)
			res, err := serviceClient.DeletePlanet(context.Background(), &pb.GetRequest{planeta, ciudad})
			if err != nil {
				panic("No se pudo obtener IP  " + err.Error())
			}
			ip = res.GetIP()
			res, err = serviceClient.DeletePlanet(context.Background(), &pb.GetRequest{planeta, ciudad})
			if err != nil {
				panic("No se pudo eliminar el planeta " + err.Error())
			}
			reloj = res.GetReloj()
			servidor = res.GetServidor()

			read := read_your_write(planeta, reloj)
			if read == false {
				res, err := serviceClient.MargeInformante(context.Background(), &pb.GetRequest{planeta: planeta, ciudad: ciudad})
				if err != nil {
					panic("No se pudo hacer el GET  " + err.Error())
				}
				reloj = res.GetReloj()
				servidor = res.GetServidor()
			}
			update_Planeta(planeta, reloj, servidor, ciudad, cant_soldados)
		} else if int(accion) == "5" {
			fmt.Println("Adios")
			activo = false
		} else {
			fmt.Println("Escriba una opción valida")
		}
	}

}
