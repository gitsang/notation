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

func NewClient(opts ...ClientOptionFunc) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) CreateNotation() (sliceId string, err error) {
	var (
		method       = http.MethodPost
		path         = "/manage/create-via-import"
		responseBody = ListResponse{
			Items: &results,
		}
		err error
	)

}
