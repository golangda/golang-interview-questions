package contracts

type Command struct {
	TraceID       string                 `json:"trace_id"`
	CorrelationID string                 `json:"correlation_id"`
	Timestamp     string                 `json:"timestamp"`
	Command       string                 `json:"command"`
	Resource      string                 `json:"resource"`
	Payload       map[string]any         `json:"payload"`
	Metadata      map[string]any         `json:"metadata"`
}

type Ack struct {
	TraceID       string         `json:"trace_id"`
	CorrelationID string         `json:"correlation_id"`
	Timestamp     string         `json:"timestamp"`
	Status        string         `json:"status"`
	Event         string         `json:"event"`
	Payload       map[string]any `json:"payload"`
	Error         *struct {
		Code   string `json:"code"`
		Detail string `json:"detail"`
	} `json:"error,omitempty"`
}