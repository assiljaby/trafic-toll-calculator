# 🚚 Traffic Toll Calculator

This project is a microservices-based system designed to calculate traffic tolls for trucks, with potential applications in broader logistics operations. The system leverages modern technologies to ensure scalability, efficiency, and real-time processing.

## Key Features
- **Microservices Architecture:** Modular design with separate services for toll calculation and invoicing.
- **Golang:** High-performance backend services developed in Go, ensuring low latency and efficient concurrency.
- **Kafka:** Real-time message streaming with Apache Kafka for seamless communication between services.
- **gRPC:** Fast and reliable communication between microservices using gRPC.
- **Scalable and Resilient:** Built with scalability in mind to handle high traffic loads and ensure fault tolerance.

## Technologies Used
- **Golang:** Core backend services
- **Apache Kafka:** Messaging and event streaming
- **gRPC:** Inter-service communication
- **Docker:** Containerization for easy deployment and management

## Use Cases
- Real-time toll calculation for trucking fleets
- Integration with logistics platforms for automated toll management
- Scalable solution for managing tolls across multiple regions or countries

## Getting Started
1. **Clone the repository:**
    ```bash
    git clone https://github.com/assiljaby/trafic-toll-calculator
    

2. **Run Kafka Container:**
    ```bash
    docker compose up -d
    
    or
    
    docker run -d -p 9092:9092 --name broker apache/kafka:latest

    
**Work in progress...**

## TODO
- [x] Simulate the OBU
- [x] Implement data receiver
- [x] Implement logger middleware for the receiver
- [x] Implement Distancec calculator
- [x] Implement logger middleware for the calculator
- [x] Implement invoice aggregator
- [x] Implement logger middleware for the invoice aggregator
- [x] Implement HTTP transport
- [ ] Implement GRPC and proto buffers