# Go example web http server

This application is a simple demonstration of how to get a basic http server up and running using Go.
It showcases some best practices of running microservices in Kubernetes, such as:

- graceful shutdown
- metrics
- health checks
- config from env variables
- structured logging

## Available endpoints

```sh
GET /health     # returns http code 200 when the pod is up and running
GET /ready      # returns http code 200 when the app can communicate with the DB otherwise it returns http code 503
GET /version    # returns the application build version
GET /metrics    # exposes metrics
```

## Run with docker

In order to run the application locally with docker, first build the image:

```sh
make build-image
```

Then run with:

```sh
make docker-run
```

You may access the app locally:

```bash
curl http://localhost:8024/ready
```

## Deploy on Minikube

You may deploy the app to a local kubernetes cluster (e.g. minikube) using **helm**.

### Prerequisites

> Minikube \
Helm

#### Prepare Minikube

Start Minikube local cluster:

```sh
minikube start
```

Then run:

```sh
make deploy
```

This will deploy the sample webapp in the current kubernetes namespace.

Also, you may enable minikube metrics server so as horizontal pod autoscaler to be able to work (disabled by default):

```sh
minikube addons enable metrics-server
```

In order to access the app, run:

```sh
minikube service service go-server-bootstrap
```

or via port forwarding:

```sh
kubectl port-forward services/go-server-bootstrap 8024:8024
```

### TODO

Write tests.
