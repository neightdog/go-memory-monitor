package main

import (
	"encoding/json"
	"log"
	"time"

	"go-memory-monitor/internal/monitor"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

func main() {
	svc := monitor.NewMemoryService()
	svc.Start()

	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()
	err = ch.ExchangeDeclare("usage", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		usage := svc.GetUsage()
		status := svc.GetStatus()
		msg := map[string]interface{}{
			"type":   "memory",
			"usage":  usage,
			"status": status,
			"time":   time.Now().Format(time.RFC3339),
		}
		body, _ := json.Marshal(msg)
		ch.Publish("usage", "", false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		time.Sleep(2 * time.Second)
	}
}
