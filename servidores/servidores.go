package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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
	escribir_archivo(nombre_planeta, nombre_accion+" "+nombre_planeta+" "+nombre_ciudad+" "+fmt.Sprint(cant_soldados))
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
			nuevo_texto = nuevo_texto + scanner.Text()
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

func obtener_rebeldes(nombre_planeta string, nombre_ciudad string) int {
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
			return cant_soldados
		}
	}

	return -1
}
