# Alumnos: 
- Stephano Helí Morales Linares
- Luis Gustavo Becerra Bisso
- Sebastián Arana del Carpio

# Introducción del Juego
El juego de nodos es una implementación en Go de un juego basado en nodos que utiliza conexiones TCP para la comunicación entre los nodos. El objetivo del juego es simular un salto de un nodo a otro, donde cada nodo representa a un jugador en un equipo y los jugadores compiten por llegar al nodo final. Este informe proporcionará una visión general del diseño del juego, las funciones principales y las tecnologías utilizadas.


## Diseño del juego:
El juego de nodos se basa en un modelo de comunicación cliente-servidor utilizando sockets TCP. Cada nodo actúa como un servidor TCP que escucha las conexiones entrantes y maneja los mensajes enviados por otros nodos. Los mensajes se codifican en formato JSON y contienen información sobre los comandos y los jugadores involucrados.

##Funciones principales:

Salto entre nodos: El juego permite a los jugadores realizar saltos entre nodos, simulando su movimiento a través del grafo de nodos. Los jugadores pueden enviar un mensaje "jump" al nodo actual, y el nodo decide si el salto es válido o no, basado en ciertas condiciones y la lógica del juego.

Comunicación entre nodos: Los nodos se comunican entre sí enviando mensajes JSON codificados a través de conexiones TCP. Los mensajes contienen información sobre los comandos, el jugador actual y los nodos vecinos.

Gestión de equipos: Cada jugador está asociado a un equipo y compite contra el equipo contrario para llegar al nodo final. El juego maneja la lógica de los equipos y determina los ganadores y perdedores basándose en las acciones de los jugadores.

Imagen Referencial
![image](https://github.com/GustavoBecerra1109/pcd-TF/assets/54639476/c03646da-8131-4778-92da-c27d57fd04d0)4

##Conclusión:
