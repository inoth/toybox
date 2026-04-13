package source

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxConfigSize = 10 << 20 // 10MB

// RemoteOption configures a RemoteSource.
type RemoteOption func(*RemoteSource)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) RemoteOption {
	return func(r *RemoteSource) {
		r.client = client
	}
}

// WithHeaders adds custom HTTP headers to config requests.
func WithHeaders(headers map[string]string) RemoteOption {
	return func(r *RemoteSource) {
		for k, v := range headers {
			r.headers[k] = v
		}
	}
}

// RemoteSource loads config from a remote HTTP endpoint.
type RemoteSource struct {
	endpoint string
	client   *http.Client
	headers  map[string]string
}

// NewRemote creates a remote HTTP config source.
func NewRemote(endpoint string, opts ...RemoteOption) *RemoteSource {
	r := &RemoteSource{
		endpoint: endpoint,
		client:   &http.Client{Timeout: 10 * time.Second},
		headers:  make(map[string]string),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Load fetches config from the remote endpoint.
func (r *RemoteSource) Load(_ string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, r.endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch config: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxConfigSize))
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}
	return string(body), nil
}

// RemoteWatcher watches a remote config endpoint for changes by polling.
type RemoteWatcher struct {
	source   *RemoteSource
	interval time.Duration
	lastHash string
}

// NewRemoteWatcher creates a remote config watcher with an optional polling interval (default 30s).
func NewRemoteWatcher(source *RemoteSource, interval ...time.Duration) *RemoteWatcher {
	d := 30 * time.Second
	if len(interval) > 0 {
		d = interval[0]
	}
	return &RemoteWatcher{source: source, interval: d}
}

// Watche polls the remote endpoint and sends a trigger when config content changes.
func (w *RemoteWatcher) Watche(ctx context.Context, trigger chan<- struct{}) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			data, err := w.source.Load("*")
			if err != nil {
				continue
			}
			h := sha256.Sum256([]byte(data))
			hash := fmt.Sprintf("%x", h)
			if w.lastHash != "" && hash != w.lastHash {
				select {
				case trigger <- struct{}{}:
				default:
				}
			}
			w.lastHash = hash
		}
	}
}
