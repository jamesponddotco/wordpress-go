package wordpress

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"git.sr.ht/~jamesponddotco/httpx-go"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"golang.org/x/time/rate"
)

const (
	// ErrConfigRequired is returned when a Client is created without a Config.
	ErrConfigRequired xerrors.Error = "config cannot be empty"
)

type (
	// Service is a common struct that can be reused instead of allocating a new
	// one for each service on the heap.
	service struct {
		client *Client
	}

	// Client is a client for the WordPress REST API.
	Client struct {
		// httpc is the underlying HTTP client used by the API client.
		httpc *httpx.Client

		// cfg specifies the configuration used by the API client.
		cfg *Config

		// common service fields shared by all services.
		common service
	}
)

// NewClient returns a new client for the WordPress REST API.
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, ErrConfigRequired
	}

	cfg.init()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	c := &Client{
		httpc: &httpx.Client{
			RateLimiter: rate.NewLimiter(rate.Limit(2), 1),
			RetryPolicy: httpx.DefaultRetryPolicy(),
			UserAgent:   cfg.Application.UserAgent(),
		},
		cfg: cfg,
	}

	c.common.client = c

	return c, nil
}

// Response represents a response from the WordPress API.
type Response struct {
	// Header contains the response headers.
	Header http.Header

	// Body contains the response body as a byte slice.
	Body []byte

	// Status is the HTTP status code of the response.
	Status int
}

// do performs an HTTP request using the underlying HTTP client.
func (c *Client) do(ctx context.Context, req *http.Request) (*Response, error) {
	if c.cfg.Debug {
		c.cfg.Logger.Printf("request: %s %s", req.Method, req.URL)
	}

	ret, err := c.httpc.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	defer func() {
		if err = httpx.DrainResponseBody(ret); err != nil {
			log.Fatal(err)
		}
	}()

	if c.cfg.Debug {
		var dump []byte

		dump, err = httputil.DumpResponse(ret, true)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		c.cfg.Logger.Printf("\n%s", dump)
	}

	body, err := io.ReadAll(ret.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	response := &Response{
		Header: ret.Header.Clone(),
		Body:   body,
		Status: ret.StatusCode,
	}

	return response, nil
}

// request is a convenience function for creating an HTTP request.
func (c *Client) request(
	ctx context.Context,
	method, uri string,
	headers map[string]string,
	body io.Reader,
) (*http.Request, error) {
	if _, ok := headers["User-Agent"]; !ok {
		headers["User-Agent"] = c.cfg.Application.UserAgent().String()
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if c.cfg.Debug {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		c.cfg.Logger.Printf("\n%s", string(dump))
	}

	return req, nil
}
