# SAGA Choreography Lab (Go + Kafka + K8s on Minikube)

Spin up a 5‑step choreography saga, simulate Step‑5 failures, and learn real‑world debugging with Prometheus + Grafana + Jaeger.

## Prereqs

- minikube, kubectl, helm, docker
- Go 1.22+ (for local builds if needed)

## 1) Start Minikube

```bash
minikube start --cpus=6 --memory=12288
eval $(minikube docker-env)
```

## 2) Install Kafka + Observability (Helm)

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm repo update

helm install kafka bitnami/kafka --set replicaCount=1 --set zookeeper.enabled=true
helm install kps prometheus-community/kube-prometheus-stack \  --set grafana.service.type=NodePort --set prometheus.service.type=NodePort
helm install jaeger jaegertracing/jaeger \  --set storage.type=memory --set query.service.type=NodePort
```

## 3) Build images (into Minikube’s Docker)

```bash
# from repo root
docker build -t saga/emitter:dev --build-arg CMD=emitter .
for s in step1 step2 step3 step4 step5 dlq-replayer; do
  docker build -t saga/$s:dev --build-arg CMD=$s .
done
```

## 4) Deploy K8s resources

```bash
kubectl apply -f k8s/00-topics-job.yaml
kubectl apply -f k8s/step1.yaml -f k8s/step2.yaml -f k8s/step3.yaml -f k8s/step4.yaml -f k8s/step5.yaml
kubectl apply -f k8s/emitter.yaml -f k8s/dlq-replayer.yaml
kubectl apply -f k8s/10-servicemonitor.yaml
```

## 5) Open Grafana & Jaeger

```bash
minikube service -n default kps-grafana --url
minikube service -n default jaeger-query --url
```

- Grafana default creds: `admin/prom-operator` (or check chart notes).  
- Explore metrics:
  - `saga_retries_total`
  - `saga_step_latency_seconds_bucket`
  - `dlq_messages_total`

In Jaeger, search traces by `x-saga-id` tag or filter by service `saga-step-5`.

## 6) Labs

### Lab A: Retryable failure storms
```bash
kubectl set env deploy/step5 FAIL_MODE=retryable
# Watch retries and consumer lag
kubectl logs deploy/step5 -f
```

### Lab B: Fatal (schema/validation) → DLQ + replay
```bash
kubectl set env deploy/step5 FAIL_MODE=fatal
# Inspect DLQ
kubectl run ktools --image=bitnami/kafka:latest -it --rm -- bash -lc   'kafka-console-consumer.sh --bootstrap-server kafka:9092 --topic saga.dlq --from-beginning --max-messages 1'
# Stop failure, then drain DLQ
kubectl set env deploy/step5 FAIL_MODE=none
# (dlq-replayer will re-emit to x-original-topic or REPLAY_TARGET)
```

### Lab C: Fix & observe recovery
```bash
kubectl set env deploy/step5 FAIL_MODE=none
# Lag should drain, retries drop
```

### Optional: Ordering bug drill
Temporarily build `step3` with a random Kafka key to observe out-of-order effects, then revert.

## 7) Useful commands

```bash
kubectl get pods
kubectl logs deploy/step5 -f

# Kafka consumer-group lag (inside kafka pod)
kubectl exec -it deploy/kafka -- bash -lc   'kafka-consumer-groups.sh --bootstrap-server kafka:9092 --describe --group svc5-group'
```

## 8) Clean up

```bash
kubectl delete -f k8s
helm uninstall jaeger
helm uninstall kps
helm uninstall kafka
minikube delete
```

---

### Notes
- Services expose Prometheus metrics at `:8080/metrics` and send traces to `jaeger-collector:14268`.
- All Kafka messages are keyed by `saga_id` to preserve per-saga ordering.
- DLQ messages include the `x-original-topic` header for safe targeted replay.


## 9) Provision the sample Grafana dashboard

```bash
kubectl apply -f k8s/20-grafana-dashboard.yaml
# Then open Grafana (make grafana) -> Dashboards -> 'SAGA Choreography Lab'
```

## 9) Quick How to Use Note:
```bash
# after `make up`
make dashboard
make grafana  # open the URL
```
In Grafana, look for “SAGA Choreography Lab”. Panels include:
- DLQ messages/min
- Retries by reason (all steps)
- Step-5 latency p95
- Step throughput (per step)
- DLQ rate (by topic)
- Step-5 retry reasons (zoom)

