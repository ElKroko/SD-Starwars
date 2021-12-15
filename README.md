# SD-Starwars
Tarea 3 de SD a implementar con GO y gRPC
Integrantes: 
- Vicente Perelli - ROL: 201756594-2
- Nicolas Maldonado - ROL: 201360619-9


### Ejecucion
Para la correcta ejecucion de estos programas, es necesario ubicarlos en las maquinas virtuales segun el siguiente esquema:

[insertar img]

Para lograr esto, es necesario abrir y ejecutar:
- 1 terminal en dist@121 - `$make broker`
- 2 terminales en dist@122 (Serv Dom y Leia) - `$make server` y `$make leia` respectivamente.
- 2 terminales en dist@123 (Serv Sub y Ashoka Tano) - `$make server` y `$make informante` respectivamente.
- 2 terminales en dist@124 (Serv sub y Almirante Thrawn) - `$make server` y `$make informante` respectivamente.

<sup>* Los archivos automaticamente disciernen que tienen que hacer, basandose en la IP de la maquina virtual. </sup>


#### Funcionamiento:
Los servidores se inicializan como vacios, por tanto hay que insertar planetas a traves de los informantes antes de consultar con Leia.

Para insertar planetas, es necesario ejecutar primeramente los servidores y el broker en sus respectivas maquinas virtuales, y los servidores deben mostrar dentro de sus CLI: `"En 2 minutos se hara un merge"` o `"Esperando 10 segundos..."` 

Luego de la insercion de los planetas, cada 2 minutos se va a ejecutar automaticamente un merge que compartira la informacion entre los servidores.

Si Leia Organa hace una consulta y el planeta consultado no se encuentra en el servidor aleatorio que le toco, se ejecutara un merge (sin romper el equilibrio de los merge realizados automaticamente) y luego volvera a consultar, para obtener la informacion de los soldados del planeta deseado.

###Servidores Fulcrum:
Existen 3 servidores:
1. **Servidor Fulcrum Dominador (IP: 10.6.43.110)** - Este servidor se encarga de recibir todos los logs y archivos de los otros servidores, compararlos y luego enviar las actualizaciones a los otros servidores.
2. **Servidor Fulcrum sub (IP: 10.6.43.111)** - Este servidor envia sus logs y recibe actualizaciones
3. **Servidor fulcrum sub (IP: 10.6.43.112)** - Este servidor envia sus logs y recibe actualizaciones

