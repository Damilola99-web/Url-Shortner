# URL Shortener Service

This Go code is part of a URL shortening service. It handles the URL shortening request and returns a response with the shortened URL.

## Technologies

-   **Go**: The service is written in Go, a statically typed, compiled language that is efficient and excellent for web services.
-   **Fiber**: Fiber is a web framework built on top of Fasthttp, the fastest HTTP engine for Go. It's designed to ease things up for fast development with zero memory allocation and performance in mind.
-   **Redis**: Redis is used for rate limiting. It's an open-source, in-memory data structure store, used as a database, cache, and message broker.
-   **Docker**: Docker is used to containerize the service. It allows you to package up an application with all of the parts it needs, such as libraries and other dependencies, and ship it all out as one package.
-   **Docker Compose**: Docker Compose is used to define and run multi-container Docker applications. It uses YAML files to configure the application's services and performs the creation and start-up process of all the containers with a single command.

## Installation

1. Install Docker and Docker Compose on your machine.
2. Clone the repository.
3. Navigate to the project directory.
4. Build the Docker image using the command `docker build -t url-shortener .`
5. Run the Docker container using the command `docker run -p 3000:3000 url-shortener`

## Usage

1. Start the service using Docker Compose with the command `docker-compose up`.
2. Send a POST request to `http://localhost:3000/shorten` with a JSON body containing the `URL` and `Expiry`.
3. The service will return a JSON response with the original `URL`, the `CustomShort` URL, the `Expiry`, and the `XRateRemaining`.

## Docker Compose and Redis

The Docker Compose file defines two services: the URL shortener service and the Redis service. The URL shortener service depends on the Redis service.

The Redis service is used for rate limiting. Each IP address is allowed a certain number of requests per time period. The number of remaining requests is stored in the Redis database and is decreased each time a request is made.

The Docker Compose file also defines a network that the two services are a part of. This allows the services to communicate with each other.

To start the services, use the command `docker-compose up`. To stop the services, use the command `docker-compose down`.
