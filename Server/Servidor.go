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

var num_servidor int32 = -1 // Cambiar 0, 1 o 2 segun el servidor a ejecutar.

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

// funcion nueva
func escribir_archivo(nombre_archivo string, texto string) {
	f, err := os.OpenFile("archivos/"+nombre_archivo+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	datawriter := bufio.NewWriter(f)

	fstat, _ := f.Stat()
	log.Println("[Escribir archivo] nombre: " + nombre_archivo + " [Texto a escribir]: " + texto)
	if fstat.Size() == 0 {
		_, _ = datawriter.WriteString(texto)
	} else {
		_, _ = datawriter.WriteString("\n" + texto)
	}

	datawriter.Flush()
	f.Close()

}

func escribir_archivo2(nombre_archivo string, texto string) {

	f, err := os.OpenFile("archivos/"+nombre_archivo+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	fi, _ := f.Stat()

	defer f.Close()

	log.Println("[Escribir archivo] nombre: " + nombre_archivo + " [Texto a escribir]: " + texto)
	if fi.Size() == 0 {
		_, err2 := f.WriteString(texto)
		if err2 != nil {
			panic(err2)
		}
	} else {
		_, err2 := f.WriteString("\n" + texto)
		if err2 != nil {
			panic(err2)
		}
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
		if i < len(planetas)-1 {
			nuevo_texto = nuevo_texto + "\n"
		}
	}
	return nuevo_texto
}

func planetas_string() string {
	var nombre_planeta string
	var planeta Planeta
	nuevo_texto := ""
	for i := 0; i < len(planetas); i++ {
		planeta = planetas[i]
		nombre_planeta = planeta.nombre_planeta

		f, err := os.Open("archivos/" + nombre_planeta + ".txt")
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
		log.Println("Crear ciudad: " + nombre_planeta)
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
					if nuevo_texto == "" {
						nuevo_texto = scanner.Text()
					} else {
						nuevo_texto = nuevo_texto + "\n" + scanner.Text()
					}
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
	if log_recibido != "" {
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
			accion = strings.Split(linea, " ")[0]
			nombre_planeta = strings.Split(linea, " ")[1]
			nombre_ciudad = strings.Split(linea, " ")[2]

			if accion == "UpdateNumber" {
				temp_int, _ = strconv.Atoi(strings.Split(linea, " ")[3])
				cant_soldados = int32(temp_int)
				actualizar_soldados_ciudad(nombre_planeta, nombre_ciudad, cant_soldados)
			} else if accion == "AddCity" {
				temp_int, _ = strconv.Atoi(strings.Split(linea, " ")[3])
				cant_soldados = int32(temp_int)
				crear_ciudad(nombre_planeta, nombre_ciudad, cant_soldados, false)
			} else if accion == "UpdateName" {
				nuevo_nombre_ciudad = strings.Split(linea, " ")[3]
				actualizar_nombre_ciudad(nombre_planeta, nombre_ciudad, nuevo_nombre_ciudad)
			} else if accion == "DeleteCity" {
				eliminar_ciudad(nombre_planeta, nombre_ciudad, false)
			}
			actualizar_reloj(nombre_planeta, num_servidor_log)
		}
	}
}

func clean_logs() {
	for i := 0; i < len(planetas); i++ {
		e := os.Remove("archivos/log_" + planetas[i].nombre_planeta + ".txt")
		if e != nil {
			log.Fatal(e)
		}

		f, err := os.Create("archivos/log_" + planetas[i].nombre_planeta + ".txt")

		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
}

func clean_planetas() {
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

func merge_conexion(IP string) {
	conn, err := grpc.Dial(IP+":8081", grpc.WithInsecure())
	if err != nil {
		panic("cannot connect with server " + err.Error())
	}
	serviceClient := pb.NewStarwarsGameClient(conn)
	res, err := serviceClient.GetLogs(context.Background(), &pb.GetLogsRequest{Numserver: num_servidor})
	if err != nil {
		panic("No se pudo hacer conexion de merge  " + err.Error())
	}
	log2 := res.GetLog()
	log.Print("[Log recibido]: " + log2)
	num_servidor_log := res.GetServidor()
	merge(log2, num_servidor_log)
}

func postmerge_conexion(IP string, reloj string, planetas_postmerge string) bool {
	conn, err := grpc.Dial(IP+":8081", grpc.WithInsecure())
	if err != nil {
		panic("cannot connect with server " + err.Error())
	}
	serviceClient := pb.NewStarwarsGameClient(conn)
	res, err := serviceClient.PostMerge(context.Background(), &pb.PostMergeRequest{Reloj: reloj, Planetas: planetas_postmerge})
	if err != nil {
		panic("No se pudo hacer conexion de merge  " + err.Error())
	}
	ack := res.GetAck()
	return ack
}

// Esta funcion es llamada x el dominatrix
func merge_todo(IP1 string, IP2 string) {

	merge_conexion(IP1)
	merge_conexion(IP2)

	clean_logs()

	reloj := reloj_string()
	planetas_merge := planetas_string()
	log.Println("[PreMerge] Planetas merge: ", planetas_merge)
	log.Println("[PreMerge] Reloj merge: ", reloj)

	ack1 := postmerge_conexion(IP1, reloj, planetas_merge)
	ack2 := postmerge_conexion(IP2, reloj, planetas_merge)

	if ack1 && ack2 {
		fmt.Println("Merge realizado")
	} else {
		fmt.Println("Merge fallido :cccccccc")
	}

}

func actualizar_merge_planetas(data string) {
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
				planeta_actual = planeta
			} else {
				crear_planeta(planeta_actual, true)
				escribir_archivo(planeta_actual, info_planeta)
				planeta_actual = planeta
			}
			info_planeta = ""
		}
		if i < len(lineas)-1 {
			info_planeta = info_planeta + lineas[i] + "\n"
		} else {
			info_planeta = info_planeta + lineas[i]
		}
	}
	if planeta_actual == "" {
		if existe_planeta(planeta_actual) {
			escribir_archivo(planeta_actual, info_planeta)
		} else {
			crear_planeta(planeta_actual, true)
			escribir_archivo(planeta_actual, info_planeta)
		}
	}
}

func actualizar_merge_reloj(data string) {
	lineas := strings.Split(data, "\n")
	var linea []string
	var reloj_tmp int
	var reloj = []int32{0, 0, 0}
	var reloj_actual = []int32{0, 0, 0}
	for i := 0; i < len(lineas); i++ {
		linea = strings.Split(lineas[i], " ")
		planeta := linea[0]
		if planeta != "" {
			if !existe_planeta(planeta) {
				crear_planeta(planeta, true)
				reloj_actual = []int32{0, 0, 0}
			} else {
				reloj_actual = planetas[buscar_Planeta(planeta)].reloj
			}
			reloj_tmp, _ = strconv.Atoi(strings.Split(linea[1], ",")[0])
			reloj[0] = int32(reloj_tmp)
			for j := reloj_actual[0]; j < reloj[0]; j++ {
				actualizar_reloj(planeta, 0)
			}
			reloj_tmp, _ = strconv.Atoi(strings.Split(linea[1], ",")[1])
			reloj[1] = int32(reloj_tmp)
			for j := reloj_actual[1]; j < reloj[1]; j++ {
				actualizar_reloj(planeta, 1)
			}
			reloj_tmp, _ = strconv.Atoi(strings.Split(linea[1], ",")[2])
			reloj[2] = int32(reloj_tmp)
			for j := reloj_actual[2]; j < reloj[2]; j++ {
				actualizar_reloj(planeta, 2)
			}
		}
	}
}

//
// Funciones GRPC
//

func (s *server) GetLogs(ctx context.Context, in *pb.GetLogsRequest) (*pb.GetLogsReply, error) {
	log.Println("El servidor", in.GetNumserver(), " esta mandando un GetLogs")

	log_a_enviar := log_string()
	log.Println("[Log a enviar]: " + log_a_enviar)

	// Revisar que enviar de servidor, si su IP o el numero  que lo identifica

	return &pb.GetLogsReply{Log: log_a_enviar, Servidor: num_servidor}, nil
}

func (s *server) PostMerge(ctx context.Context, in *pb.PostMergeRequest) (*pb.PostMergeReply, error) {

	log.Println("PostMerge!")

	reloj := in.GetReloj()
	planetas_merge := in.GetPlanetas()
	log.Println("[PostMerge] Reloj Merge: ", reloj, "\tPlanetas merge: ", planetas_merge)

	clean_planetas()
	actualizar_merge_planetas(planetas_merge)
	actualizar_merge_reloj(reloj)
	clean_logs()

	return &pb.PostMergeReply{Ack: true}, nil
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

func main() {
	fmt.Println()

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
	ip := GetIP()
	if ip == "10.6.43.110" {
		log.Printf("Bienvenido al Fulcrum Dominante, iniciando servicios...")
		num_servidor = 0
	} else if ip == "10.6.43.111" {
		log.Printf("Bienvenido al Servidor Fulcrum, iniciando servicios...")
		num_servidor = 1
	} else if ip == "10.6.43.112" {
		log.Printf("Bienvenido al Servidor Fulcrum, iniciando servicios...")
		num_servidor = 2
	}

	log.Print("Servicios iniciados, escuchando red...")

	flag_opcion := true
	if num_servidor == 0 {
		for flag_opcion {
			fmt.Println("En 2 minutos se hara un merge")
			time.Sleep(30 * time.Second) //cambiar esto!!!!!!****!*!*!*!*!*!
			merge_todo("10.6.43.111", "10.6.43.112")
			fmt.Println("Merge Realizado")
			fmt.Println(planetas)
			fmt.Println("")
		}
	} else if num_servidor == 1 || num_servidor == 2 {
		for flag_opcion {
			fmt.Println("Esperando 10 segundos...")
			fmt.Println(planetas)
			time.Sleep(10 * time.Second)
		}
	} else {
		fmt.Println("Por favor, escribe una opcion correcta!")
	}

	log.Println("Se ha cerrado el proceso.")
}
