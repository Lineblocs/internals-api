package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultHTTPTimeout = 5 * time.Second
	defaultHTTPRetries = 3
	defaultBackoffBase = 200 * time.Millisecond
)

func PostFormWithRetry(ctx context.Context, endpoint string, form url.Values) (*http.Response, error) {
	client := &http.Client{Timeout: readHTTPTimeout()}
	retries := readRetryCount()

	encoded := form.Encode()
	var lastErr error

	for attempt := 1; attempt <= retries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(encoded))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		if resp != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			lastErr = fmt.Errorf("downstream returned non-success status: %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		if attempt < retries {
			backoff := time.Duration(attempt) * defaultBackoffBase
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("request failed")
	}

	logrus.WithFields(logrus.Fields{
		"endpoint": endpoint,
		"retries":  retries,
	}).WithError(lastErr).Error("downstream HTTP call failed after retries")

	return nil, lastErr
}

func readHTTPTimeout() time.Duration {
	value := os.Getenv("HTTP_CLIENT_TIMEOUT_MS")
	if value == "" {
		return defaultHTTPTimeout
	}

	ms, err := strconv.Atoi(value)
	if err != nil || ms <= 0 {
		return defaultHTTPTimeout
	}

	return time.Duration(ms) * time.Millisecond
}

func readRetryCount() int {
	value := os.Getenv("HTTP_CLIENT_RETRIES")
	if value == "" {
		return defaultHTTPRetries
	}

	tries, err := strconv.Atoi(value)
	if err != nil || tries <= 0 {
		return defaultHTTPRetries
	}

	return tries
}
