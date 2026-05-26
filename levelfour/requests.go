package levelfour

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/LevelFourAI/levelfour-go/core"
)

// Execute sends an HTTP request to the LevelFour API and decodes the JSON
// response into the provided response value. Use this for undocumented or
// beta endpoints not yet covered by the generated client methods.
//
// The path is appended to the client's base URL. It should start with "/".
// If body is nil, no request body is sent. If response is nil, the response
// body is discarded.
//
// Errors returned by Execute are *core.APIError and work with StatusCode,
// IsNotFound, IsRateLimited, and all other error helpers.
//
//	var result map[string]interface{}
//	err := client.Execute(ctx, http.MethodGet, "/api/v1/some-endpoint", nil, &result)
func (c *Client) Execute(ctx context.Context, method, path string, body, response interface{}) error {
	resp, err := c.do(ctx, method, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return core.NewAPIError(resp.StatusCode, resp.Header, errors.New(string(respBody)))
	}

	if response != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("levelfour: failed to decode response: %w", err)
		}
	}

	return nil
}

// RawRequest sends an HTTP request to the LevelFour API and returns the
// raw *http.Response. The caller is responsible for closing the response body.
//
// Use this when you need access to response headers (rate limit info,
// request IDs, etc.) or want full control over response handling.
//
//	resp, err := client.RawRequest(ctx, http.MethodGet, "/api/v1/health")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer resp.Body.Close()
//	fmt.Println("Status:", resp.StatusCode)
//	fmt.Println("Request-ID:", resp.Header.Get("X-Request-Id"))
func (c *Client) RawRequest(ctx context.Context, method, path string, body ...interface{}) (*http.Response, error) {
	var reqBody interface{}
	if len(body) > 0 {
		reqBody = body[0]
	}
	return c.do(ctx, method, path, reqBody)
}

// Get is a convenience method for Execute with http.MethodGet.
func (c *Client) Get(ctx context.Context, path string, response interface{}) error {
	return c.Execute(ctx, http.MethodGet, path, nil, response)
}

// Post is a convenience method for Execute with http.MethodPost.
func (c *Client) Post(ctx context.Context, path string, body, response interface{}) error {
	return c.Execute(ctx, http.MethodPost, path, body, response)
}

// Put is a convenience method for Execute with http.MethodPut.
func (c *Client) Put(ctx context.Context, path string, body, response interface{}) error {
	return c.Execute(ctx, http.MethodPut, path, body, response)
}

// Patch is a convenience method for Execute with http.MethodPatch.
func (c *Client) Patch(ctx context.Context, path string, body, response interface{}) error {
	return c.Execute(ctx, http.MethodPatch, path, body, response)
}

// Delete is a convenience method for Execute with http.MethodDelete.
func (c *Client) Delete(ctx context.Context, path string, response interface{}) error {
	return c.Execute(ctx, http.MethodDelete, path, nil, response)
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	u := strings.TrimRight(c.baseURL, "/") + path

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("levelfour: failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return nil, fmt.Errorf("levelfour: failed to create request: %w", err)
	}

	for key, values := range c.header {
		req.Header[key] = values
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.retrier.Run(c.httpClient.Do, req, nil)
	return resp, err
}
