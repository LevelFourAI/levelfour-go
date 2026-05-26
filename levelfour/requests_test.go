package levelfour

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/LevelFourAI/levelfour-go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	client, err := NewClient("l4_test_key",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
		WithNoRetries(),
	)
	require.NoError(t, err)
	return client, server
}

func TestExecute_Get(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/test", r.URL.Path)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer l4_test_key")
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	defer server.Close()

	var result map[string]string
	err := client.Execute(context.Background(), http.MethodGet, "/api/v1/test", nil, &result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestExecute_Post(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, _ := io.ReadAll(r.Body)
		var payload map[string]string
		require.NoError(t, json.Unmarshal(body, &payload))
		assert.Equal(t, "bar", payload["foo"])

		w.WriteHeader(201)
		_, _ = w.Write([]byte(`{"id":"123"}`))
	})
	defer server.Close()

	var result map[string]string
	err := client.Execute(context.Background(), http.MethodPost, "/api/v1/create",
		map[string]string{"foo": "bar"}, &result)
	require.NoError(t, err)
	assert.Equal(t, "123", result["id"])
}

func TestExecute_ReturnsAPIError(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"error":{"code":"NOT_FOUND","message":"not found"}}`))
	})
	defer server.Close()

	err := client.Execute(context.Background(), http.MethodGet, "/api/v1/missing", nil, nil)
	require.Error(t, err)

	assert.True(t, IsNotFound(err))
	assert.Equal(t, 404, StatusCode(err))
	assert.Equal(t, "not found", ErrorMessage(err))
	assert.Equal(t, "NOT_FOUND", ErrorCode(err))

	var apiErr *core.APIError
	assert.True(t, errors.As(err, &apiErr))
}

func TestExecute_NoContent(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	var result map[string]string
	err := client.Execute(context.Background(), http.MethodDelete, "/api/v1/resource", nil, &result)
	require.NoError(t, err)
}

func TestExecute_RetriesOn503(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) <= 2 {
			w.WriteHeader(503)
			_, _ = w.Write([]byte("unavailable"))
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client, err := NewClient("l4_test_key",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
		WithMaxRetries(2),
	)
	require.NoError(t, err)

	var result map[string]bool
	err = client.Get(context.Background(), "/api/v1/test", &result)
	require.NoError(t, err)
	assert.True(t, result["ok"])
	assert.Equal(t, int32(3), attempts.Load())
}

func TestGet(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer server.Close()

	var result map[string]bool
	err := client.Get(context.Background(), "/api/v1/health", &result)
	require.NoError(t, err)
	assert.True(t, result["ok"])
}

func TestPost(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"created":true}`))
	})
	defer server.Close()

	var result map[string]bool
	err := client.Post(context.Background(), "/api/v1/items", map[string]string{"name": "test"}, &result)
	require.NoError(t, err)
	assert.True(t, result["created"])
}

func TestPut(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	})
	defer server.Close()

	err := client.Put(context.Background(), "/api/v1/items/1", map[string]string{"name": "updated"}, nil)
	require.NoError(t, err)
}

func TestPatch(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	})
	defer server.Close()

	err := client.Patch(context.Background(), "/api/v1/items/1", map[string]string{"name": "patched"}, nil)
	require.NoError(t, err)
}

func TestDelete(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.Delete(context.Background(), "/api/v1/items/1", nil)
	require.NoError(t, err)
}

func TestRawRequest(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req-abc")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"data":"raw"}`))
	})
	defer server.Close()

	resp, err := client.RawRequest(context.Background(), http.MethodGet, "/api/v1/raw")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "req-abc", resp.Header.Get("X-Request-Id"))

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"data":"raw"}`, string(body))
}

func TestRawRequest_WithBody(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]string
		_ = json.Unmarshal(body, &payload)
		assert.Equal(t, "value", payload["key"])
		w.WriteHeader(200)
	})
	defer server.Close()

	resp, err := client.RawRequest(context.Background(), http.MethodPost, "/api/v1/raw",
		map[string]string{"key": "value"})
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestExecute_ContextCancelled(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.Execute(ctx, http.MethodGet, "/api/v1/test", nil, nil)
	require.Error(t, err)
}
