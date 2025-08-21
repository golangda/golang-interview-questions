SHELL := /bin/bash

SERVICES := emitter step1 step2 step3 step4 step5 dlq-replayer

.PHONY: help
help:
	@echo "make up        - install helm deps, build images, apply k8s manifests"
	@echo "make down      - delete k8s resources"
	@echo "make images    - build docker images into minikube daemon"
	@echo "make topics    - create kafka topics (job)"
	@echo "make grafana   - open grafana URL via minikube"
	@echo "make jaeger    - open jaeger URL via minikube"

.PHONY: up
up: helm images topics deploy

.PHONY: helm
helm:
	helm repo add bitnami https://charts.bitnami.com/bitnami || true
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts || true
	helm repo add jaegertracing https://jaegertracing.github.io/helm-charts || true
	helm repo update
	helm install kafka bitnami/kafka --set replicaCount=1 --set zookeeper.enabled=true || true
	helm install kps prometheus-community/kube-prometheus-stack --set grafana.service.type=NodePort --set prometheus.service.type=NodePort || true
	helm install jaeger jaegertracing/jaeger --set storage.type=memory --set query.service.type=NodePort || true

.PHONY: images
images:
	@eval $$(minikube docker-env); \	for s in $(SERVICES); do \		echo "==> Building $$s"; \		docker build -t saga/$$s:dev --build-arg CMD=$$s . || exit 1; \	done

.PHONY: topics
topics:
	kubectl apply -f k8s/00-topics-job.yaml

.PHONY: deploy
deploy:
	kubectl apply -f k8s/step1.yaml -f k8s/step2.yaml -f k8s/step3.yaml -f k8s/step4.yaml -f k8s/step5.yaml
	kubectl apply -f k8s/emitter.yaml -f k8s/dlq-replayer.yaml
	kubectl apply -f k8s/10-servicemonitor.yaml

.PHONY: down
down:
	kubectl delete -f k8s || true
	-helm uninstall jaeger
	-helm uninstall kps
	-helm uninstall kafka

.PHONY: grafana
grafana:
	minikube service -n default kps-grafana --url

.PHONY: jaeger
jaeger:
	minikube service -n default jaeger-query --url


.PHONY: dashboard
dashboard:
	kubectl apply -f k8s/20-grafana-dashboard.yaml


.PHONY: loki
loki:
	helm repo add grafana https://grafana.github.io/helm-charts || true
	helm repo update
	helm install loki grafana/loki-stack --set grafana.enabled=false --set prometheus.enabled=false --set fluent-bit.enabled=false --set promtail.enabled=true || true

.PHONY: datasources
datasources:
	kubectl apply -f k8s/30-grafana-datasources.yaml

.PHONY: logs-dashboard
logs-dashboard:
	kubectl apply -f k8s/21-grafana-logs-traces.yaml
