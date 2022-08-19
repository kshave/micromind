package middleware

import (
	"time"

	"github.com/go-kit/log"
	"github.com/micromind/services"
)

func LoggingMiddleware(
	logger log.Logger,
) services.ServiceMiddleware {
	return func(next services.ZenService) services.ZenService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	next   services.ZenService
}

func (mw logmw) Quote() (quote string, author string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "quote",
			"quote", quote,
			"author", author,
			"took", time.Since(begin),
		)
	}(time.Now())

	quote, author, err = mw.next.Quote()
	return
}

func (mw logmw) Question() (question string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "question",
			"question", question,
			"took", time.Since(begin),
		)
	}(time.Now())

	question, err = mw.next.Question()
	return
}
