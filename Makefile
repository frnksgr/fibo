
IMAGE 			?: frnksgr/fibo-dev
BASEIMAGE 		?: alpine\:3.9

.PHONY: build
build: 
	go build -v ./pkg/...

bin/fibo: build
	go build  -o ./bin/fibo ./cmd/fibo/main.go

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE) --build-arg BASEIMAGE=$(BASEIMAGE) .

.PHONY: docker-build
docker-push:
	docker push -t $(IMAGE)

.PHONY: cf-push
cf-push: bin/fibo
	cf push -f config/cf-manifest.yaml