APP?=apisvc
CONSUMER?=consumersvc
REGISTRY?=
IMAGE_TAG?=local

.PHONY: build build-apisvc build-consumersvc build-mysql docker minikube-load \
        k8s-apply dev-up dev-down test lint migrate logs-apisvc logs-consumersvc pf-apisvc

# --- Build Go binaries locally (useful for unit tests) ---
build:
	go build -o bin/$(APP) ./cmd/$(APP)
	go build -o bin/$(CONSUMER) ./cmd/$(CONSUMER)

# --- Build OCI images directly from service Dockerfiles ---
build-apisvc:
	docker build -t $(REGISTRY)$(APP):$(IMAGE_TAG) -f cmd/$(APP)/Dockerfile .

build-consumersvc:
	docker build -t $(REGISTRY)$(CONSUMER):$(IMAGE_TAG) -f cmd/$(CONSUMER)/Dockerfile .

build-mysql:
	docker build -t $(REGISTRY)mysql:$(IMAGE_TAG) -f deploy/mysql/Dockerfile .

# Build all images (apisvc, consumersvc, and optional mysql)
docker: build-apisvc build-consumersvc build-mysql

# Load images into Minikube (if using Docker driver, this may be a no-op)
minikube-load: docker
	minikube image load $(REGISTRY)$(APP):$(IMAGE_TAG) || minikube image load $(APP):$(IMAGE_TAG)
	minikube image load $(REGISTRY)$(CONSUMER):$(IMAGE_TAG) || minikube image load $(CONSUMER):$(IMAGE_TAG)
	- minikube image load $(REGISTRY)mysql:$(IMAGE_TAG) || true

# Apply Kubernetes manifests
k8s-apply:
	kubectl apply -f k8s/mysql.yaml
	kubectl apply -f k8s/kafka.yaml
	kubectl apply -f k8s/$(APP).yaml
	kubectl apply -f k8s/$(CONSUMER).yaml

# Run SQL migration inside the MySQL pod (expects MYSQL_ROOT_PASSWORD in the pod env)
migrate:
	# apply SQL migration into the mysql pod
	POD=$$(kubectl get po -l app=mysql -o jsonpath='{.items[0].metadata.name}'); \
	kubectl cp migrations/0001_init.sql $$POD:/tmp/0001.sql; \
	kubectl exec $$POD -- sh -c "mysql -uroot -p$${MYSQL_ROOT_PASSWORD} app < /tmp/0001.sql"

# Port-forward API service locally
pf-apisvc:
	kubectl port-forward svc/$(APP) 8080:80

# Logs
logs-apisvc:
	kubectl logs -l app=$(APP) -f --tail=200

logs-consumersvc:
	kubectl logs -l app=$(CONSUMER) -f --tail=200

# One-shot dev bring-up / tear-down
dev-up: docker minikube-load k8s-apply

dev-down:
	-kubectl delete -f k8s/$(CONSUMER).yaml
	-kubectl delete -f k8s/$(APP).yaml
	-kubectl delete -f k8s/mysql.yaml
	-kubectl delete -f k8s/kafka.yaml

# Tests & lint
test:
	go test ./...

lint:
	golangci-lint run || true
