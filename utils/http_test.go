package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostFormWithRetryEventuallySucceeds(t *testing.T) {
	t.Setenv("HTTP_CLIENT_RETRIES", "3")
	t.Setenv("HTTP_CLIENT_TIMEOUT_MS", "500")

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&attempts, 1)
		if current < 3 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := PostFormWithRetry(ctx, server.URL, url.Values{"sip_msg": {"hello"}})
	require.NoError(t, err)
	require.NotNil(t, resp)
	_ = resp.Body.Close()
	assert.Equal(t, int32(3), atomic.LoadInt32(&attempts))
}

func TestPostFormWithRetryTimeout(t *testing.T) {
	t.Setenv("HTTP_CLIENT_RETRIES", "1")
	t.Setenv("HTTP_CLIENT_TIMEOUT_MS", "20")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := PostFormWithRetry(ctx, server.URL, url.Values{"sip_msg": {"hello"}})
	assert.Error(t, err)
}
