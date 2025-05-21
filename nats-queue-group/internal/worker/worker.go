package worker

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/fawadmazhar/nats-queue-group/internal/config"
)

type Worker struct {
	id string
	nc *nats.Conn
	sub *nats.Subscription
	wg  sync.WaitGroup
}

func NewWorker(id string) *Worker {
	nc, err := nats.Connect(config.DefaultNatsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	return &Worker{
		id: id,
		nc: nc,
	}
}

func (w *Worker) Run() error {
	defer w.nc.Close()

	// Small delay to ensure proper queue group formation
	time.Sleep(100 * time.Millisecond)

	// Subscribe to tasks subject as part of workers queue group
	sub, err := w.nc.QueueSubscribe(config.Subject, config.QueueGroup, w.processTask)
	if err != nil {
		return fmt.Errorf("error subscribing: %w", err)
	}
	w.sub = sub

	log.Printf("Worker %s started and joined '%s' queue group\n", w.id, config.QueueGroup)

	// Wait indefinitely (until shutdown is called)
	select {}
	
	return nil
}

func (w *Worker) Shutdown() {
	if w.sub != nil {
		w.sub.Unsubscribe()
	}
}

func (w *Worker) processTask(msg *nats.Msg) {
	// Simulate processing time
	processingTime := rand.Intn(
		int(config.TaskProcessingMaxTime-config.TaskProcessingMinTime),
	) + int(config.TaskProcessingMinTime)
	time.Sleep(time.Duration(processingTime))
	
	log.Printf("Worker %s processed message: %s (took %dms)\n", 
		w.id, string(msg.Data), processingTime/int(time.Millisecond))
	
	// Send an acknowledgment
	msg.Respond([]byte("Processed by worker " + w.id))
}

