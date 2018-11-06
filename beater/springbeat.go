package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"errors"
	"github.com/SietseVisser/springbeat/config"
	"net/url"
)

// Springbeat configuration.
type Springbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client

	urls []*url.URL

	metricsStats bool
	healthStats  bool
}

// New creates an instance of springbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Springbeat{
		done:   make(chan struct{}),
		config: c,
	}

	//define default URL if none provided
	var urlConfig []string
	if config.URLs != nil {
		urlConfig = config.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1"}
	}

	bt.urls = make([]*url.URL, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u, err := url.Parse(urlConfig[i])
		if err != nil {
			logp.Err("Invalid Spring Boot URL: %v", err)
			return nil, err
		}
		bt.urls[i] = u
	}

	if config.Stats.Metrics != nil {
		bt.metricsStats = *config.Stats.Metrics
	} else {
		bt.metricsStats = true
	}

	if config.Stats.Health != nil {
		bt.healthStats = *config.Stats.Health
	} else {
		bt.healthStats = true
	}

	if !bt.metricsStats && !bt.metricsStats {
		return nil, errors.New("Invalid statistics configuration")
	}

	logp.Debug("springbeat", "Init springbeat")
	logp.Debug("springbeat", "Period %v\n", bt.config.Period)
	logp.Debug("springbeat", "Watch %v", bt.urls)
	logp.Debug("springbeat", "Metrics statistics %t\n", bt.metricsStats)
	logp.Debug("springbeat", "Health statistics %t\n", bt.healthStats)

	return bt, nil
}

// Run starts springbeat.
func (bt *Springbeat) Run(b *beat.Beat) error {
	logp.Info("springbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"counter": counter,
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
		counter++
	}
}

// Stop stops springbeat.
func (bt *Springbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
