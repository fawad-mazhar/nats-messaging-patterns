package publisher

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/fawadmazhar/nats-queue-group/internal/config"
)

type Publisher struct {
	nc *nats.Conn
	wg sync.WaitGroup
}

func NewPublisher() *Publisher {
	nc, err := nats.Connect(config.DefaultNatsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	return &Publisher{
		nc: nc,
	}
}

func (p *Publisher) Run() error {
	defer p.nc.Close()

	// Short delay to ensure workers have time to start and join the queue group
	time.Sleep(500 * time.Millisecond)
	
	p.wg.Add(1)
	go p.publishTasks()
	
	p.wg.Wait()
	return nil
}

func (p *Publisher) Shutdown() {
	// Implementation for graceful shutdown if needed
}

func (p *Publisher) publishTasks() {
	defer p.wg.Done()
	
	log.Println("Starting to publish tasks...")
	
	// Publish tasks
	for i := 1; i <= config.NumTasks; i++ {
		taskMsg := "Task " + strconv.Itoa(i)
		
		// Publish with request (expecting a response)
		response, err := p.nc.Request(
			config.Subject, 
			[]byte(taskMsg), 
			config.TaskResponseTimeout,
		)
		if err != nil {
			log.Printf("Error publishing task %d: %v\n", i, err)
			continue
		}
		
		log.Printf("Published: %s, Response: %s\n", taskMsg, string(response.Data))
		
		// Wait between tasks to allow better distribution
		time.Sleep(config.TaskPublishInterval)
	}
	
	log.Println("Finished publishing tasks")
}

