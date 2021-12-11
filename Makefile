leia:
	echo "Ejecutando Leia"
	go run Leia/Leia.go

broker:
	echo "Ejecutando Broker"
	go run Server/Broker.go

pozo:
	echo "Ejecutando pozo"
	go run server/Pozo.go

datanode:
	echo "Ejecutando datanode"
	go run server/DataNode.go
	
jugador:
	echo "Ejecutando jugadores"
	go run client/Jugador.go
