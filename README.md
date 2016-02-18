# metrics-collector

metrics-collector is a simple Go script that can be run to poll the metrics
endpoint of a [TwitterServer](https://twitter.github.io/twitter-server/Features.html#metrics)
process at an established interval and publish those metrics to a configurable
stats backend. The script is optimized for collecting metrics from a
[linkerd](https://github.com/BuoyantIO/linkerd) process.

Try publishing a limited set of metrics from a linkerd admin server running
locally on port 9990 to stdout (useful for debugging):

```bash
$ go run metrics-collector.go -metrics-url http://localhost:9990/admin/metrics.json -publisher debug -period 5s -trimmed true
jvm/uptime: 8600.000000
rt/int/srv/0.0.0.0/4140/request_latency_ms.p50: 0.000000
rt/int/srv/0.0.0.0/4140/request_latency_ms.p99: 0.000000
...
```

Or send the full set of metrics to a datadog-agent process, like so:

```
$ go run metrics-collector.go -publisher datadog -datadog-statsd=localhost:8125 -namespace=dev -service=mysvc
```

Check out the help text for all available command line flags:

```bash
$ go run metrics-collector.go -help
Usage: metrics-collector [flags]
  -datadog-statsd string
      Address of StatsD process used by the datadog publisher (default "127.0.0.1:8125")
  -metrics-url string
      Address of TwitterServer metrics to collect (default "http://127.0.0.1:9990/admin/metrics.json")
  -namespace string
      Namespace where metrics are being collected
  -period duration
      Polling period (default 1m0s)
  -publisher string
      Stats publisher; supported publisher: datadog, debug
  -service string
      Logical name for service that is being collected
  -source string
      Common name for entity that is being collected (default os.Hostname())
  -trimmed
      Reduce metrics before sending them to the publisher
exit status 2
```

Am gladly accepting pull requests for additional publishers. :rainbow:
