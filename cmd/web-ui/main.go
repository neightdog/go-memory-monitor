package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

var (
	memoryData = struct {
		Usage  uint64 `json:"usage"`
		Status int32  `json:"status"`
	}{}
	diskData = struct {
		Usage  uint64 `json:"usage"`
		Status int32  `json:"status"`
	}{}
	cpuData = struct {
		Usage  uint64 `json:"usage"`
		Status int32  `json:"status"`
	}{}
	mu sync.RWMutex
)

func main() {
	go subscribeToUsageEvents()

	http.HandleFunc("/api/memory", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		json.NewEncoder(w).Encode(memoryData)
	})
	http.HandleFunc("/api/disk", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		json.NewEncoder(w).Encode(diskData)
	})
	http.HandleFunc("/api/cpu", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		json.NewEncoder(w).Encode(cpuData)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head><title>Go Memory, Disk & CPU Monitor</title></head>
<body>
<h1>System Monitor</h1>
<div id="memory"></div>
<div id="disk"></div>
<div id="cpu"></div>
<script>
async function fetchData() {
  try {
    let mem = await fetch('/api/memory').then(r => r.json());
    let disk = await fetch('/api/disk').then(r => r.json());
    let cpu = await fetch('/api/cpu').then(r => r.json());
    document.getElementById('memory').innerText = 'Memory usage: ' + mem.usage + '% (' + (mem.status === 1 ? 'above' : 'below') + ' 50%)';
    document.getElementById('disk').innerText = 'Disk usage: ' + disk.usage + '% (' + (disk.status === 1 ? 'above' : 'below') + ' 50%)';
    document.getElementById('cpu').innerText = 'CPU usage: ' + cpu.usage + '% (' + (cpu.status === 1 ? 'above' : 'below') + ' 50%)';
  } catch (e) {
    document.getElementById('memory').innerText = 'Error fetching memory data';
    document.getElementById('disk').innerText = 'Error fetching disk data';
    document.getElementById('cpu').innerText = 'Error fetching CPU data';
  }
}
setInterval(fetchData, 2000);
fetchData();
</script>
</body>
</html>
        `))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func subscribeToUsageEvents() {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = ch.QueueBind(q.Name, "", "usage", false, nil)
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", true, true, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for msg := range msgs {

		log.Printf("Received message: %s", string(msg.Body))

		var data map[string]interface{}
		json.Unmarshal(msg.Body, &data)
		mu.Lock()
		switch data["type"] {
		case "memory":
			memoryData.Usage = uint64(data["usage"].(float64))
			memoryData.Status = int32(data["status"].(float64))
		case "disk":
			diskData.Usage = uint64(data["usage"].(float64))
			diskData.Status = int32(data["status"].(float64))
		case "cpu":
			cpuData.Usage = uint64(data["usage"].(float64))
			cpuData.Status = int32(data["status"].(float64))
		}
		mu.Unlock()
	}
}
