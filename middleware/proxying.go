package middleware

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/micromind/services"
	"github.com/micromind/transports"

	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
)

func ProxyingMiddlewere(ctx context.Context, instances string, logger log.Logger) services.ServiceMiddleware {
	// If instances is empty, dont proxy
	if instances == "" {
		logger.Log("proxy_to", "none")
		return func(next services.ZenService) services.ZenService { return next }
	}

	// Set some parameters for our client
	var (
		qps         = 100                    // beyond which we will return an error
		maxAttempts = 3                      // per request, before giving up
		maxTime     = 250 * time.Millisecond // wallclock time, before giving up
	)

	// Construct an endpoint for each instance in the list, and add it to a fixed
	// set of endpoints. In a real service, rather than doing this by hand
	// you would probably use package sd's support for your service discovery system
	var (
		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)
	logger.Log("proxy_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeQuestionProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	// Now, build a single, retrying, load-balancing endpoint out of all of those
	// individual components.
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	// Finally, return the ServiceMiddleware, implemented by proxymw
	return func(next services.ZenService) services.ZenService {
		return proxymw{ctx, next, retry}
	}

}

func makeQuestionProxy(ctx context.Context, instance string) endpoint.Endpoint {
	// Parse instance url
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/question"
	}
	return httptransport.NewClient(
		"GET",
		u,
		transports.EndcodeRequest,
		transports.DecodeQuestionResponse,
	).Endpoint()
}

// Proxy middleware implements ZenSerivce, serving Quote (and all other) requests on to another ZenService
// but forwading Question requests to the provided endpoint.
type proxymw struct {
	ctx      context.Context
	next     services.ZenService
	question endpoint.Endpoint
}

// Client endpoint (used to invoke, rather than serve, a request)
func (mw proxymw) Question() (string, error) {
	response, err := mw.question(mw.ctx, transports.QuestionRequest{})
	if err != nil {
		return "", err
	}
	resp := response.(transports.QuestionResponse)
	return resp.Q, nil
}

func (mw proxymw) Quote() (string, string, error) {
	return mw.next.Quote()
}

func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}
