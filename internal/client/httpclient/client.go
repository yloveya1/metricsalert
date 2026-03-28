package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/yloveya1/metricsalert/internal/model"
)

var (
	updateCounterEndpoint = "/update/%s/%s/%d"
	updateGaugeEndpoint   = "/update/%s/%s/%g"
	updateEndpoint        = "/update"
)

type Config struct {
	Host string
}

type HTTPClient struct {
	cfg    Config
	client *http.Client
}

func NewClient(cfg Config) *HTTPClient {
	return &HTTPClient{
		cfg:    cfg,
		client: http.DefaultClient, // todo настроить
	}
}

func (h *HTTPClient) SendMetric(metric *models.Metrics) error {
	resp, err := h.sendRequest(metric)
	if err != nil {
		return fmt.Errorf("request error, err: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	return nil
}

func (h *HTTPClient) sendRequest(metrics *models.Metrics) (*http.Response, error) {
	resURL := h.cfg.Host + updateEndpoint

	body, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf("marshal metrics error, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return resp, nil
}

func formURL(url string, metrics *models.Metrics) (string, error) {
	switch metrics.MType {
	case models.Counter:
		return url + fmt.Sprintf(updateCounterEndpoint, metrics.MType, metrics.ID, *metrics.Delta), nil
	case models.Gauge:
		return url + fmt.Sprintf(updateGaugeEndpoint, metrics.MType, metrics.ID, *metrics.Value), nil
	default:
		return "", fmt.Errorf("unsupported metrics type: %s", metrics.MType)
	}
}
