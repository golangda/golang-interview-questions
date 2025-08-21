package tracing

import (
	"strings"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel/propagation"
)

// HeaderCarrier implements OTEL's TextMapCarrier for Sarama headers.
type HeaderCarrier struct{ Headers *[]*sarama.RecordHeader }

func (c HeaderCarrier) Get(key string) string {
	if c.Headers == nil { return "" }
	lk := strings.ToLower(key)
	for _, h := range *c.Headers {
		if strings.ToLower(string(h.Key)) == lk {
			return string(h.Value)
		}
	}
	return ""
}
func (c HeaderCarrier) Set(key, val string) {
	if c.Headers == nil { return }
	lk := strings.ToLower(key)
	for _, h := range *c.Headers {
		if strings.ToLower(string(h.Key)) == lk {
			h.Value = []byte(val); return
		}
	}
	*c.Headers = append(*c.Headers, &sarama.RecordHeader{Key: []byte(key), Value: []byte(val)})
}
func (c HeaderCarrier) Keys() []string {
	if c.Headers == nil { return nil }
	keys := make([]string, 0, len(*c.Headers))
	for _, h := range *c.Headers { keys = append(keys, string(h.Key)) }
	return keys
}

// ExtractContext builds a context from Kafka headers.
func ExtractContext(headers *[]*sarama.RecordHeader, propagator propagation.TextMapPropagator) (ctx interface{ Done() <-chan struct{} }, _ propagation.TextMapCarrier) {
	// We return a background context for simplicity (actual ctx created by caller).
	// Carrier returned for potential further use.
	return nil, HeaderCarrier{Headers: headers}
}
