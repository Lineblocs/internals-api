## Prerequisites

Minimum software required to run the service:
* [go](https://go.dev/doc/install)

## Clone repository

```bash
git clone https://github.com/Lineblocs/internals-api.git
```

## Structure of Code

internals-api uses echo framework which is high performance, extensible, minimalist Go framework.
1. router
   initiate echo and set basic configuration.
2. handler
   configure routings and bind all services including admin, call, carrier, debit, fax, logger, recording, user
3. store
   each services are declared with interface in their own package and implemenations are defined in store package.
4. model
   includes all models which can be used in services.


## Running Newsman tests

### Install Node

Newman is built on Node.js. To run Newman, make sure you have Node.js installed.

```bash
sudo apt update
sudo apt install nodejs
node -v
```

### Install Newsman

```bash
npm install -g newman
```

### Test API endpoints

Running Newsman via Postman collection

```bash
newman run https://api.getpostman.com/collections/myPostmanCollectionUid?apikey=myPostmanApiKey
collectionUid = 25298469-50f49cbd-43c7-459a-9052-706e8d7c002f
apiKey = PMAK-63c11b3766698c5953a3333b-bc04bedd1b5fb17cba12cae7c1de9018ec
```

## Debugging

### Configure HTTP and HTTPS PORT
API uses default 80 port for HTTP and 443 for HTTPS
You Configure ports with HTTP_PORT and HTTPS_PORT environment variable.
export HTTP_PORT=80
export HTTPS_PORT=443

### Configure log channels
Debugging issues by tracking logs

There are 4 log channels including console, file, cloudwatch, logstash
Set LOG_DESTINATIONS variable in .env file

ex: export LOG_DESTINATIONS=console,file

## Reliability and Fault Tolerance

The service includes built-in safeguards to reduce cascading failures:

1. Panic recovery middleware
Handles panics in HTTP handlers, logs stack traces with structured fields, and returns HTTP 500.

2. Graceful shutdown
On SIGINT/SIGTERM the server stops accepting new requests and drains in-flight work before exit.

3. Context-based query deadlines
Database queries use context timeouts through the shared MySQL wrapper.

4. Retry with backoff for downstream webhooks
SIP webhook delivery uses bounded retries with per-request timeout.

### Resilience Environment Variables

Set these to tune fault-tolerance behavior:

```bash
# Shutdown drain window (Go duration format)
export GRACEFUL_SHUTDOWN_TIMEOUT=10s

# DB timeout per operation (milliseconds)
export DB_QUERY_TIMEOUT_MS=5000

# HTTP client timeout for downstream calls (milliseconds)
export HTTP_CLIENT_TIMEOUT_MS=5000

# Retry count for downstream HTTP calls
export HTTP_CLIENT_RETRIES=3
```

## Linting and pre-comit hook

### Go lint
Config .golangci.yaml file to add or remote lint options

### pre-commit hook
Config .pre-commit-config.yaml file to enable or disable pre-commit hook

## load test with k6

### Install k6
```bash
sudo apt-get update
sudo apt-get upgrade
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

### Run script with k6
```bash
k6 run loadtest.js
```

## Deploy

### Deploy Steps
1. Install Docker on the machines you want to use it;
2. Set up a registry at Docker Hub;
3. Initiate Docker build to create your Docker Image;
4. Set up your ’Dockerized‘ machines;
5. Deploy your built image or application.

### Deploy Command

```bash
docker build -t internals-api
```
