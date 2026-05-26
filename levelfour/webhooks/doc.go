// Package webhooks provides HMAC-SHA256 signature verification for
// incoming LevelFour webhook payloads.
//
// Create a [Verifier] with your webhook signing secret, then call
// [Verifier.Verify] on each incoming request:
//
//	verifier, err := webhooks.NewVerifier("whsec_...")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	payload, err := verifier.Verify(r.Header, body)
//	if err != nil {
//	    http.Error(w, "invalid signature", http.StatusUnauthorized)
//	    return
//	}
//
// Signatures are verified using HMAC-SHA256 with replay protection
// via timestamp validation (default 5-minute tolerance). Use
// [Verifier.VerifyWithTolerance] to customize the tolerance window.
//
// This package handles incoming payload verification. To manage webhook
// endpoint subscriptions (register, list, delete), use the webhooks
// sub-client available as client.Webhooks on the [levelfour.Client].
package webhooks
