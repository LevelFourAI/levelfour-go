package levelfour

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	baseclient "github.com/LevelFourAI/levelfour-go/client"
	"github.com/LevelFourAI/levelfour-go/internal"
	"github.com/LevelFourAI/levelfour-go/option"
)

// Version is the current SDK version.
const Version = "0.2.0" // x-release-please-version

const DefaultBaseURL = "https://api.levelfour.ai"

const defaultTimeout = 30 * time.Second

var validPrefixes = []string{"l4_live_", "l4_test_"}

// Client wraps the generated base client with convenience methods
// for raw HTTP access, error inspection, and future extensions.
// All generated service fields (Recommendations, Audit, Costs, Providers, etc.)
// are accessible directly through embedding.
type Client struct {
	*baseclient.Client
	httpClient *http.Client
	retrier    *internal.Retrier
	baseURL    string
	header     http.Header
}

// ClientOption configures the LevelFour client.
type ClientOption func(*clientConfig)

type clientConfig struct {
	baseURL    string
	maxRetries uint
	noRetries  bool
	httpClient *http.Client
	logger     Logger
}

// WithBaseURL overrides the default API base URL (https://api.levelfour.ai).
func WithBaseURL(baseURL string) ClientOption {
	return func(c *clientConfig) {
		c.baseURL = baseURL
	}
}

// WithMaxRetries sets the maximum number of retry attempts for failed requests.
// Defaults to 2. Retries apply to 429, 408, 409, and 5xx status codes.
func WithMaxRetries(n uint) ClientOption {
	return func(c *clientConfig) {
		c.maxRetries = n
		c.noRetries = false
	}
}

// WithNoRetries disables automatic retries entirely. By default the client
// retries up to 2 times on 429, 408, 409, and 5xx responses.
func WithNoRetries() ClientOption {
	return func(c *clientConfig) {
		c.noRetries = true
	}
}

// WithHTTPClient sets a custom http.Client for all API requests.
// When provided, the default 30-second timeout is not applied.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *clientConfig) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new LevelFour API client.
//
// The apiKey must start with "l4_live_" or "l4_test_". If empty, the
// LEVELFOUR_API_KEY environment variable is used as a fallback.
//
//	client, err := levelfour.NewClient("l4_live_...")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	summary, err := client.Recommendations.GetSavingsByProvider(ctx)
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		apiKey = os.Getenv("LEVELFOUR_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("levelfour: no API key provided, set LEVELFOUR_API_KEY or pass apiKey to NewClient")
	}

	valid := false
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(apiKey, prefix) {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("levelfour: invalid API key format, must start with one of: %s", strings.Join(validPrefixes, ", "))
	}

	cfg := &clientConfig{
		baseURL:    DefaultBaseURL,
		maxRetries: 2,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	httpClient := cfg.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	if cfg.logger != nil {
		base := httpClient.Transport
		if base == nil {
			base = http.DefaultTransport
		}
		httpClient = &http.Client{
			Transport: &loggingTransport{base: base, logger: cfg.logger},
			Timeout:   httpClient.Timeout,
		}
	}

	header := http.Header{}
	header.Set("Authorization", "Bearer "+apiKey)
	header.Set("User-Agent", "levelfour-go/"+Version+" "+runtime.Version())

	attempts := cfg.maxRetries + 1
	if cfg.noRetries {
		attempts = 1
	}

	clientOpts := []option.RequestOption{
		option.WithBaseURL(cfg.baseURL),
		option.WithHTTPHeader(header),
		option.WithMaxAttempts(attempts),
		option.WithHTTPClient(httpClient),
	}

	var retryOpts []internal.RetryOption
	if attempts > 0 {
		retryOpts = append(retryOpts, internal.WithMaxAttempts(attempts))
	}

	return &Client{
		Client:     baseclient.NewClient(clientOpts...),
		httpClient: httpClient,
		retrier:    internal.NewRetrier(retryOpts...),
		baseURL:    cfg.baseURL,
		header:     header.Clone(),
	}, nil
}
