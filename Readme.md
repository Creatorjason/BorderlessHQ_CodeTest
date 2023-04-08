# BoderlessHQ Golang Microservice with gRPC, Protobuf, MongoDB and NATS using Docker Compose

This is the solution to the code test from BorderlessHQ .

## Prerequisites

- Docker installed on your local machine
- Docker Compose installed on your local machine

## Getting Started

1. Clone this repository to your local machine.

2. Make sure you have Docker and Docker Compose installed.

3. Update the Golang application code to listen on the appropriate port (e.g., 9091) in the `main.go` file.

4. Update the MongoDB and NATS connection configurations in the Golang application code, if necessary, to match the service names defined in the Docker Compose file (`mongo` and `nats`).

5. Build the Docker containers using Docker Compose with the following command:
   
         docker-compose build


6. Start the Docker containers with the following command:
        
        docker-compose up

7. Access the Golang microservice locally at `http://localhost:9091` in your web browser or through a REST client.

8. Stop the Docker containers with the following command:
    
        docker-compose down

## Docker Compose Configuration

The Docker Compose file (`docker-compose.yml`) in this project defines three services:

- `borderlessHQ_service`: The Golang microservice container that is built from the current directory, exposing port 9091 on the host machine.
- `mongo`: The MongoDB container that uses the official MongoDB image from Docker Hub, exposing port 27017 on the host machine for MongoDB connections.
- `nats`: The NATS container that uses the official NATS image from Docker Hub, exposing port 4222 on the host machine for NATS connections.



