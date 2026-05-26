package levelfour

import (
	"strings"
	"testing"
)

func TestNewClient_ExplicitKey(t *testing.T) {
	c, err := NewClient("l4_live_test123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.Client == nil {
		t.Fatal("expected embedded base client to be non-nil")
	}
	if c.retrier == nil {
		t.Fatal("expected retrier to be non-nil")
	}
}

func TestNewClient_TestPrefix(t *testing.T) {
	c, err := NewClient("l4_test_abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_EnvFallback(t *testing.T) {
	t.Setenv("LEVELFOUR_API_KEY", "l4_live_from_env")

	c, err := NewClient("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_MissingKey(t *testing.T) {
	t.Setenv("LEVELFOUR_API_KEY", "")

	_, err := NewClient("")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "no API key provided") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestNewClient_InvalidPrefix(t *testing.T) {
	t.Setenv("LEVELFOUR_API_KEY", "")

	_, err := NewClient("sk_bad_key")
	if err == nil {
		t.Fatal("expected error for invalid prefix")
	}
	if !strings.Contains(err.Error(), "invalid API key format") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	c, err := NewClient("l4_live_test",
		WithBaseURL("https://custom.api.com"),
		WithMaxRetries(5),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_WithNoRetries(t *testing.T) {
	c, err := NewClient("l4_test_noretry", WithNoRetries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.retrier == nil {
		t.Fatal("expected retrier to be non-nil even with no retries")
	}
}
