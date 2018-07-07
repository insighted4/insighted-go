/*
Package kit provides experimental packages to put together server and pubsub daemons with the following features:

* Standardized configuration and logging
* Health check endpoints with configurable strategies
* Configuration for managing pprof endpoints and log levels
* Basic interfaces to define expectations and vocabulary
* Structured logging containing basic request information
* Useful metrics for endpoints
* Graceful shutdowns

This is an experimental reference for creating Microservices in Go.

The rationale behind this package:

* A more opinionated server with fewer choices.
* go-kit is used for serving HTTP/JSON & gRPC is used for serving HTTP2/RPC
* Monitoring and metrics are handled by a sidecar (ie. Cloud Endpoints)
* Logs always go to stdout/stderr
* Using Go's 1.8 graceful HTTP shutdown
* Services using this package are meant for deploy to GCP with GKE and Cloud Endpoints.

If you experience any issues please create an [issue](https://github.com/insighted4/insighted-go/issues).

## Examples

Several reference implementations utilizing `server` and `pubsub` are available in the [`examples`](examples/) subdirectory.

*/
package kit
