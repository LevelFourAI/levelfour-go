package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "whsec_" + "MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"

func sign(secret []byte, msgID string, ts int64, body []byte) string {
	toSign := []byte(fmt.Sprintf("%s.%d.", msgID, ts))
	toSign = append(toSign, body...)
	mac := hmac.New(sha256.New, secret)
	mac.Write(toSign)
	return "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func makeHeaders(msgID string, ts int64, sig string) http.Header {
	h := http.Header{}
	h.Set("webhook-id", msgID)
	h.Set("webhook-timestamp", strconv.FormatInt(ts, 10))
	h.Set("webhook-signature", sig)
	return h
}

func TestVerify_ValidSignature(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`{"event":"test"}`)
	ts := time.Now().Unix()
	sig := sign(v.secret, "msg_123", ts, body)
	headers := makeHeaders("msg_123", ts, sig)

	result, err := v.Verify(headers, body)
	require.NoError(t, err)
	assert.Equal(t, "test", result["event"])
}

func TestVerify_InvalidSignature(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`{"event":"test"}`)
	ts := time.Now().Unix()
	headers := makeHeaders("msg_123", ts, "v1,invalidsignature")

	_, err = v.Verify(headers, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signature")
}

func TestVerify_TamperedBody(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`{"event":"test"}`)
	ts := time.Now().Unix()
	sig := sign(v.secret, "msg_123", ts, body)
	headers := makeHeaders("msg_123", ts, sig)

	tampered := []byte(`{"event":"hacked"}`)
	_, err = v.Verify(headers, tampered)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signature")
}

func TestVerify_ExpiredTimestamp(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`{"event":"test"}`)
	ts := time.Now().Add(-10 * time.Minute).Unix()
	sig := sign(v.secret, "msg_123", ts, body)
	headers := makeHeaders("msg_123", ts, sig)

	_, err = v.Verify(headers, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timestamp too old")
}

func TestVerify_FutureTimestamp(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`{"event":"test"}`)
	ts := time.Now().Add(10 * time.Minute).Unix()
	sig := sign(v.secret, "msg_123", ts, body)
	headers := makeHeaders("msg_123", ts, sig)

	_, err = v.Verify(headers, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timestamp too old")
}

func TestVerify_CustomTolerance(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`{"event":"test"}`)
	ts := time.Now().Add(-3 * time.Minute).Unix()
	sig := sign(v.secret, "msg_123", ts, body)
	headers := makeHeaders("msg_123", ts, sig)

	_, err = v.Verify(headers, body)
	require.NoError(t, err)

	_, err = v.VerifyWithTolerance(headers, body, 1*time.Minute)
	assert.Error(t, err)
}

func TestVerify_MissingHeaders(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)
	body := []byte(`{"event":"test"}`)

	tests := []struct {
		name    string
		headers http.Header
	}{
		{"missing id", func() http.Header {
			h := http.Header{}
			h.Set("webhook-timestamp", "123")
			h.Set("webhook-signature", "v1,abc")
			return h
		}()},
		{"missing timestamp", func() http.Header {
			h := http.Header{}
			h.Set("webhook-id", "msg_123")
			h.Set("webhook-signature", "v1,abc")
			return h
		}()},
		{"missing signature", func() http.Header {
			h := http.Header{}
			h.Set("webhook-id", "msg_123")
			h.Set("webhook-timestamp", "123")
			return h
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := v.Verify(tt.headers, body)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "missing required webhook headers")
		})
	}
}

func TestVerify_SvixHeaders(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`{"event":"test"}`)
	ts := time.Now().Unix()
	sig := sign(v.secret, "msg_svix", ts, body)

	h := http.Header{}
	h.Set("svix-id", "msg_svix")
	h.Set("svix-timestamp", strconv.FormatInt(ts, 10))
	h.Set("svix-signature", sig)

	result, err := v.Verify(h, body)
	require.NoError(t, err)
	assert.Equal(t, "test", result["event"])
}

func TestVerify_MultipleSignatures(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`{"event":"test"}`)
	ts := time.Now().Unix()
	validSig := sign(v.secret, "msg_123", ts, body)
	multiSig := "v1,invalidsig " + validSig + " v1,anotherinvalid"
	headers := makeHeaders("msg_123", ts, multiSig)

	result, err := v.Verify(headers, body)
	require.NoError(t, err)
	assert.Equal(t, "test", result["event"])
}

func TestVerify_InvalidTimestamp(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`{"event":"test"}`)
	h := http.Header{}
	h.Set("webhook-id", "msg_123")
	h.Set("webhook-timestamp", "not-a-number")
	h.Set("webhook-signature", "v1,abc")

	_, err = v.Verify(h, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid timestamp")
}

func TestVerify_InvalidJSON(t *testing.T) {
	v, err := NewVerifier(testSecret)
	require.NoError(t, err)

	body := []byte(`not json`)
	ts := time.Now().Unix()
	sig := sign(v.secret, "msg_123", ts, body)
	headers := makeHeaders("msg_123", ts, sig)

	_, err = v.Verify(headers, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON payload")
}

func TestNewVerifier_InvalidSecret(t *testing.T) {
	_, err := NewVerifier("not-base64-!!!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid webhook secret")
}

func TestNewVerifier_WithoutPrefix(t *testing.T) {
	rawSecret := "MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"
	v, err := NewVerifier(rawSecret)
	require.NoError(t, err)

	body := []byte(`{"ok":true}`)
	ts := time.Now().Unix()
	sig := sign(v.secret, "msg_1", ts, body)
	headers := makeHeaders("msg_1", ts, sig)

	result, err := v.Verify(headers, body)
	require.NoError(t, err)
	assert.Equal(t, true, result["ok"])
}
