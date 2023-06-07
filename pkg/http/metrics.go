package http

import (
    "time"
)

type Metrics interface {
    IncRequestsCount()
    IncEventSendsCount(eventType string, appID string, appInstance string, appEnv string, project string, ack bool)
    IncSuccessfulRequestsCount(httpCode int)
    IncFailedRequestsCount(httpCode int)
    ObserveSuccessfulRequestDuration(dur time.Duration)
    IncKafkaRequestsCount()
    IncKafkaErrorsCount()
    ObserveKafkaRequestDuration(dur time.Duration)
}
