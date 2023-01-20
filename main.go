package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/micromind/middleware"
	"github.com/micromind/repositories"
	"github.com/micromind/services"
	"github.com/micromind/transports"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy question requests")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "zen_service",
		Name:      "request_count",
		Help:      "Number of requests recieved.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "zen_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	quoteLength := kitprometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: "my_group",
		Subsystem: "zen_service",
		Name:      "quote_length",
		Help:      "The length of each quote result",
	}, []string{})

	// Mongodb
	zr := repositories.NewInstanceOfZenRepository()

	var svc services.ZenService
	svc = services.NewInstanceOfZenService(zr)
	svc = middleware.ProxyingMiddlewere(context.Background(), *proxy, logger)(svc)
	svc = middleware.LoggingMiddleware(logger)(svc)
	svc = middleware.InstrumentingMiddleware(requestCount, requestLatency, quoteLength)(svc)

	quoteHandler := httptransport.NewServer(
		transports.MakeQuoteEndpoint(svc),
		transports.DecodeQuoteRequest,
		transports.EncodeResponse,
	)

	questionHandler := httptransport.NewServer(
		transports.MakeQuestionEndpoint(svc),
		transports.DecodeQuestionRequest,
		transports.EncodeResponse,
	)

	rootHandler := httptransport.NewServer(
		transports.MakeRootEndpoint(),
		httptransport.NopRequestDecoder,
		transports.EncodeResponse,
	)

	http.Handle("/", rootHandler)
	http.Handle("/quote", quoteHandler)
	http.Handle("/question", questionHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
