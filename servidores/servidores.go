package main

import (
	"bufio"
	"context"
	"fmt"
	pb "lab3/proto"
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
	pb.UnimplementedStarwarsGameServer
}

type Planeta struct {
	nombre_planeta string
	reloj          []int32
}

var planetas []Planeta

func buscar_Planeta(nombre_buscado string) int32 {
	var planeta Planeta
	for i := 0; i < len(planetas); i++ {
		planeta = planetas[i]
		if nombre_buscado == planeta.nombre_planeta {
			return int32(i)
		}
	}
	return -1
}

func escribir_archivo(nombre_archivo string, texto string) {

	f, err := os.OpenFile(nombre_archivo+".txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(texto); err != nil {
		panic(err)
	}

}

func add_log(nombre_planeta string, nombre_ciudad string, cant_soldados int32, nombre_accion string) {
	escribir_archivo("log_"+nombre_planeta, nombre_accion+" "+nombre_planeta+" "+nombre_ciudad+" "+fmt.Sprint(cant_soldados))
}

func log_string() string {
	var nombre_planeta string
	var planeta Planeta
	nuevo_texto := ""
	for i := 0; i < len(planetas); i++ {
		planeta = planetas[i]
		nombre_planeta = planeta.nombre_planeta

		f, err := os.Open("log_" + nombre_planeta + ".txt")
		if err != nil {
			return ""
		}
		defer f.Close()

		// Splits on newlines by default.
		scanner := bufio.NewScanner(f)

		line := ""

		for scanner.Scan() {
			line = scanner.Text()
			nuevo_texto = nuevo_texto + line + "\n"
		}
	}

	return nuevo_texto
}

func crear_planeta(nombre_planeta string) {

	_, err := os.Create(nombre_planeta + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	_, err2 := os.Create("log_" + nombre_planeta)

	if err2 != nil {
		log.Fatal(err2)
	}
	var nuevo_planeta Planeta

	nuevo_planeta.nombre_planeta = nombre_planeta
	nuevo_planeta.reloj = []int32{0, 0, 0}

	planetas = append(planetas, nuevo_planeta)
}

func actualizar_reloj(nombre_planeta string, reloj []int32) {
	for i := 0; i < len(planetas); i++ {
		if planetas[i].nombre_planeta == nombre_planeta {
			planetas[i].reloj = reloj
		}
	}
}

func reloj_string() string {

}

func existe_ciudad(nombre_planeta string, nombre_ciudad string) bool {
	f, err := os.Open(nombre_planeta + ".txt")
	if err != nil {
		return false
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), nombre_ciudad) {
			return true
		}

		line++
	}

	return false
}

func existe_planeta(nombre_planeta string) bool {
	if _, err := os.Stat(nombre_planeta + ".txt"); err == nil {
		return true
	}
	return false
}

func crear_ciudad(nombre_planeta string, nombre_ciudad string, cant_soldados int32) {

	if existe_planeta(nombre_planeta) {
		if existe_ciudad(nombre_planeta, nombre_ciudad) {
			actualizar_soldados_ciudad(nombre_planeta, nombre_ciudad, cant_soldados)
			add_log(nombre_planeta, nombre_ciudad, cant_soldados, "UpdateNumber")
		} else {
			escribir_archivo(nombre_planeta, nombre_planeta+" "+nombre_ciudad+" "+fmt.Sprint(cant_soldados))
			add_log(nombre_planeta, nombre_ciudad, cant_soldados, "AddCity")
		}

	} else {
		crear_planeta(nombre_planeta)
		escribir_archivo(nombre_planeta, nombre_planeta+" "+nombre_ciudad+" "+fmt.Sprint(cant_soldados))
		add_log(nombre_planeta, nombre_ciudad, cant_soldados, "AddCity")
	}
}

func actualizar_nombre_ciudad(nombre_planeta string, nombre_ciudad string, nuevo_nombre_ciudad string) {

	if existe_planeta(nombre_planeta) {
		if existe_ciudad(nombre_planeta, nombre_ciudad) {
			eliminar_ciudad(nombre_planeta, nombre_ciudad)
		}
		crear_ciudad(nombre_planeta, nuevo_nombre_ciudad, 0)
	} else {
		crear_planeta(nombre_planeta)
		crear_ciudad(nombre_planeta, nuevo_nombre_ciudad, 0)
	}

}

