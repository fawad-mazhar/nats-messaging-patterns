# NATS Queue Groups Example

A demonstration of NATS Queue Groups for load balancing work across multiple subscribers.

## Features

- Multiple workers processing tasks concurrently
- Load balancing via NATS Queue Groups
- Request-Reply pattern for task acknowledgment
- Dedicated publisher for task distribution
- Configurable task processing simulation
- Worker identification and monitoring

## Prerequisites

- Go 1.21+
- Docker (for running NATS server)

## How It Works

1. **Queue Groups**: Workers subscribe to the "tasks" subject as part of the "workers" queue group
2. **Load Balancing**: NATS distributes tasks among available workers
3. **Task Processing**: Each worker:
   - Receives tasks
   - Simulates processing time
   - Sends acknowledgment
4. **Publishing**: A dedicated publisher component:
   - Sends test tasks
   - Receives and logs worker responses
   - Ensures even task distribution

## Quick Start

1. Start the NATS server:
   ```bash
   make docker-nats
   ```

2. Build the components:
   ```bash
   make build
   ```

3. Run multiple workers:
   ```bash
   # In separate terminals:
   make run-worker ID=1
   make run-worker ID=2
   make run-worker ID=3
   
   # Or run all at once:
   make run-workers N=3
   ```

4. Run the publisher in a separate terminal:
   ```bash
   make run-publisher
   ```

5. For a complete demo (runs everything):
   ```bash
   make start-demo
   ```

## Configuration

The following can be configured in `internal/config/config.go`:

- Task processing time range
- Number of tasks to publish
- Publishing interval
- Response timeout
- NATS connection settings

## Example Output

```
Worker 1 started and joined 'workers' queue group
Worker 2 started and joined 'workers' queue group
Worker 3 started and joined 'workers' queue group
Starting to publish tasks...
Worker 2 processed message: Task 1 (took 234ms)
Published: Task 1, Response: Processed by worker 2
Worker 3 processed message: Task 2 (took 156ms)
Published: Task 2, Response: Processed by worker 3
Worker 1 processed message: Task 3 (took 292ms)
Published: Task 3, Response: Processed by worker 1
...
```

## Cleanup

To stop everything:

```bash
make stop-nats
make clean
```
