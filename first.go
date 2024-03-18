package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	//Iniciar Docker con RabbitMQ
	//docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
	//docker network create --subnet=192.168.0.0/16 cluster-network

	//node1
	//docker run -d -h node1.rabbit --net cluster-network --ip 192.168.0.10 --add-host node2.rabbit:192.168.0.11 --add-host node3.rabbit:192.168.0.12 -p "4369:4369" -p "5672:5672" -p "15672:15672" -p "25672:25672" -p "35672:35672" -e "RABBITMQ_USE_LONGNAME=true" -e RABBITMQ_ERLANG_COOKIE="cookie" rabbitmq:3-management

	//node2
	//docker run -d -h node2.rabbit --net cluster-network --ip 192.168.0.11 --name rabbitNode2 --add-host node1.rabbit:192.168.0.10 --add-host node3.rabbit:192.168.0.12 -p "4370:4369" -p "5673:5672" -p "15673:15672" -p "25673:25672" -p "35673:35672" -e "RABBITMQ_USE_LONGNAME=true" -e RABBITMQ_ERLANG_COOKIE="cookie" rabbitmq:3-management

	//Node3
	//docker run -d -h node3.rabbit --net cluster-network --ip 192.168.0.12 --name rabbitNode3 --add-host node1.rabbit:192.168.0.10 --add-host node2.rabbit:192.168.0.11 -p "4371:4369" -p "5674:5672" -p "15674:15672" -p "25674:25672" -p "35674:35672" -e "RABBITMQ_USE_LONGNAME=true" -e RABBITMQ_ERLANG_COOKIE="cookie"   rabbitmq:3-management

	//stop de node >>> docker exec rabbitNode2 rabbitmqctl join_cluster rabbit@node1.rabbit

	//join the nodes >>> docker exec rabbitNode3 rabbitmqctl join_cluster rabbit@node1.rabbit
	//user; guest, password: guest

	//mysql >> docker run -d -p 3307:3306 --name mysql-container -e MYSQL_ROOT_PASSWORD=pass123 mysql:latest

	// URL de conexión a RabbitMQ
	url := "amqp://guest:guest@localhost:5672/"

	// Conectarse a RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Conexión establecida con RabbitMQ")

}
