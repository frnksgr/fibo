
IMAGE 			= frnksgr/fibo
GCR_IMAGE		= gcr.io/sap-cp-gke/fibo
# only us-central1 has beta run
# region need to be set as default region in gcloud config 
#REGION			= us-central1
BASEIMAGE 		= alpine:3.9

k8s-domain 		= default.example.com
knative-gw		= $(shell scripts/get-gateway.sh istio-system istio-ingressgateway)
nginx-gw		= $(shell scripts/get-gateway.sh nginx nginx-ingress-controller)
gcp-token		= $(shell scripts/get-gcp-token.sh)

# wow, simple self documenting makefile
# see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.DEFAULT_GOAL := help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build local to create bin/fibo executable
	go build -v ./pkg/...
	go build -o ./bin/fibo ./cmd/fibo/...

.PHONY: docker-build
docker-build: ## build docker image
	docker build -t $(IMAGE) --build-arg BASEIMAGE=$(BASEIMAGE) .

.PHONY: docker-push
docker-push: docker-build ## push docker-image
	docker push $(IMAGE)

.PHONY: gcr-build
gcr-build: ## build container image for gcr
	docker build -t $(GCR_IMAGE) --build-arg BASEIMAGE=$(BASEIMAGE) .

.PHONY: gcr-push
gcr-push: gcr-build ## push container-image to gcr
	docker push $(GCR_IMAGE)

#curl -H "Authorization: Bearer $(gcloud config config-helper --format 'value(credential.id_token)')" <service url>
.PHONY: gcp-run-deploy
gcp-run-deploy: gcr-push ## deploy to gcp-run
	# yet only available in us-central1
	gcloud beta run deploy fibo --image=$(GCR_IMAGE)
	@echo service URI:
	@gcloud beta run routes describe fibo \
		--format 'value(status.address.hostname)'
	@echo curl -H \"Authorization Bearer \<token\>\" \<service URI\>

gcp-run-deploy-only:
	# yet only available in us-central1
	gcloud beta run deploy fibo --image=$(GCR_IMAGE)
	@echo service URI:
	@gcloud beta run routes describe fibo \
		--format 'value(status.address.hostname)'
	@echo curl -H \"Authorization Bearer \<token\>\" \<service URI\>

.PHONY: gke-run-deploy
gke-run-deploy: gcr-push ## deploy to gke-run
	@echo not implemented

.PHONY: cf-push 
cf-push: build ## push application to CF
	cf push -f config/cf-manifest.yaml

.PHONY: knative-push 
knative-push: docker-push ## push application to Knative on current K8S cluster
	kubectl apply -f config/knative.yaml
	@echo to call: curl -H \"Host: fibo.$(k8s-domain)\" http://$(knative-gw)/	

.PHONY: clean
clean: ## clean up
	go clean -i ./...
	rm -f bin/fibo
