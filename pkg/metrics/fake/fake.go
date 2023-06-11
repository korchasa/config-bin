package fake

import "time"

type Fake struct{}

func (f Fake) IncRequestsCount()                                {}
func (f Fake) IncEventSendsCount(_, _, _, _, _ string, _ bool)  {}
func (f Fake) IncSuccessfulRequestsCount(_ int)                 {}
func (f Fake) IncFailedRequestsCount(_ int)                     {}
func (f Fake) ObserveSuccessfulRequestDuration(_ time.Duration) {}
func (f Fake) IncKafkaRequestsCount()                           {}
func (f Fake) IncKafkaErrorsCount()                             {}
func (f Fake) ObserveKafkaRequestDuration(_ time.Duration)      {}
