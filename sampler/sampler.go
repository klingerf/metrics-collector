package sampler

type Metrics map[string]float64

type Sample struct {
	Metrics  Metrics
	Earliest int64 // Unix epoch seconds
}

type Sampler interface {
	Sample() (Sample, error)
	Trim(*Sample)
}
