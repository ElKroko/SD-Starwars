package main

import ("fmt"
		"os"
		"log"
        "bufio"
        "strings"
)

func ActualAmmount() string {
    f, err := os.Open("JugadoresEliminados.txt")
    var actualAmmount string
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        actualAmmount = strings.Split(scanner.Text(), " ")[2]
    }
    return actualAmmount
}

func main() {
	f, err := os.Create("JugadoresEliminados.txt")
    var index int64 = 0
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    next := false
	for next == false {
        val := "Jugador_1 Ronda_1 10000\n"//espera que manden datos
        data := []byte(val)
        _, err := f.WriteAt(data, index)
        index = int64(len(data))+index
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println(ActualAmmount())

}