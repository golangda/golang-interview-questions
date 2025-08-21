package retry

import "time"

const (
	HeaderAttempt = "x-retry-attempt"
	HeaderError   = "x-error"
)

type Stage struct {
	Topic string
	Delay time.Duration
}

var Stages = []Stage{
	{Topic: "events.v1.retry.5s",  Delay: 5 * time.Second},
	{Topic: "events.v1.retry.30s", Delay: 30 * time.Second},
	{Topic: "events.v1.retry.2m",  Delay: 2 * time.Minute},
}

func Next(attempt int) (Stage, bool) {
	if attempt < len(Stages) {
		return Stages[attempt], true
	}
	return Stage{}, false
}
