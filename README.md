# simple-load-balancer

A simple load balancer implemented in Go that routes incoming HTTP requests to multiple backend servers. This project serves as an educational example and is not production-ready.

## Features

- Least Connections Algorithm for Load Balancing
- Random Algorithm for Load Balancing
- Weighted Random Algorithm for Load Balancing
- Dynamic Weighted Random Algorithm for Load Balancing
- Least Response Time Algorithm for Load Balancing
- Atomic updates for active connection counts
- Simple HTTP backend servers for demonstration
- Supports multiple backend servers
- Thread-safe operations

## Prerequisites

- Go 1.20.5

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/araujo88/simple-load-balancer
cd simple-load-balancer
```

### Start the Backend Servers

Navigate to the `/server` directory:

```bash
cd server
```

Run the example servers on ports 8081 and 8082:

```bash
go run main.go 8081
go run main.go 8082
```

### Start the Load Balancer

Navigate back to the root directory and start the load balancer:

```bash
cd ..
go run . -algorithm=<algorithm-type>
```

Your load balancer should now be running and forwarding incoming HTTP requests to the backend servers.

## Structure

- `main.go`: The load balancer code.
- `/server/main.go`: Example HTTP server code.

## Algorithm

This load balancer supports multiple algorithms for distributing incoming HTTP requests among backend servers. Below are the algorithms implemented:

### Least Connections

In the Least Connections method, incoming requests are routed to the server with the fewest active connections. This helps ensure a more equitable distribution of load.

### Random

The Random method randomly selects a backend server for each incoming request. All servers have an equal chance of being chosen, regardless of their current load or performance.

### Weighted Random

The Weighted Random method assigns a static weight to each backend server. The probability of selecting a particular server is proportional to its weight. Servers with higher weights will receive more requests than those with lower weights.

### Dynamic Weighted Random

This is an extension of the Weighted Random algorithm. In this method, the weight of each server is determined dynamically based on the inverse of its current number of active connections. This allows the load balancer to adapt to the real-time load on each server.

### Least Response Time

The Least Response Time algorithm selects the server that has the lowest response time for a new request. Implementing this algorithm would involve measuring the response time of each server and directing incoming requests to the server with the least recent response time.

## How to Contribute

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

## License

This project is open source and available under the [GNU GENERAL PUBLIC LICENSE](LICENSE).

## TODOs

 - Incorporate moving averages to measure response time
 - Handle failed requests
 - Incorporate health checks
 - Measure CPU and disk usage
 - Metrics and logging
 - Rate limiting
 - Connection draining
 - Circuit breaker pattern
 - Sticky sessions
 - DDoS Protection
 - Response Caching
 - IPv6 Support
 - WebSockets Support
 - API for Management
 - Automated Tests
