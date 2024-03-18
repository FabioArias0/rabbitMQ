package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
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

	// Declarar la cola de entrada
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

	// Conexión a la base de datos PostgreSQL
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:32772/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	// Procesar mensajes recibidos
	for d := range msgs {
		// Decodificar el mensaje
		var message Message
		err := json.Unmarshal(d.Body, &message)
		if err != nil {
			log.Printf("Error al decodificar el mensaje: %v", err)
			continue
		}

		// Guardar el mensaje con la hora en que se recibió en la base de datos PostgreSQL
		_, err = db.Exec("INSERT INTO messages(content, created_at) VALUES($1, $2)", message.Content, message.CreatedAt)
		if err != nil {
			log.Printf("Error ejecutando la consulta SQL: %v", err)
			continue
		}

		// Enviar un mensaje de confirmación a la primera aplicación
		err = ch.Publish(
			"",         // Exchange
			"response", // Key de enrutamiento
			false,      // Mandatorio
			false,      // Inmediato
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("Mensaje procesado y guardado en la base de datos"),
			},
		)
		if err != nil {
			log.Printf("Error enviando el mensaje de confirmación: %v", err)
			continue
		}

		fmt.Println("Mensaje procesado y guardado:", message.Content)
	}
}
