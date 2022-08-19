package middleware

import (
	"time"

	"github.com/micromind/services"

	"github.com/go-kit/kit/metrics"
)

func InstrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	quoteLength metrics.Gauge,
) services.ServiceMiddleware {
	return func(next services.ZenService) services.ZenService {
		return instrmw{requestCount, requestLatency, quoteLength, next}
	}
}

type instrmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	quoteLength    metrics.Gauge
	next           services.ZenService
}

func (mw instrmw) Quote() (quote string, author string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "quote"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.quoteLength.Set(float64(len(quote)))
	}(time.Now())

	quote, author, err = mw.next.Quote()
	return
}

func (mw instrmw) Question() (question string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "question"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	question, err = mw.next.Question()
	return
}
