package levelfour

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testLogger struct {
	mu       sync.Mutex
	debugLog []string
	warnLog  []string
}

func (l *testLogger) Debug(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.debugLog = append(l.debugLog, msg)
}

func (l *testLogger) Warn(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.warnLog = append(l.warnLog, msg)
}

func TestLogger_DebugOnSuccess(t *testing.T) {
	logger := &testLogger{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient("l4_test_key",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
		WithNoRetries(),

		WithLogger(logger),
	)
	require.NoError(t, err)

	_ = client.Get(context.Background(), "/api/v1/test", nil)

	logger.mu.Lock()
	defer logger.mu.Unlock()
	assert.Len(t, logger.debugLog, 1)
	assert.Contains(t, logger.debugLog[0], "request")
	assert.Empty(t, logger.warnLog)
}

func TestLogger_WarnOnError(t *testing.T) {
	logger := &testLogger{}

	client, err := NewClient("l4_test_key",
		WithBaseURL("http://localhost:1"),
		WithNoRetries(),

		WithLogger(logger),
	)
	require.NoError(t, err)

	_ = client.Get(context.Background(), "/unreachable", nil)

	logger.mu.Lock()
	defer logger.mu.Unlock()
	assert.Len(t, logger.debugLog, 1)
	assert.Len(t, logger.warnLog, 1)
	assert.Contains(t, logger.warnLog[0], "failed")
}

func TestLogger_GeneratedSubclientLogs(t *testing.T) {
	logger := &testLogger{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"success":true,"data":{}}`))
	}))
	defer server.Close()

	client, err := NewClient("l4_test_key",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
		WithNoRetries(),
		WithLogger(logger),
	)
	require.NoError(t, err)

	_, _ = client.Recommendations.GetSavingsByProvider(context.Background())

	logger.mu.Lock()
	defer logger.mu.Unlock()
	assert.NotEmpty(t, logger.debugLog, "generated sub-client calls must produce debug logs")
	assert.Contains(t, logger.debugLog[0], "levelfour: request")
}

func TestLogger_NilLoggerNoOp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient("l4_test_key",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
		WithNoRetries(),
	)
	require.NoError(t, err)

	err = client.Get(context.Background(), "/api/v1/test", nil)
	require.NoError(t, err)
}
