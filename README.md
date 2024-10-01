# Webhook Adapter

Webhook Adapter is a Go-based application that serves as an intermediary for handling webhook requests and responses. It uses NATS as a message broker and provides a REST API for interaction.

## Features

- REST API for receiving webhook requests
- NATS integration for message brokering
- Configurable via environment variables and YAML files
- Logging with request ID tracking for easy debugging
- Docker support for easy deployment

## Prerequisites

- Go 1.22 or higher
- NATS server
- Docker with Docker Compose (optional, for e2e tests)

## Installation

1. Clone the repository:
   ```sh
   git clone git@github.com:KubeLambda/kl-webhook-adapter.git
   ```

2. Navigate to the project directory:
   ```sh
   cd webhook-adapter
   ```

3. Install dependencies:
   ```sh
   go mod tidy
   ```

## Configuration

The application can be configured using environment variables or a YAML configuration file. For production, use the configs/production.yaml file:

```yaml
name: webhook-adapter
credentials:
  key: _
  secret: _
server:
  addr: 0.0.0.0
  port: 3001
broker:
  addr: 0.0.0.0
  port: 4222
  stream: 'request-response'
credentials:
  key: _
  secret: _
```

For development, you can use the configs/local.yaml.

Environment variables can override these settings:

- WEBHOOK_ADAPTER_CREDENTIALS_KEY
- WEBHOOK_ADAPTER_CREDENTIALS_SECRET
- WEBHOOK_ADAPTER_SERVER_ADDR
- WEBHOOK_ADAPTER_SERVER_PORT
- WEBHOOK_ADAPTER_BROKER_ADDR
- WEBHOOK_ADAPTER_BROKER_PORT

## Building and Running

To build the application:

```sh
make build
```

To run the application:

```sh
make run deployment=local
```


## Access REST API

Generated application uses REST protocol to store and fetch address book records.
Once you have the application launched, you can perform HTTP calls to test REST APIs exposed by API server.

Please note that each HTTP response contains **X-Request-Id** header with value that is displayed with application logs (as **requestId** field). It helps you to troubleshoot application, because logger provided with generated code prints request id with every log line.

### API Usage

#### Get service version 

Request:

```sh
curl --location 'http://localhost:8080/api/version'
```
Response (could be slightly different):

```json
{
  "service": "rest-net/http",
  "version": "0.1.0",
  "build": "1"
}
```

#### Send request

Request:

```sh
curl --data '{"name":"bob"}' --header 'Content-Type: application/json' http://0.0.0.0:3001/api/adapter
'
```

Response:

```json
```

### Logging

Each HTTP request returns `X-Request-Id` header as part of response. This `X-Request-Id` is always unique, unless you specify it explicitly as part of request. What makes it useful is that each application log line contains `{requestId="...."}` tag, and it matches `X-Request-Id` value. It makes debugging code much easier because you can filter logs scoped to specific request.

## Docker

To run the application in a Docker container, use the following command:

```sh
docker build -t webhook-adapter .
docker run -p 3001:3001 webhook-adapter
```

## Testing

To run unit tests:

```sh
make test
```

For test coverage:

```sh
make test-coverage
```

To run e2e tests:

```sh
make test-e2e
```

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.


## Contributing


Contributions are welcome! Please feel free to submit a Pull Request.