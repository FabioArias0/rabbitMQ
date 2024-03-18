package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Lista de endpoints de RabbitMQ
	endpoints := []string{"amqp://guest:guest@localhost:5672/", "amqp://guest:guest@localhost:5673/", "amqp://guest:guest@localhost:5674/"}

	// Intentar conectar a los endpoints hasta que uno sea exitoso
	var conn *amqp.Connection
	var err error
	for _, endpoint := range endpoints {
		conn, err = amqp.Dial(endpoint)
		if err == nil {
			break // Se conectó exitosamente, sal del bucle
		}
	}

	// Comprobar si hubo un error al conectar
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Crear un canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creando el canal: %v", err)
	}
	defer ch.Close()

	// Declarar la cola
	q, err := ch.QueueDeclare(
		"Cola", // Nombre de la cola
		false,  // Durabilidad
		false,  // Exclusividad
		false,  // Autoeliminación
		false,  // No-Wait
		nil,    // Argumentos adicionales
	)
	if err != nil {
		log.Fatalf("Error declarando la cola: %v", err)
	}

	// Consumir mensajes de la cola
	msgs, err := ch.Consume(
		q.Name, // Nombre de la cola
		"",     // Etiqueta del consumidor
		true,   // Auto-Ack
		false,  // Exclusivo
		false,  // No-Local
		false,  // No-Wait
		nil,    // Argumentos adicionales
	)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
	}

	// Procesar mensajes recibidos
	for d := range msgs {
		// Mostrar el mensaje recibido
		fmt.Printf("Mensaje recibido: %s\n", d.Body)
	}
}
