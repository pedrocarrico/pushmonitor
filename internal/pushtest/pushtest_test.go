package pushtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/pedrocarrico/pushmonitor/internal/logger"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Silence logger for tests
	logger.Init("error", io.Discard)
}

func TestPushTestShouldRun(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{
			name:    "No command should return true",
			command: "",
			want:    true,
		},
		{
			name:    "Valid command should return true",
			command: "echo 'test'",
			want:    true,
		},
		{
			name:    "Invalid command should return false",
			command: "nonexistentcommand",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PushTest{
				Name:    "test",
				Command: tt.command,
			}
			assert.Equal(t, tt.want, p.shouldRun(), "shouldRun() should return expected value")
		})
	}
}

func TestPushTestExecuteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errorServer.Close()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Successful request",
			url:     server.URL,
			wantErr: false,
		},
		{
			name:    "Failed request",
			url:     errorServer.URL,
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			url:     "http://invalid-url",
			wantErr: true,
		},
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PushTest{
				Name: "test",
				URL:  tt.url,
			}
			err := p.executeRequest(client)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPushTest(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := &PushTest{
		Name:     "test",
		URL:      server.URL,
		Interval: 1,
		Retries:  1,
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go p.Run(client, &wg, ctx)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		assert.GreaterOrEqual(t, requestCount, 1, "Expect at least one request to have been made")
	case <-time.After(3 * time.Second):
		assert.Fail(t, "Test timed out")
	}
}
