package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithLevelFour_InjectsClient(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := ClientFromContext(r.Context())
		if client == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	wrapped := WithLevelFour(handler, Config{
		APIKey: "l4_test_middleware_key",
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestWithLevelFour_WithCustomBaseURL(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := ClientFromContext(r.Context())
		require.NotNil(t, client)
		w.WriteHeader(http.StatusOK)
	})

	wrapped := WithLevelFour(handler, Config{
		APIKey:  "l4_test_custom_url",
		BaseURL: "https://custom.api.example.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestWithLevelFour_PanicsOnInvalidKey(t *testing.T) {
	assert.Panics(t, func() {
		WithLevelFour(http.NotFoundHandler(), Config{
			APIKey: "invalid_key",
		})
	})
}

func TestWithLevelFour_PanicsOnEmptyKey(t *testing.T) {
	t.Setenv("LEVELFOUR_API_KEY", "")
	assert.Panics(t, func() {
		WithLevelFour(http.NotFoundHandler(), Config{})
	})
}

func TestClientFromContext_NilWithoutMiddleware(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	client := ClientFromContext(req.Context())
	assert.Nil(t, client)
}

func TestClientFromContext_HasEscapeHatches(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := ClientFromContext(r.Context())
		require.NotNil(t, client)
		assert.NotNil(t, client.Recommendations)
		w.WriteHeader(http.StatusOK)
	})

	wrapped := WithLevelFour(handler, Config{
		APIKey: "l4_test_escape",
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
