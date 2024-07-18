package soundslice

import (
	"log/slog"
	"net/http"
)

type Client struct {
	address  string
	username string
	password string
	sesn     string

	logHandler slog.Handler
	httpClient *http.Client
}

type ClientOptionFunc func(*Client)

func WithAddr(addr string) ClientOptionFunc {
	return func(c *Client) {
		c.address = addr
	}
}

func WithCredentials(username, password string) ClientOptionFunc {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

func WithSesn(sesn string) ClientOptionFunc {
	return func(c *Client) {
		c.sesn = sesn
	}
}

func WithLogHandler(h slog.Handler) ClientOptionFunc {
	return func(c *Client) {
		c.logHandler = h
	}
}

func WithHTTPClient(h *http.Client) ClientOptionFunc {
	return func(c *Client) {
		c.httpClient = h
	}
}

func defaultClient() *Client {
	return &Client{
		address:    "https://www.soundslice.com",
		username:   "",
		password:   "",
		sesn:       "",
		logHandler: nil,
		httpClient: http.DefaultClient,
	}
}

func NewClient(opts ...ClientOptionFunc) *Client {
	c := defaultClient()
	for _, opt := range opts {
		opt(c)
	}
	return c
}
