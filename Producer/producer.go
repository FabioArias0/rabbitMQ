package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// Message representa la estructura del mensaje
type Message struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// Establecer la conexión a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	// Escuchar la entrada del usuario para detener el proceso
	go func() {
		log.Println("Presiona [Enter] para detener el proceso...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		log.Println("Proceso detenido.")
		os.Exit(0)
	}()

	for {
		// Mensaje a enviar
		message := Message{
			Content:   "Mensaje de prueba",
			CreatedAt: time.Now(),
		}

		// Convertir mensaje a JSON
		body, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("Error al codificar el mensaje como JSON: %v", err)
		}

		// Publicar el mensaje
		err = ch.Publish(
			"",     // Exchange
			q.Name, // Key de enrutamiento
			false,  // Mandatorio
			false,  // Inmediato
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Fatalf("Error publicando el mensaje: %v", err)
		}

		log.Println("Mensaje enviado:", string(body))

		// Esperar un segundo antes de enviar el siguiente mensaje
		time.Sleep(time.Second)
	}
}
