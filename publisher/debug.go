package publisher

import (
	"fmt"
	"os"

	"github.com/klingerf/metrics-collector/sampler"
)

type DebugPublisher struct{}

func (publisher DebugPublisher) Publish(sample *sampler.Sample) error {
	for k, v := range sample.Metrics {
		fmt.Fprintf(os.Stdout, "%s: %f\n", k, v)
	}
	return nil
}
