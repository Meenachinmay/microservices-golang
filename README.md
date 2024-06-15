# Microservices with Golang

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Table of Contents

- [Introduction](#introduction)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Services](#running-the-services)
- [Services](#services)
  - [API Gateway](#api-gateway)
  - [Authentication Service](#authentication-service)
  - [Logger Service](#logger-service)
  - [Mailer Service](#mailer-service)
  - [Listener Service](#listener-service)
- [Technologies Used](#technologies-used)
- [Contributing](#contributing)
- [License](#license)

## Introduction

This repository contains a microservices-based application built with Golang. The application demonstrates how to use microservices to create a scalable and maintainable architecture. It includes services for authentication, logging, mailing, and a listener service that processes events from RabbitMQ.

## Architecture

The architecture is composed of multiple services that communicate with each other through REST APIs and message queues. The services are:

- **API Gateway**: The entry point for all client requests.
- **Authentication Service**: Handles user authentication.
- **Logger Service**: Logs application events.
- **Mailer Service**: Sends emails.
- **Listener Service**: Listens for events from RabbitMQ and processes them.


## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- [Golang](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/Meenachinmay/microservices-golang.git
    cd microservices-golang
    ```

2. Create an `.env` file in the root directory and set the necessary environment variables. An example `.env` file is provided as `.env.example`.

    ```sh
    cp .env.example .env
    ```

### Running the Services

1. Build and run the services using Docker Compose using Makefile:

    ```sh
    make up_build
    ```

2. The services will be available at the following ports:

    - API Gateway: `http://localhost:8080`
    - Authentication Service: `http://localhost:8081`
    - Logger Service: `http://localhost:8082`
    - Mailer Service: `http://localhost:8083`
    - Listener Service: `http://localhost:8084`

## Services

### API Gateway (Broker-service)

The API Gateway is the main entry point for all incoming requests. It routes requests to the appropriate service based on the action specified.

### Authentication Service

The Authentication Service handles user authentication. It provides endpoints for login, signup, and password reset. (as of now just login is there for simulation, I will add more like signup, session, cookies and everything)

### Logger Service

The Logger Service logs application events. It listens for log events from RabbitMQ and stores them in a log file.

### Mailer Service

The Mailer Service sends emails. It listens for mail events from RabbitMQ and sends emails using the configured SMTP server. (I am using mailhog as of now, later I will use sendgrid)

### Listener Service

The Listener Service listens for events from RabbitMQ and processes them. It handles log and mail events, delegating tasks to the appropriate services.

## Technologies Used

- **Golang**: The primary language used for building the services.
- **Docker**: Used for containerizing the services.
- **Docker Compose**: Used for orchestrating the multi-container Docker application.
- **RabbitMQ**: Used for messaging between services.
- **PostgreSQL**: The database used by the Authentication Service.
- **MongoDB**: The database used by the Logger Service.
- **gRPC: Communication between api-gateway and services.

### Command to generate gRPC code 
``` sh 
     protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <name of file>.proto
```

## Contributing

Contributions are welcome! Please read the [CONTRIBUTING.md](./CONTRIBUTING.md) file for more information on how to contribute to this project.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE.md) file for details.
