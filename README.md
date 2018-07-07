# Insighted Go - Microservice Toolkit

This toolkit provides an experimental packages to put together server and pubsub daemons with the following features:

* Standardized configuration and logging
* Health check endpoints with configurable strategies
* Configuration for managing pprof endpoints and log levels
* Basic interfaces to define expectations and vocabulary
* Structured logging containing basic request information
* Useful metrics for endpoints
* Graceful shutdowns

## Goals

- RPC as the primary messaging pattern
- Event sourcing library
- Supporting messaging patterns other than RPC — e.g. Pub/Sub, CQRS, etc.

## Non-goals

- Re-implementing functionality that can be provided by adapting existing software
- Having opinions on operational concerns: deployment, configuration, process supervision, orchestration, etc.


## The Kit Package

This is an experimental reference for creating microservices in Go. 

The rationale behind this package:

* A more opinionated server with fewer choices.
* Gin-Gonic is used for serving HTTP/JSON & gRPC is used for serving HTTP2/RPC
* Monitoring and metrics are handled by a sidecar (ie. Cloud Endpoints)
* Logs always go to stdout/stderr
* Using Go's 1.8 graceful HTTP shutdown
* Services using this package are meant for deploy to Kuberntes.

If you experience any issues please create an [issue](https://github.com/insighted4/insighted-go/issues). 

## Examples

Several reference implementations utilizing `server` and `pubsub` are available in the [`examples`](examples/) subdirectory.

## Contributing

Please see [CONTRIBUTING.md](/CONTRIBUTING.md).

## Related projects

This project is highly influenced by [gizmo](https://github.com/nytimes/gizmo) and [go-kit](https://github.com/go-kit/kit) design.

### Service frameworks

- [go-kit](https://github.com/go-kit/kit), a programming toolkit for building microservices  ★
- [gizmo](https://github.com/nytimes/gizmo), a microservice toolkit from The New York Times ★
- [go-micro](https://github.com/myodc/go-micro), a microservices client/server library ★

### Individual components

- [grpc/grpc-go](https://github.com/grpc/grpc-go), HTTP/2 based RPC
- [sirupsen/logrus](https://github.com/sirupsen/logrus), structured, pluggable logging for Go ★

### Web frameworks

- [Gorilla](http://www.gorillatoolkit.org)
- [Gin](https://gin-gonic.github.io/gin/)

## Additional reading

- [Architecting for the Cloud](https://slideshare.net/stonse/architecting-for-the-cloud-using-netflixoss-codemash-workshop-29852233) — Netflix
- [Dapper, a Large-Scale Distributed Systems Tracing Infrastructure](http://research.google.com/pubs/pub36356.html) — Google
- [Your Server as a Function](http://monkey.org/~marius/funsrv.pdf) (PDF) — Twitter

---

Development supported by [Insighted4](https://github.com/insighted4).