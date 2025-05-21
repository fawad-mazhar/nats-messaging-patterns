# NATS Messaging Patterns

This repository contains implementations of different messaging patterns using NATS and JetStream. Each subdirectory is a complete, self-contained project demonstrating a specific messaging pattern.

## Projects

### [nats-pubsub](./nats-pubsub)
Implementation of the Publish-Subscribe pattern using NATS JetStream.

**Features:**
- Publisher for sending messages to a NATS stream
- Subscriber with concurrent message processing
- Real-time monitoring capabilities
- Configurable stream settings

### [nats-reqrep](./nats-reqrep)
Implementation of the Request-Reply pattern using NATS JetStream.

**Features:**
- Server that processes requests and sends responses
- Client that sends requests and waits for responses
- JetStream persistence for reliability
- Proper timeout handling

### [nats-queue-group](./nats-queue-group)
Implementation of the Queue Group (Competing Consumers) pattern using NATS JetStream.

**Features:**
- Multiple consumers in a queue group share the message load
- Ensures each message is processed by only one consumer in the group
- Demonstrates horizontal scalability for message processing
- Configurable group and stream settings

## Getting Started

Each project contains its own README with detailed instructions, but generally:

1. Make sure you have the prerequisites:
   - Go 1.21+
   - Docker
   - Make

2. Choose a pattern and navigate to its directory:
   ```bash
   cd nats-pubsub   # or nats-reqrep or nats-queue-group
   ```

3. Follow the project-specific README for detailed instructions.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
