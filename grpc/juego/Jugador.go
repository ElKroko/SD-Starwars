package main

import (
	"fmt"
	"math/rand"
	"time"
)

func preguntar_juego() {
	fmt.Println("Preguntar al lider si puedo jugar...")
}

func juego_1(jugador bool, ronda int) int {
	fmt.Println("Hola!")
	rand.Seed(time.Now().UnixNano())

	suma := 0
	var numero int
	gano := 0

	for ronda < 4 {

		// Seleccionar numero para jugar
		if jugador {
			fmt.Println("Elija un numero del 1 al 10")
			fmt.Print("-> ")
			fmt.Scanln(&numero)

		} else {
			numero = rand.Intn(10)

		}

		// Chequear si el numero es igual o mayor al del lider

		muere := false

		// Fin

		if !muere {
			suma += numero
			fmt.Println("La suma da: ", suma)
		} else {
			gano = 0
			break
		}

		// sumar al numero anterior
		if suma >= 21 {
			gano = 1
			break
		}

		fmt.Println(numero)
		ronda = ronda + 1
	}

	return gano

}

func juego2_jugador() int {
	fmt.Println("Hola!")
	fmt.Println("Elija un numero del 1 al 4")
	fmt.Print("-> ")
	fmt.Scanln(&numero)
	fmt.Println(resultado)
	return 1
}

func juego3_jugador() int {
	fmt.Println("Hola!")
	fmt.Println("Elija un numero del 1 al 10")
	fmt.Print("-> ")
	fmt.Scanln(&numero)
	fmt.Println(resultado)
	return 1
}

func v_o_f(resultado int) bool {
	if resultado == 1 {
		fmt.Println("gane!")
		return true
	} else {
		fmt.Println("perdi :c")
		return false
	}
}

func opc_pozo() int {
	fmt.Println("Elija que quiere hacer: ")
	fmt.Println("[1] Preguntar por el valor acumulado en el pozo")
	fmt.Println("[2] Continuar con la siguiente etapa")

	var opcion int
	fmt.Print("-> ")
	fmt.Scanln(&opcion)

	return opcion
}

func preguntar_pozo() {
	fmt.Println("Preguntamos al pozo por el total")

}

func main() {

	// Enviar peticion para unirse al juego y recibir true o false
	preguntar_juego()

	var opcion, resultado int

	// if preguntar juego
	if true {
		// Como voy a saber si es un jugador o un bot?
		resultado = juego1_jugador(true, 0)
		v_o_f(resultado)

		// Segundo juego, si es que sigue vivo
		opcion = opc_pozo()

		if opcion == 1 {
			preguntar_pozo()
			resultado = juego2_jugador()
		} else {
			resultado = juego2_jugador()
		}
		v_o_f(resultado)


		// Tercer juego, si es que sigue vivo
		opcion = opc_pozo()

		if opcion == 1 {
			preguntar_pozo()
			resultado = juego3_jugador()
		} else {
			resultado = juego3_jugador()
		}
		v_o_f(resultado)

	} else {
		fmt.Println("No se pudo unir a un juego, te tendremos que matar.")
	}

}
