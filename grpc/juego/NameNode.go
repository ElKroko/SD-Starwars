package main

import ("math/rand"
		"time"
		"log"
		"strconv"
		"os"
)

// func SearchPlayerPlays(player int) int[] {
// 	f, err := os.Open("lugarJugador.txt")
//     var splitText string
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer f.Close()
//     scanner := bufio.NewScanner(f)
//     for scanner.Scan() {
// 		splitText = strings.Split(scanner.Text(), " ")
// 		if splitText[0] == "Jugador_" + player {
//         	llamarDataNode(splitText[0][-1], splitText[1][-1], splitText[2])
// 		}
//     }
// }

func main() {
	rand.Seed(time.Now().UnixNano())
	f, err := os.Create("lugarJugadores.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
	next := false
	var index int64 = 0
	var ip string
	for next == false {
		//esperar datos
		jugador, ronda := 1, 2
		selectDataNode := rand.Intn(3)
		if selectDataNode == 0 {
			ip = ip1
		} else if selectDataNode == 1 {
			ip = ip2
		} else {
			ip = ip3s
		}
		val := "Jugador_" + strconv.Itoa(jugador) + " Ronda_" + strconv.Itoa(ronda) + " " + ip + "\n"
		data := []byte(val)
		_, err := f.WriteAt(data, index)
		index = int64(len(data))+index
		if err != nil {
			log.Fatal(err)
		}
		//mandar jugador y ronda a ip
	}
	
}