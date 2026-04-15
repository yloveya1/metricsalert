package httpclient

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	models "github.com/yloveya1/metricsalert/internal/model"
)

var (
	updateEndpoint     = "/update/"
	updateListEndpoint = "/updates/"
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
		cfg: cfg,
		client: &http.Client{Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     20 * time.Second,
			},
		},
	}
}

func (h *HTTPClient) SendMetric(metric *models.Metrics) error {
	resp, err := h.sendRequest(updateEndpoint, metric)
	if err != nil {
		return fmt.Errorf("request error, err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	return nil
}

func (h *HTTPClient) SendMetricList(metrics []*models.Metrics) error {
	resp, err := h.sendRequest(updateListEndpoint, metrics)
	if err != nil {
		return fmt.Errorf("request error, err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	return nil
}

func (h *HTTPClient) sendRequest(endpoint string, data any) (*http.Response, error) {
	resURL := h.cfg.Host + endpoint

	body, err := json.Marshal(&data)
	if err != nil {
		return nil, fmt.Errorf("marshal metrics error, err: %w", err)
	}

	compressBody, err := compress(body)
	if err != nil {
		return nil, fmt.Errorf("compress metrics error, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resURL, compressBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return resp, nil
}

func compress(data []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer, err: %w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return &buf, nil
}
