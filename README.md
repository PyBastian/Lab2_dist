# Laboratorio 3: Star Wars

Este documento mostrará y documentará informacion relacionada a la tarea y a las etapas del proceso, como correr los archivos, supuestos que se tengan y cuál fue la lógica usada tras los scripts implementados en Go.

## Integrantes grupo 54
- Bastián Vivar 201773109-5
- Rodrigo Hernandez 201730036-1
- Luis Blanco 201573027-K

### La distribución de las máquinas virtuales y sus componentes son los siguentes:
- dist213 ---> Informante Ahsoka Tano / Servidor Fulcrum 
- dist214 ---> Broker
- dist215 ---> Informante Almirante Thrawn / Servidor Fulcrum 2
- dist216 ---> Leia Organa / Servidor Fulcrum 3
  
## Informantes
- Inicialmente envia los comando al Broker, el cual respondera con la dirección del servidor fulcrum al cual debera reenviar el comando 
- Para luego recibir la dirección y comunicarse con este, donde cabe desatcar que el serviidor puede estar en su maquina y por ende se ejecuta el comando localmente, y en los otros casos mandara el comando al servidor fulcrum.

## Broker/Server
- Recibir las comunicaciones de los dos informantes y elegir según el criterio de consistencia para que estos comandos tengan sentido al servidor fulcrum que direcciona, 
- Recibe los comandos de Leia y este consulta a los servidores Fulcrum, donde utilizando Monotonic Reads, el valor que este retorne siempre sera igual o más actualizado que las lecturas anterior.
  
## Servidores Fulcrum
- Recibir comandos...LISTO
- Ejecutar funcionalidades de Add, Update, Delete...LISTO
- Guardar registro de planetas en .txt...LISTO
- Retornar vector reloj del cambio que se hizo PENDIENTE
- Revisar consistencia entre servidores PENDIENTE
- Merge PENDIENTE

## Leia Organa
- Leia envia comandos al Broker, donde ella guardara los comandos que ha solicitado, y además guardara el reloj del planeta de lainformación que solicito y cual fue la maquina que tenia guardado este registro.

## Ejecución Codigo

Todo el desarrollo de nuestro sistema se lleva a cabo con github, por lo  que las maquinas tienen el codigo de todo el sistema pero cada una ejecuta un go run especifico, es por esto que nuestras maquina dentro del directorio Lab2_dist, cada una posee un Makefile especifico que ejecuta su logica correspondiente, por lo que para el testeo y funcionamiento del sistema se debe ejecutar un "make" en cada maquina partiendo por la maquina 214, la cual posee el make que ejecutara "go run greeter_server/main.go" o sea el broker, para ejecutar el make en cada una de las otras maquinas sin un orden necesario especifico.




