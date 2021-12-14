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

var num_servidor int32 = 0 // Cambiar 0, 1 o 2 segun el servidor a ejecutar.

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

	f, err := os.OpenFile("archivos/"+nombre_archivo+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(texto + "\n"); err != nil {
		panic(err)
	}

}

func add_log(nombre_planeta string, nombre_ciudad string, dato_extra string, nombre_accion string) {
	escribir_archivo("log_"+nombre_planeta, nombre_accion+" "+nombre_planeta+" "+nombre_ciudad+" "+dato_extra)
}

func log_string() string {
	var nombre_planeta string
	var planeta Planeta
	nuevo_texto := ""
	for i := 0; i < len(planetas); i++ {
		planeta = planetas[i]
		nombre_planeta = planeta.nombre_planeta

		f, err := os.Open("archivos/" + "log_" + nombre_planeta + ".txt")
		if err != nil {
			return ""
		}
		defer f.Close()

		// Splits on newlines by default.
		scanner := bufio.NewScanner(f)

		line := ""

		for scanner.Scan() {
			line = scanner.Text()
			nuevo_texto = nuevo_texto + line
			if i < len(planetas)-1 {
				nuevo_texto = nuevo_texto + "\n"
			}
		}
	}

	return nuevo_texto
}

func crear_planeta(nombre_planeta string, logear bool) {

	_, err := os.Create("archivos/" + nombre_planeta + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	if logear {
		_, err2 := os.Create("archivos/" + "log_" + nombre_planeta + ".txt")

		if err2 != nil {
			log.Fatal(err2)
		}
	}
	var nuevo_planeta Planeta

	nuevo_planeta.nombre_planeta = nombre_planeta
	nuevo_planeta.reloj = []int32{0, 0, 0}

	planetas = append(planetas, nuevo_planeta)
}

func actualizar_reloj(nombre_planeta string, num_servidor int32) []int32 {
	for i := 0; i < len(planetas); i++ {
		if planetas[i].nombre_planeta == nombre_planeta {
			planetas[i].reloj[num_servidor] = planetas[i].reloj[num_servidor] + 1
			return planetas[i].reloj
		}
	}
	return []int32{-1, -1, -1}
}

func obtener_reloj(nombre_planeta string) []int32 {
	return planetas[buscar_Planeta(nombre_planeta)].reloj
}

func reloj_string() string {
	var nombre_planeta string
	var planeta Planeta
	var reloj string
	nuevo_texto := ""
	for i := 0; i < len(planetas); i++ {
		planeta = planetas[i]
		nombre_planeta = planeta.nombre_planeta
		reloj = strconv.Itoa(int(planeta.reloj[0])) + "," + strconv.Itoa(int(planeta.reloj[1])) + "," + strconv.Itoa(int(planeta.reloj[2]))

		nuevo_texto = nuevo_texto + nombre_planeta + " " + reloj
		if i < len(planetas)-2 {
			nuevo_texto = nuevo_texto + "\n"
		} else if i == len(planetas)-1 {
			nuevo_texto = nuevo_texto
		}
	}
	return nuevo_texto
}

func existe_ciudad(nombre_planeta string, nombre_ciudad string) bool {
	f, err := os.Open("archivos/" + nombre_planeta + ".txt")
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
	if _, err := os.Stat("archivos/" + nombre_planeta + ".txt"); err == nil {
		return true
	}
	return false
}

func crear_ciudad(nombre_planeta string, nombre_ciudad string, cant_soldados int32, logear bool) {

	if existe_planeta(nombre_planeta) {
		if existe_ciudad(nombre_planeta, nombre_ciudad) {
			actualizar_soldados_ciudad(nombre_planeta, nombre_ciudad, cant_soldados)
		} else {
			escribir_archivo(nombre_planeta, nombre_planeta+" "+nombre_ciudad+" "+fmt.Sprint(cant_soldados))
			if logear {
				add_log(nombre_planeta, nombre_ciudad, fmt.Sprint(cant_soldados), "AddCity")
			}
		}

	} else {
		crear_planeta(nombre_planeta, true)
		escribir_archivo(nombre_planeta, nombre_planeta+" "+nombre_ciudad+" "+fmt.Sprint(cant_soldados))
		if logear {
			add_log(nombre_planeta, nombre_ciudad, fmt.Sprint(cant_soldados), "AddCity")
		}
	}
}

func actualizar_nombre_ciudad(nombre_planeta string, nombre_ciudad string, nuevo_nombre_ciudad string) {

	if existe_planeta(nombre_planeta) {
		if existe_ciudad(nombre_planeta, nombre_ciudad) {
			eliminar_ciudad(nombre_planeta, nombre_ciudad, false)
		}
		crear_ciudad(nombre_planeta, nuevo_nombre_ciudad, 0, false)
		add_log(nombre_planeta, nombre_ciudad, nuevo_nombre_ciudad, "UpdateName")
	} else {
		crear_planeta(nombre_planeta, true)
		crear_ciudad(nombre_planeta, nuevo_nombre_ciudad, 0, false)
		add_log(nombre_planeta, nombre_ciudad, nuevo_nombre_ciudad, "AddCity")
	}

}

func actualizar_soldados_ciudad(nombre_planeta string, nombre_ciudad string, cant_soldados int32) {

	if existe_planeta(nombre_planeta) {
		if existe_ciudad(nombre_planeta, nombre_ciudad) {
			eliminar_ciudad(nombre_planeta, nombre_ciudad, false)
		}
		crear_ciudad(nombre_planeta, nombre_ciudad, cant_soldados, false)

		add_log(nombre_planeta, nombre_ciudad, fmt.Sprint(cant_soldados), "UpdateNumber")

	} else {
		crear_planeta(nombre_planeta, true)
		crear_ciudad(nombre_planeta, nombre_ciudad, cant_soldados, false)
		add_log(nombre_planeta, nombre_ciudad, fmt.Sprint(cant_soldados), "AddCity")

	}

}

func eliminar_ciudad(nombre_planeta string, nombre_ciudad string, logear bool) {

	if existe_planeta(nombre_planeta) {
		if existe_ciudad(nombre_planeta, nombre_ciudad) {
			f, _ := os.Open("archivos/" + nombre_planeta + ".txt")
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
			e := os.Remove("archivos/" + nombre_planeta + ".txt")
			if e != nil {
				log.Fatal(e)
			}
			crear_planeta(nombre_planeta, false)
			escribir_archivo(nombre_planeta, nuevo_texto)
			if logear {
				add_log(nombre_planeta, nombre_ciudad, "0", "DeleteCity")
			}
		}
	}
}

func obtener_rebeldes(nombre_planeta string, nombre_ciudad string) int32 {
	f, err := os.Open("archivos/" + nombre_planeta + ".txt")
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

func merge(log_recibido string, num_servidor_log int32) {
	lineas := strings.Split(log_recibido, "\n")
	var nombre_planeta string
	var nombre_ciudad string
	var nuevo_nombre_ciudad string
	var accion string
	var linea string
	var cant_soldados int32
	var temp_int int
	for i := 0; i < len(lineas); i++ {
		linea = lineas[i]
		nombre_planeta = strings.Split(linea, " ")[1]
		nombre_ciudad = strings.Split(linea, " ")[2]
		accion = strings.Split(linea, " ")[0]
		if accion == "UpdateNumber" {
			temp_int, _ = strconv.Atoi(strings.Split(linea, " ")[3])
			cant_soldados = int32(temp_int)
			actualizar_soldados_ciudad(nombre_planeta, nombre_ciudad, cant_soldados)
		} else if accion == "AddCity" {
			temp_int, _ = strconv.Atoi(strings.Split(linea, " ")[3])
			cant_soldados = int32(temp_int)
			crear_ciudad(nombre_planeta, nombre_ciudad, cant_soldados, true)
		} else if accion == "UpdateName" {
			nuevo_nombre_ciudad = strings.Split(linea, " ")[3]
			actualizar_nombre_ciudad(nombre_planeta, nombre_ciudad, nuevo_nombre_ciudad)
		} else if accion == "DeleteCity" {
			eliminar_ciudad(nombre_planeta, nombre_ciudad, true)
		}
		actualizar_reloj(nombre_planeta, num_servidor_log)
	}

}

func clean_logs() {
	for i := 0; i < len(planetas); i++ {
		e := os.Remove("archivos/" + planetas[i].nombre_planeta + ".txt")
		if e != nil {
			log.Fatal(e)
		}

		f, err := os.Create("archivos/" + planetas[i].nombre_planeta + ".txt")

		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
}

func merge_todo() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure()) // Conectamos al IP de 10.6.43.109:8080, el lider.
	if err != nil {
		panic("cannot connect with server " + err.Error())
	}
	serviceClient := pb.NewStarwarsGameClient(conn)
	res, err := serviceClient.GetLogs(context.Background(), &pb.GetLogsRequest{})
	if err != nil {
		panic("No se pudo hacer conexion de merge  " + err.Error())
	}
	log2 := res.GetLog()
	num_servidor_log := res.GetServidor()
	merge(log2, num_servidor_log)
	res2, err2 := serviceClient.GetLogs(context.Background(), &pb.GetLogsRequest{})
	if err2 != nil {
		panic("No se pudo hacer conexion de merge  " + err.Error())
	}
	log3 := res2.GetLog()
	num_servidor_log = res.GetServidor()
	merge(log3, num_servidor_log)
	// res, err := serviceClient.PostReloj(context.Background(), &pb.PostReloj{})
	// if err != nil {
	// 	panic("No se pudo hacer conexion de merge  " + err.Error())
	// }
	fmt.Println("Merge realizado")
}

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
				crear_planeta(planeta_actual, false)
				escribir_archivo(planeta_actual, info_planeta)
			}
			info_planeta = ""
		}
		info_planeta = info_planeta + lineas[i] + "\n"
	}
	if existe_planeta(planeta_actual) {
		escribir_archivo(planeta_actual, info_planeta)
	} else {
		crear_planeta(planeta_actual, false)
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
		actualizar_reloj(planeta, 0)
		reloj_tmp, _ = strconv.Atoi(strings.Split(linea[1], ",")[1])
		reloj[1] = int32(reloj_tmp)
		actualizar_reloj(planeta, 1)
		reloj_tmp, _ = strconv.Atoi(strings.Split(linea[1], ",")[2])
		reloj[2] = int32(reloj_tmp)
		actualizar_reloj(planeta, 2)
	}
}