func actualizar_soldados_ciudad(nombre_planeta string, nombre_ciudad string, cant_soldados int32) {

	if existe_planeta(nombre_planeta) {
		if existe_ciudad(nombre_planeta, nombre_ciudad) {
			eliminar_ciudad(nombre_planeta, nombre_ciudad)
		}
		crear_ciudad(nombre_planeta, nombre_ciudad, cant_soldados)

	} else {
		crear_planeta(nombre_planeta)
		crear_ciudad(nombre_planeta, nombre_ciudad, cant_soldados)
	}

}

func eliminar_ciudad(nombre_planeta string, nombre_ciudad string) bool {

	f, err := os.Open(nombre_planeta + ".txt")
	if err != nil {
		return false
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	nuevo_texto := ""
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), nombre_ciudad) {

		} else {
			nuevo_texto = nuevo_texto + scanner.Text() + "\n"
		}
	}
	e := os.Remove(nombre_planeta + ".txt")
	if e != nil {
		log.Fatal(e)
	}
	crear_planeta(nombre_planeta)
	escribir_archivo(nombre_planeta, nuevo_texto)
	return true

}

func obtener_rebeldes(nombre_planeta string, nombre_ciudad string) int32 {
	f, err := os.Open(nombre_planeta + ".txt")
	if err != nil {
		return -1
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), nombre_ciudad) {
			line := strings.Split(scanner.Text(), " ")
			cant_soldados, _ := strconv.Atoi(line[2])
			return int32(cant_soldados)
		}
	}

	return -1
}

func merge(log_recibido string) {
	lineas := strings.Split(log_recibido, "\n")
	var nombre_planeta string
	for i := 0; i < len(lineas); i++ {
		nombre_planeta = ""
		f, err := os.Open(nombre_planeta + ".txt")
		if err != nil {
			//ver si retornar algo
			log.Println("Error en el merge!")
		}
		defer f.Close()

		// Splits on newlines by default.
		scanner := bufio.NewScanner(f)

		line := ""
		nuevo_texto := ""
		var nombre_ciudad string
		for scanner.Scan() {
			line = scanner.Text()
			nombre_ciudad = strings.Split(line, " ")[1]
			nuevo_texto = nuevo_texto + line + "\n" + nombre_ciudad //sacar nombre_ciudad
		}
	}

}

