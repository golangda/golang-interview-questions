package trace

import "context"

type ctxKey string

const TraceIDKey ctxKey = "trace_id"

func WithTraceID(ctx context.Context, id string) context.Context { return context.WithValue(ctx, TraceIDKey, id) }
func GetTraceID(ctx context.Context) (string, bool) {
	v := ctx.Value(TraceIDKey)
	if s, ok := v.(string); ok { return s, true }
	return "", false
}