//
// Funciones GRPC
//

func (s *server) GetLogs(ctx context.Context, in *pb.GetLogsRequest) (*pb.GetLogsReply, error) {
	log.Println("Recibi GetLogs de: ", in.GetNumserver())
	log_string := log_string()
	return &pb.GetLogsReply{Log: log_string, Servidor: num_servidor}, nil
}

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

	var reloj []int32

	// Leer el comando, y enviar a las funciones correspondientes...
	// Todas tienen que responder un reloj.

	splitted_comando := strings.Split(comando, " ")

	if splitted_comando[0] == "AddCity" {
		log.Println("[AddCity]")
		nombre_planeta := splitted_comando[1]
		nombre_ciudad := splitted_comando[2]
		string_soldados := splitted_comando[3]
		int_soldados, _ := strconv.Atoi(string_soldados)
		cant_soldados := int32(int_soldados)
		crear_ciudad(nombre_planeta, nombre_ciudad, cant_soldados, true)
		reloj = actualizar_reloj(nombre_planeta, num_servidor)

	} else if splitted_comando[0] == "UpdateName" {
		log.Println("[UpdateName]")
		nombre_planeta := splitted_comando[1]
		nombre_ciudad := splitted_comando[2]
		nuevo_nombre_ciudad := splitted_comando[3]
		actualizar_nombre_ciudad(nombre_planeta, nombre_ciudad, nuevo_nombre_ciudad)
		reloj = actualizar_reloj(nombre_planeta, num_servidor)

	} else if splitted_comando[0] == "UpdateNumber" {
		log.Println("[UpdateNumber]")
		nombre_planeta := splitted_comando[1]
		nombre_ciudad := splitted_comando[2]
		string_soldados := splitted_comando[3]
		int_soldados, _ := strconv.Atoi(string_soldados)
		cant_soldados := int32(int_soldados)
		actualizar_soldados_ciudad(nombre_planeta, nombre_ciudad, cant_soldados)
		reloj = actualizar_reloj(nombre_planeta, num_servidor)

	} else if splitted_comando[0] == "DeleteCity" {
		log.Println("[DeleteCity]")
		nombre_planeta := splitted_comando[1]
		nombre_ciudad := splitted_comando[2]
		eliminar_ciudad(nombre_planeta, nombre_ciudad, true)
		reloj = actualizar_reloj(nombre_planeta, num_servidor)
	}

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
	os.RemoveAll("archivos/")
	err := os.Mkdir("archivos", 0755)
	if err != nil {
		log.Fatal(err)
	}
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
