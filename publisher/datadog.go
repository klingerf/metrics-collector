package publisher

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/klingerf/metrics-collector/sampler"
)

type DatadogPublisher struct {
	client *statsd.Client
	tags   []string
}

func NewDatadog(addr, namespace, source, service string) *DatadogPublisher {
	client, err := statsd.New(addr)
	if err != nil {
		panic(err)
	}
	tags := []string{
		"namespace:" + namespace,
		"source:" + source,
		"service:" + service,
	}
	return &DatadogPublisher{client, tags}
}

func (publisher DatadogPublisher) Publish(sample *sampler.Sample) (err error) {
	for k, v := range sample.Metrics {
		err = publisher.client.Gauge(k, v, publisher.tags, 1)
		if err != nil {
			return err
		}
	}
	return nil
}
