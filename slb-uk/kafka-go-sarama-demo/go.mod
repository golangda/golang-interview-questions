module example.com/kafka-go-sarama-demo

go 1.22.0

require (
	github.com/IBM/sarama v1.45.0
	github.com/dnwe/otelsarama v0.4.3
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/semconv v1.26.0
)
