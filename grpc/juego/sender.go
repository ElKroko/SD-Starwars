package main

import (
  "log"
  amqp "github.com/streadway/amqp"
  "net"
  "grpc"
  "pb"
  "server"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func main() {
	go func() {
		listner, err := net.Listen("tcp", ":8083")

		if err != nil {
			panic("cannot create tcp connection " + err.Error())
		}

		serv := grpc.NewServer()
		pb.RegisterCalamardoGameServer(serv, &server{})
		if err = serv.Serve(listner); err != nil {
			panic("cannot initialize the server" + err.Error())
		}
	}()

	conn, err := amqp.Dial("amqp://test:test@10.6.43.110:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()


	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	i := "1"
	s := "2"

	body := "Jugador_" + i + " Ronda_" + s

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf("Ha muerto: %s ", body)
}