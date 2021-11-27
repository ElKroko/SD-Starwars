lider:
	echo "Ejecutando lider"
	go run server/Lider.go

namenode:
	echo "Ejecutando namenode"
	go run server/NameNode.go

pozo:
	echo "Ejecutando pozo"
	go run server/Pozo.go

datanode:
	echo "Ejecutando datanode"
	go run server/DataNode.go
	
jugador:
	echo "Ejecutando jugadores"
	go run client/Jugador.go
