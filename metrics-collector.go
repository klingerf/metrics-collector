package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/klingerf/metrics-collector/publisher"
	"github.com/klingerf/metrics-collector/sampler"
)

func main() {
	defaultSource, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	metricsURL := flag.String("metrics-url",
		"http://127.0.0.1:9990/admin/metrics.json",
		"Address of TwitterServer metrics to collect")
	publisherStr := flag.String("publisher", "",
		"Stats publisher; supported publishers: datadog, debug")
	trimmed := flag.Bool("trimmed", false,
		"Reduce metrics before sending them to the publisher")
	namespace := flag.String("namespace", "",
		"Namespace where metrics are being collected")
	source := flag.String("source", defaultSource,
		"Common name for entity that is being collected")
	service := flag.String("service", "",
		"Logical name for service that is being collected")
	period := flag.Duration("period", 60*time.Second,
		"Polling period")
	datadogStatsd := flag.String("datadog-statsd", "127.0.0.1:8125",
		"Address of StatsD process used by the datadog publisher")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	var p publisher.Publisher
	switch *publisherStr {
	case "":
		fmt.Fprintf(os.Stderr, "Publisher is required\n")
		flag.Usage()
		os.Exit(1)
	case "datadog":
		p = publisher.NewDatadog(*datadogStatsd, *namespace, *source, *service)
	case "debug":
		p = publisher.DebugPublisher{}
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized publisher: %s\n", *publisherStr)
		os.Exit(1)
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, os.Kill)

	t := time.NewTicker(*period)
	s := sampler.NewTwitterServerSampler(*metricsURL)

	for {
		select {
		case <-t.C:
			sample, err := s.Sample()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			if *trimmed {
				s.Trim(sample)
			}
			err = p.Publish(sample)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		case <-exitChan:
			os.Exit(0)
		}
	}
}
