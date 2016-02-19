# metrics-collector

metrics-collector is a simple Go script that can be run to poll the metrics
endpoint of a [TwitterServer](https://twitter.github.io/twitter-server/Features.html#metrics)
process at an established interval and publish those metrics to a configurable
stats backend. The script is optimized for collecting metrics from
[linkerd](https://github.com/BuoyantIO/linkerd) (which uses TwitterServer).

## Examples

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

```bash
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
      Stats publisher; supported publishers: datadog, debug
  -service string
      Logical name for service that is being collected
  -source string
      Common name for entity that is being collected (default os.Hostname())
  -trimmed
      Reduce metrics before sending them to the publisher
exit status 2
```

## Docker

Use the included Dockerfile to build a metrics-collector Docker image, like so:

```bash
$ docker build -t metrics-collector .
...
$ docker run --rm metrics-collector -publisher debug
...
```

You can also pull a pre-built Docker image from:

https://hub.docker.com/r/klingerf/metrics-collector/

## Kubernetes

The metrics-collector Docker image from above makes for a great sidecar process
when running linkerd in [Kubernetes](https://github.com/kubernetes/kubernetes).

Here's a sample replication controller that shows how to run them together:

```yaml
---
kind: ReplicationController
apiVersion: v1
metadata:
  namespace: prod
  name: linkerd
spec:
  replicas: 3
  selector:
    app: linkerd
  template:
    metadata:
      labels:
        app: linkerd
    spec:
      dnsPolicy: ClusterFirst
      volumes:
      - name: linkerd-config
        secret:
          secretName: "linkerd-web-config"
      containers:
      - name: linkerd
        image: buoyantio/linkerd:latest
        imagePullPolicy: Always
        args:
        - "/io.buoyant/linkerd/config/config.yaml"
        ports:
        - name: router
          containerPort: 4140
        - name: admin
          containerPort: 9990
        volumeMounts:
        - name: "linkerd-config"
          mountPath: "/io.buoyant/linkerd/config"
          readOnly: true
      - name: collector
        image: klingerf/metrics-collector:latest
        imagePullPolicy: Always
        command:
        - "app"
        args:
        - "-publisher=datadog"
        - "-datadog-statsd=dd-agent.datadog.svc.cluster.local:8125"
        - "-trimmed=true"
        - "-namespace=prod"
        - "-service=linkerd"
```

This setup assumes that you have the dd-agent process running in a separate pod
with a NodePort service running in front of it.

## Contributing

Am gladly accepting pull requests for additional publishers, or any other
features that might be lacking. :rainbow:
