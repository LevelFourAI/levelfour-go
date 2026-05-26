package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Verifier validates incoming webhook signatures using HMAC-SHA256.
type Verifier struct {
	secret []byte
}

// NewVerifier creates a webhook verifier from a signing secret.
// The secret may optionally be prefixed with "whsec_" (Svix format).
func NewVerifier(secret string) (*Verifier, error) {
	secret = strings.TrimPrefix(secret, "whsec_")
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook secret: %w", err)
	}
	return &Verifier{secret: decoded}, nil
}

// Verify validates the webhook signature and returns the parsed JSON payload.
// Uses a default tolerance of 5 minutes for timestamp replay protection.
func (v *Verifier) Verify(headers http.Header, body []byte) (map[string]interface{}, error) {
	return v.VerifyWithTolerance(headers, body, 5*time.Minute)
}

// VerifyWithTolerance validates the webhook signature with a custom timestamp
// tolerance for replay protection.
func (v *Verifier) VerifyWithTolerance(headers http.Header, body []byte, tolerance time.Duration) (map[string]interface{}, error) {
	msgID := firstNonEmpty(headers.Get("webhook-id"), headers.Get("svix-id"))
	msgTS := firstNonEmpty(headers.Get("webhook-timestamp"), headers.Get("svix-timestamp"))
	msgSig := firstNonEmpty(headers.Get("webhook-signature"), headers.Get("svix-signature"))

	if msgID == "" || msgTS == "" || msgSig == "" {
		return nil, fmt.Errorf("missing required webhook headers")
	}

	ts, err := strconv.ParseInt(msgTS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	diff := math.Abs(float64(time.Now().Unix() - ts))
	if diff > tolerance.Seconds() {
		return nil, fmt.Errorf("timestamp too old or too new")
	}

	toSign := []byte(fmt.Sprintf("%s.%s.", msgID, msgTS))
	toSign = append(toSign, body...)

	mac := hmac.New(sha256.New, v.secret)
	mac.Write(toSign)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	valid := false
	for _, sig := range strings.Split(msgSig, " ") {
		parts := strings.SplitN(sig, ",", 2)
		sigValue := sig
		if len(parts) == 2 {
			sigValue = parts[1]
		}
		valid = valid || hmac.Equal([]byte(expected), []byte(sigValue))
	}

	if !valid {
		return nil, fmt.Errorf("invalid signature")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON payload: %w", err)
	}
	return result, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