func clean_logs() {
	for i := 0; i < len(planetas); i++ {
		e := os.Remove(planetas[i].nombre_planeta + ".txt")
		if e != nil {
			log.Fatal(e)
		}

		f, err := os.Create(planetas[i].nombre_planeta + ".txt")

		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
}

// func merge_todo() {
// 	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure()) // Conectamos al IP de 10.6.43.109:8080, el lider.
// 	if err != nil {
// 		panic("cannot connect with server " + err.Error())
// 	}
// 	serviceClient := pb.NewStarwarsGameClient(conn)
// 	res, err := serviceClient.GetLogs(context.Background(), &pb.GetLogs{})
// 	if err != nil {
// 		panic("No se pudo hacer conexion de merge  " + err.Error())
// 	}
// 	log2 := res.GetLog()
// 	res2, err2 := serviceClient.GetLogs(context.Background(), &pb.GetLogs{})
// 	if err2 != nil {
// 		panic("No se pudo hacer conexion de merge  " + err.Error())
// 	}
// 	log3 := res2.GetLog()
// 	// res, err := serviceClient.PostReloj(context.Background(), &pb.PostReloj{})
// 	// if err != nil {
// 	// 	panic("No se pudo hacer conexion de merge  " + err.Error())
// 	// }
// 	fmt.Println("Merge realizado")
// }

func actualizar_merge_archivos(data string) {
	lineas := strings.Split(data, "\n")
	var linea []string
	var planeta string
	planeta_actual := strings.Split(lineas[0], " ")[0]
	info_planeta := ""

	for i := 0; i < len(lineas); i++ {
		linea = strings.Split(lineas[i], " ")
		planeta = linea[0]
		if planeta != planeta_actual {
			if existe_planeta(planeta_actual) {
				escribir_archivo(planeta_actual, info_planeta)
			} else {
				crear_planeta(planeta_actual)
				escribir_archivo(planeta_actual, info_planeta)
			}
			info_planeta = ""
		}
		info_planeta = info_planeta + lineas[i] + "\n"
	}
	if existe_planeta(planeta_actual) {
		escribir_archivo(planeta_actual, info_planeta)
	} else {
		crear_planeta(planeta_actual)
		escribir_archivo(planeta_actual, info_planeta)
	}
}

func actualizar_merge_reloj(data string) {
	lineas := strings.Split(data, "\n")
	var linea []string
	var reloj_tmp int
	var reloj = []int32{0, 0, 0}
	for i := 0; i < len(lineas); i++ {
		linea = strings.Split(lineas[i], " ")
		planeta := linea[0]
		reloj_tmp, _ = strconv.Atoi(strings.Split(linea[1], ",")[0])
		reloj[0] = int32(reloj_tmp)
		reloj_tmp, _ = strconv.Atoi(strings.Split(linea[1], ",")[1])
		reloj[1] = int32(reloj_tmp)
		reloj_tmp, _ = strconv.Atoi(strings.Split(linea[1], ",")[2])
		reloj[2] = int32(reloj_tmp)
		actualizar_reloj(planeta, reloj)
	}
}

//
// Funciones GRPC
//

func (s *server) GetCantSoldadosServer(ctx context.Context, in *pb.GetServerRequest) (*pb.GetServerReply, error) {
	log.Printf("Estan haciendo request!")
	planeta := in.GetPlaneta()
	ciudad := in.GetCiudad()

	cant_soldados := obtener_rebeldes(planeta, ciudad)
	log.Printf("Planeta: %s \t Ciudad: %s", planeta, ciudad)
	var reloj = []int32{1, 0, 0}

	return &pb.GetServerReply{Rebeldes: cant_soldados, Reloj: reloj}, nil

}

func (s *server) AskedServer(ctx context.Context, in *pb.AskedServerRequest) (*pb.AskedServerReply, error) {
	log.Printf("El informante %s esta haciendo un comando!", in.GetInformante())

	comando := in.GetComando()
	log.Println("El comando es: \t", comando)

	var reloj = []int32{1, 0, 0}

	// Leer el comando, y enviar a las funciones correspondientes...
	// Todas tienen que responder un reloj.

	return &pb.AskedServerReply{Reloj: reloj}, nil

}

func main() {
	fmt.Println()
	log.Printf("Bienvenido al Servidor Fulcrum, iniciando servicios...")
	rand.Seed(time.Now().UnixNano())

	go func() {
		listner, err := net.Listen("tcp", ":8081")

		if err != nil {
			panic("cannot create tcp connection " + err.Error())
		}

		serv := grpc.NewServer()
		pb.RegisterStarwarsGameServer(serv, &server{})
		if err = serv.Serve(listner); err != nil {
			panic("cannot initialize the server" + err.Error())
		}
	}()

	var opcion string
	fmt.Println("Eres el servidor dominante?")
	fmt.Println("\t [1] Si")
	fmt.Println("\t [2] No")
	fmt.Scanln(&opcion)

	log.Print("Servicios iniciados, escuchando red...")

	flag_opcion := true
	if opcion == "1" {
		for flag_opcion {
			fmt.Println("En 2 minutos se hara un merge")
			time.Sleep(10 * time.Second) //cambiar esto!!!!!!****!*!*!*!*!*!
			//merge_todo()
			fmt.Println("Merge Realizado")
		}
	} else if opcion == "2" {
		for flag_opcion {
			fmt.Println("Esperando 10 segundos...")
			time.Sleep(10 * time.Second)
		}
	} else {
		fmt.Println("Por favor, escribe una opcion correcta!")
	}

	log.Println("Se ha cerrado el proceso.")
}
