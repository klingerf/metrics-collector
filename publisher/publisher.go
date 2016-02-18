package publisher

import (
	"github.com/klingerf/metrics-collector/sampler"
)

type Publisher interface {
	Publish(*sampler.Sample) error
}
