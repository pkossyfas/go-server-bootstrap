VERSION?=$(shell git describe --tags)
APP_IMAGE=go-server-bootstrap

build-image:
	docker build --build-arg VERSION=$(VERSION) -t $(APP_IMAGE):$(VERSION) .

docker-run:
	docker run --rm -p 8024:8024 $(APP_IMAGE):$(VERSION)

push-minikube:
	minikube image load $(APP_IMAGE):$(VERSION)

helm-deploy:
	helm upgrade --install --wait --set image.repository=$(APP_IMAGE) --set image.tag=$(VERSION) go-server-bootstrap ./helm-chart/go-server-bootstrap 

deploy: build-image push-minikube helm-deploy
