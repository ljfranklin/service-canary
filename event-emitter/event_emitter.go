package event_emitter

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ljfranklin/service-canary/config"
	"github.com/pivotal-golang/lager"
)

type Emitter interface {
	Emit(serviceName string, value int, tags map[string]string) error
}

type datadogEmitter struct {
	apiKey string
	logger lager.Logger
}

func NewDataDogEmitter(config *config.Config) Emitter {
	return &datadogEmitter{
		apiKey: config.DatadogApiKey,
		logger: config.Logger,
	}
}

func (e *datadogEmitter) Emit(serviceName string, value int, tags map[string]string) error {

	url := fmt.Sprintf(
		"https://app.datadoghq.com/api/v1/series?api_key=%s",
		e.apiKey,
	)

	client := http.DefaultClient

	formattedTags := []string{}
	for k, v := range tags {
		formattedTags = append(formattedTags, fmt.Sprintf(`"%s:%s"`, k, v))
	}

	jsonString := []byte(fmt.Sprintf(`{"series":`+
		`[{`+
		`"metric":"services-canary.%s.instance",`+
		`"points":[[%d, %d]],`+
		`"tags":[%s]`+
		`}]`+
		`}`,
		serviceName,
		time.Now().Unix(),
		value,
		strings.Join(formattedTags, ","),
	))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Received a non-200 response from datadog: %d", resp.StatusCode)
	}

	return nil
}
