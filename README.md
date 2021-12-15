# Laboratorio 3: Star Wars

Este documento mostrará y documentará informacion relacionada a la tarea y a las etapas del proceso, como correr los archivos, supuestos que se tengan y cuál fue la lógica usada tras los scripts implementados en Go.

## Integrantes grupo 54
- Bastián Vivar 201773109-5
- Rodrigo Hernandez 201730036-1
- Luis Blanco 201573027-K
  
## Informantes
- Inicialmente pide los comandos por terminal para luego enviar los comando al Broker, el cual respondera con una maquina que se le asigne al informante mediante la logica interna que manejara el broker. Asi el broken contendra los comandos junto con los relojes de informacion almacenados en memoria para entregar la consulta y escritura que se le haga dependiendo de la maquina que la requiera.
- Estos son capaces de recibir la dirección y comunicarse con este, donde cabe destcar que el servidor puede estar en su maquina y por ende se ejecuta el comando localmente asociado a crear, eliminar o modificar los .txt, y en los otros casos mandara el comando al servidor fulcrum.

## Broker/Server
- Recibir las comunicaciones de los dos informantes y elegir según nuestro criterio de consistencia para que estos comandos tengan sentido al servidor fulcrum que direcciona. Se reciben y guardan los comandos que enviar los informantes en memoria y se establece un orden de procedencia de los comandos dependiendo del reloj del planeta asociado. Asi en memoria se encontraran los relojes asociados a cada planeta en caso de que leia requiera consultar el numero de informates mas actualizado con la maquina que hizo el ultimo cambio al reloj.
- Recibe los comandos de Leia y este consulta a los servidores Fulcrum con los comandos antes ingresados, donde utilizando Monotonic Reads, el valor que este retorne siempre sera igual o más actualizado que las lecturas anterior con la ayuda de los relojes de vectores que se enviar al broker.
  
## Servidores Fulcrum
- Recibir comandos de los informantes y el broker.
- Mantiene un registro de los comandos que se hacen en esta maquina, con el usuario que los hizo y el reloj correspondiete. Este archivo se crea automaticamente en la carpeta del servidor fulcrum y se llama "historial.txt".
- Ejecutar funcionalidades de Add, Update, Delete asociadas a .txt según el comando que reciban
- Cada vez que se ejecuta un comando de cambio a un planeta, se suma un 1 al reloj asociado y se guarda en historial.txt
- Retornar vector reloj del planeta que modifico o creo 
- Revisar consistencia entre servidores
- Para mantener la consistencia Read-your-Writes cada maquina mantiene y lee los archivos locales que aun no se han hecho merge, para no sobrecargar el servidor con peticiones. Asi solamente leia realiza peticiones de Read al servidor y se redirige segun el reloj que lleva localmente el broker con los planetas.
- Merge cada 2 minutos mediante una peticion al broker de update. El broken envia los comandos que se hicieron a los servidores fulcrum paralelos durante los 2 minutos que no se ha hecho el merge y los ejecuta en cada maquina para mantenerlas actualizadas y mantener la consistencia Read your Writes.

## Leia Organa
- Leia envia comandos al Broker, donde ella guardara los comandos que ha solicitado en su misma maquina en la carpeta greeter_leia/. con formato .txt, y además almacenara el reloj del planeta de la información que solicito y cual fue la maquina de la cual proviene este registro.
- Leia se redirecciona a la ultima actualizacion del planeta en cuestion que tenga el reloj mas alto (ya que fue el con mas cambios) y se le entrega esta informacion, manteniendo la consistencia Monotic Reads

## Merge
El merge se realiza simplemente enviando los comandos ingresados en los 2 ultimos minutos a todas las maquinas que no han ejecutado los comandos ingresados por la redireccion del broker. Los comandos son leidos de los log de cada server fulcrum y junto con el reloj de versiones que manjeja el broker, siempre elegiendo el mejor reloj, envia una lista de comandos a las maquinas que no ejecutaron los comandos, con los que corriendo en secuencia, al final se logra un merge completo de los comandos.

## Ejecución Codigo

Todo el desarrollo de nuestro sistema se lleva a cabo con github, por lo  que las maquinas tienen el codigo de todo el sistema pero cada una ejecuta un go run especifico, es por esto que nuestras maquina dentro del directorio Lab2_dist, cada una posee un Makefile especifico que ejecuta su logica correspondiente, por lo que para el testeo y funcionamiento del sistema se debe ejecutar un "make" en cada maquina partiendo por la maquina 214, la cual posee el make que ejecutara "go run greeter_server/main.go" o sea el broker, para ejecutar el make en cada una de las otras maquinas sin un orden necesario especifico.

Dentro de cada maquina se encuentra un archivo Make el cual se encarga de ejecutar los servidores fulclrum. Para correr los archivos es necesario correr los archivos Make en cada maquina virtual dentro de un intervalo minimo de 2 minutos (ya que se envia una actualizacion a los 2 minutos) y luego simplemente enviar comandos con los informantes o leia. Notar que tiene que existir el planeta para que Leia pueda pedir el numero de rebeldes en el. 

La distribucion de las maquinas virtuales es la siguiente:
- dist213 ---> Informante Thrawn / Fulcrum2
- dist214 ---> Broker Mos Eisley
- dist215 ---> Leia / Fulcrum3
- dist216 ---> Jugadores -->  Informante Ahsoka / Fulcrum1 (Nodo dominante)