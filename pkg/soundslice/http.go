package soundslice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type RequestOptions struct {
	Method  string
	Path    string
	Query   url.Values
	Headers http.Header
	Body    any
	Timeout time.Duration
}

func (c *Client) HTTPRequest(ctx context.Context, opts RequestOptions, respBody any) error {
	var (
		err       error
		entryTime = time.Now()
		logger    = slog.New(c.logHandler).With(
			slog.Any("opts", opts),
		)
	)

	defer func() {
		logger = logger.With(slog.String("cost", time.Since(entryTime).String()))
		logger.Debug("end")
	}()

	// request url
	u, err := url.JoinPath(c.address, opts.Path)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "invalid request url"))
	}

	// request body
	var body io.Reader
	if opts.Body != nil {
		bodyBytes, err := json.Marshal(opts.Body)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "invalid request body"))
		}
		body = bytes.NewBuffer([]byte(bodyBytes))
	}

	// new request
	req, err := http.NewRequest(opts.Method, u, body)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "invalid request"))
	}

	// request query
	if opts.Query != nil {
		req.URL.RawQuery = opts.Query.Encode()
	}

	// request header
	if opts.Headers != nil {
		req.Header = opts.Headers
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json;charset=utf8")
	}

	// authorize
	req.AddCookie(&http.Cookie{
		Name:  "sesn",
		Value: c.sesn,
	})

	// do request
	logger = logger.With(
		slog.Any("request.method", req.Method),
		slog.Any("request.url", req.URL.String()),
		slog.Any("request.headers", req.Header),
		slog.Any("request.query", req.URL.RawQuery),
		slog.Any("request.body", opts.Body),
	)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = errors.WithStack(errors.Wrap(err, "request failed"))
		logger = logger.With(slog.Any("err", err))
		return err
	}
	logger = logger.With(slog.Int("response.status", resp.StatusCode))

	// response body
	defer resp.Body.Close()
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.WithStack(errors.Wrap(err, "response body read failed"))
		logger = logger.With(slog.Any("err", err))
		return err
	}
	logger = logger.With(slog.String("response.body", string(respBodyBytes)))

	// response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("response status %s: %s", resp.Status, string(respBodyBytes))
		logger = logger.With(slog.Any("err", err))
		return err
	}

	// parse response
	if err := json.Unmarshal(respBodyBytes, &respBody); err != nil {
		err = errors.WithStack(errors.Wrap(err, "response body unmarshal failed"))
		logger = logger.With(slog.Any("err", err))
		return err
	}

	return nil
}
