# CalamardoDistribuido

#### Integrantes: 
```
Vicente Perelli Tassara - Rol: 201756594-2
Nicolas Maldonado Fernandez - Rol: 201360619-9
```
Contabamos con un tercer integrante, pero dejo de contestar nuestros mensajes. El estaba encargado de investigar e implementar gRPC.
No mostro ningun tipo de avance, ni los problemas que encontro en su investigacion a pesar de que le preguntamos directamente.

Se envio un correo al ayudante *Jorge Diaz* el dia *viernes 5* para notificar esta situacion.

### Comandos
Para compilar el .proto, usamos:
`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/whishlist.proto`

Para ejecutar los archivos, es necesario estar ubicado en la carpeta grpc/

- Jugador: `make jugador`
- Lider: `make lider`
- Datanode: `make datanode`
- Namenode: `make namenode`
- Pozo: `make pozo`

Todo debe estar ejecutandose al mismo tiempo ANTES de introducir el **tipo de jugador**.

### Consideraciones:

Para el caso de Jugador:
- Al momento de iniciar el juego, preguntara el **tipo de jugador**:
  - Jugador significa que el usuario introducira las jugadas por cmd
  - De lo contrario, funcionara como un bot.
- Cada instancia de jugador debe ser abierta en una nueva terminal.

Para **Namenode** y **Datanode**:
- Los programas deben ejecutarse, y para terminar con su funcionamiento es necesario introducir *exactamente* **"T"**.

### Lugar de Ejecucion:
A continuacion, se especificara donde se deben ejecutar los diferentes procesos:

- **Lider:** maquina virtual **dist121** (ip: *10.6.43.109*)
- **Jugador:** en cualquier maquina virtual **excepto la del lider** (*dist121*), puede ser todos en la misma o distribuido en diferentes maquinas.
- **NameNode:** maquina virtual **dist124** (ip: *10.6.43.112*)
- **Datanode:** maquinas **dist121**, **dist122**, **dist123** (ips: *10.6.43.109* - *10.6.43.110* - *10.6.43.111*)

### Supuestos:

  El usuario siempre entregara el input correctamente, sin caracteres demas.


#### Problemas encontrados:
Existe un problema de asignacion de listas, por el cual la funcion *indexOfPlayers()* no encuentra a un jugador, lo que gatilla un `panic:runtime error: index out of range [-1]`, por mas que tratamos de investigar no funciono. 

El problema es random, y hemos notado que para la mayoria de las situaciones se soluciona al ejecutar nuevamente los programas.

Para el programa **Pozo.go**, esta implementada toda la logica pero no alcanzamos a conectar su funcionalidad.
