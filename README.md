# micromind
Daily dose of zen.

## TODO
- TESTING !!
    - Unit tests
    - Component tests?
    - Performance tests? (In go?)
- Add list of quotes / questions to init-mongodb
- Proper error handling from zenService up to transports layer, ie return some form of error (or cool default quote) when DB returns an empty string for a quote
    - https://hackernoon.com/handling-errors-in-golang-grpc-and-go-kit-services-d0fa0a112449
- Run on local kubernetes cluster! Using minikube? 
    - Micromind and mongodb running in seperate pods? 
    - Replica sets for proxying request around to several instances of micromind? All replicas connect to same instance of mongodb
    - Prometheus? Grafana? Jaeger? Envoy? Fluentd?
        - https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang
- Deploy to a cloud provider
    - https://ccloud.google.com

## code References
- [Go kit stringsvc example](https://github.com/go-kit/kit)
- [kiethweaver go boilerplate](https://github.com/keithweaver/go-boilerplate/tree/v1.0.1)
- [MongoDB running on docker](https://faun.pub/initialize-mongodb-running-on-a-docker-container-889a43c5668a)
