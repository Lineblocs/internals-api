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

Debugging issues by tracking logs

### Configure log channels

There are 4 log channels including console, file, cloudwatch, logstash
Set LOG_DESTINATIONS variable in .env file

ex: export LOG_DESTINATIONS=console,file

## Linting and pre-comit hook

### Go lint
Config .golangci.yaml file to add or remote lint options

### pre-commit hook
Config .pre-commit-config.yaml file to enable or disable pre-commit hook

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
