﻿# Laboratorio 3: Star Wars

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
- Inicialmente pide los comandos por terminal para luego enviar los comando al Broker, el cual respondera con la dirección del servidor fulcrum al cual debera reenviar el comando 
- Para luego estos sean capaces de recibir la dirección y comunicarse con este, donde cabe destcar que el servidor puede estar en su maquina y por ende se ejecuta el comando localmente asociado a crear, eliminar o modificar los .txt, y en los otros casos mandara el comando al servidor fulcrum.

## Broker/Server
- Recibir las comunicaciones de los dos informantes y elegir según nuestro criterio de consistencia para que estos comandos tengan sentido al servidor fulcrum que direcciona, 
- Recibe los comandos de Leia y este consulta a los servidores Fulcrum, donde utilizando Monotonic Reads, el valor que este retorne siempre sera igual o más actualizado que las lecturas anterior.
  
## Servidores Fulcrum
- Recibir comandos de los informantes y el broker.
- Ejecutar funcionalidades de Add, Update, Delete asociadas a .txt según el comando que reciban
- Guardar registro de planetas en .txt
- Retornar vector reloj del planeta que modifico o creo 
- Revisar consistencia entre servidores 
- Merge cada 2 minutos

## Leia Organa
- Leia envia comandos al Broker, donde ella guardara los comandos que ha solicitado en su misma maquina en la carpeta greeter_leia/. con formato .txt, y además almacenara el reloj del planeta de la información que solicito y cual fue la maquina de la cual proviene este registro.

## Ejecución Codigo

Todo el desarrollo de nuestro sistema se lleva a cabo con github, por lo  que las maquinas tienen el codigo de todo el sistema pero cada una ejecuta un go run especifico, es por esto que nuestras maquina dentro del directorio Lab2_dist, cada una posee un Makefile especifico que ejecuta su logica correspondiente, por lo que para el testeo y funcionamiento del sistema se debe ejecutar un "make" en cada maquina partiendo por la maquina 214, la cual posee el make que ejecutara "go run greeter_server/main.go" o sea el broker, para ejecutar el make en cada una de las otras maquinas sin un orden necesario especifico.




