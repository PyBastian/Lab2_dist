# Squid Game

Este documento mostrará y documentará informacion relacionada a la tarea y a las etapas del proceso, como correr los archivos, supuestos que se tengan y cúal fue la lógica usada tras los scripts implementados en Go.
### Integrantes grupo 54
- Bastián Vivar 201773109-5
- Rodrigo Hernandez 201730036-1
- Luis Blanco

#### La distribución de las máquinas virtuales y sus componentes son los siguentes:
- dist213 ---> Node/ 
- dist214 ---> Lider / Server
- dist215 ---> Pozo
- dist216 ---> Jugadores
## Jugadores
Los jugadores actuan como clientes del sistema, donde para cada jugador se conectara al lider (greeter_server) mediante gRCP. Solamente utilizaran un servicio de mensajeria definido en helloworld.proto el cual se utilizara para hacer las conexiones y comunicaciones entre todos los sistemas en la red. Los jugadores ingresaran al juego mediante una conexion con el lider en el puerto 50071, a los cuales se envia un mensaje de bienvenida y se le preguta si quieren ingresar al juego del calamar.
En caso de que su respuesta sea afimativa, se le comunica al lider y el contador de jugadores interno del lider aumenta en 1 y asi hasta llegar al numero maximo de jugadores. Cada jugador debera esperar a que se complete el cupo de jugadores disponibles y entonces el lider informara a cada jugador el comienzo del juego del calamar, a lo cual se seleccionará un juego y los jugadores empezarán a enviar sus jugadas al lider. Si el jugador es eliminado, gana o muere, se le informa por consola y se actualiza el pozo con un llamado que queda a cargo del lider. Debido a temas de tiempo y por lo explicado en el discord del curso, no fue posible por temas de tiempo implementar a los jugadores bots, pero la lógica que seguirian estos seria identica a la de un jugador por consola, solamente que cuando se pregunte una accion, el sistema envie automaticamente esta accion con un selector Random(). Estos estarian ejecutandose igualmente en esta maquina.

## Lider/Server
Nuestra entidad Lider es utilziada como servidor inicial donde esta sera la que interactuara con nosotros y se comunicara con todas las entidades y a su vez ellas a el a través de strings y además las entidades también generaran procesos que se comunicaran con el, donde el Lider en un inicio esperara que los jugadores se contacten con el, nosotros definimos la espera de un usuario, ya que los bots si bien pueden ser copias de la entidad Jugadores con otros puertos, y también podrian ser llevados como metodos dentro de jugadores y lo manejamos de manera local, pero nosotros desarrollamos el juego con un jugador para fines practicos de testeo.

Luego podemos ver que la comunicación entre el cliente (Jugadores), se llevara a cabo con cada jugada y se analizan los strings con mensajes convenientes, para luego ser procesados por el Lider (greater_ser/main.go) y dependiendo de los diferentes outputs el Lider se comunicara con la entidad correspondiente, donde si bien estan implementados no lograron llevar a cabo la correcta comunicación entre las maquinas

## Pozo

El pozo se mantiene comunicado con el lider y mantiene el conteo del pozo según lo que informe el lider por jugador eliminado. Cada jugador se registra en el archivo "Pozo_Acumulado.txt" a medida que es eliminado se registra y actualiza. La conexión es realizada mediante RabbitMQ y las peticiones de la cantidad de dinero acumulado se realizan con gRPC.

## NameNode - Datanode
El namenode contiene las direcciones de los 3 datanodes que se encargan de mantener las jugadas por id de cada juego mediante la actualizacion. Para escribir en los archivos se utiliza la misma funcion build-in de Go, WriteinTXT. Así cada vez que se actualiza busca el datanode que corresponde a la instancia que se esta jugando y ejecuta una actualizacion del archivo txt. Este proceso deberá estar corriendo en 3 maquinas virtuales pero por temas de tiempo antes explicados y por el desarrollo en local que se hizo de la tarea, no se pudo hacer conectar los 3 datanodes administrador por cada una **en las maquinas virtuales** donde localmente corrian por consola junto con los otros procesos. Las conexiones se hacen de manera sincrona con gRCP con la misma lógica implementada anteriormente y espera una respuesta del lider para poder trabajr actualizando los archivos. Sus puertos de defecto eran los 50081, 50082 y 50083.

## Ejecución Codigo
Al ejecutar el juego sin utilizar make, nosotros primero hacemos en /Distribuidos_w en la maquina Dist214 go run greeter_server/main.go inicializando el Lider, para luego en Dist216 en /Distribuidos_w (todas seran ahi ya que utilizamos git para subir codigo a las maquinas) go run greeter_client/main.go y luego en dist215 go run greeter_pozo/main.go y finalmente en dist215 go run greeter_pozo/main.go, utilizando los Maklefile de cada maquina simplemente es el make en /Distribuidos_w