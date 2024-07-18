package soundslice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

func (c *Client) CreateNotation() (sliceId string, err error) {
	var (
		entryTime = time.Now()
		logger    = slog.New(c.logHandler)
	)
	defer func() {
		logger = logger.With(slog.String("cost", time.Since(entryTime).String()))
		logger.Debug("end")
	}()

	const (
		method = http.MethodPost
		path   = "/manage/create-via-import/"
	)

	req, err := http.NewRequest(method, c.address+path, nil)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return "", err
	}
	req.AddCookie(&http.Cookie{Name: "sesn", Value: c.sesn})
	req.Header.Set("Referer", c.address)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.WithStack(fmt.Errorf("response status %s", resp.Status))
		logger = logger.With(slog.Any("err", err))
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return "", err
	}

	sliceId, exists := doc.Find("slice-practice-lists[id='title-practice-lists']").Attr("slice")
	if !exists {
		err = errors.WithStack(errors.New("slice id not found"))
		logger = logger.With(slog.Any("err", err))
		return "", err
	}

	return sliceId, nil
}

type UploadNotationResponse struct {
	Name string `json:"name"`
}

func (c *Client) UploadNotation(ctx context.Context, sliceId string, filename string) (*UploadNotationResponse, error) {
	var (
		entryTime = time.Now()
		logger    = slog.New(c.logHandler)
	)
	defer func() {
		logger = logger.With(slog.String("cost", time.Since(entryTime).String()))
		logger.Debug("end")
	}()

	var (
		method = http.MethodPost
		path   = fmt.Sprintf("/api/v1/slices/%s/notation/", sliceId)
	)

	// new multipart request
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	defer writer.Close()

	// create form file
	formWriter, err := writer.CreateFormFile("score", filename)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}

	// open file
	fh, err := os.Open(filename)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}
	defer fh.Close()

	// copy body
	_, err = io.Copy(formWriter, fh)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}

	// set form type
	writer.WriteField("type", "application/octet-stream")

	// new request
	req, err := http.NewRequest(method, c.address+path, reqBody)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// authenticate
	req.AddCookie(&http.Cookie{Name: "sesn", Value: c.sesn})
	req.Header.Set("Referer", c.address)

	// do request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}
	defer resp.Body.Close()

	// check response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		err = errors.WithStack(fmt.Errorf("response status %s: %s", resp.Status, string(respBody)))
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}

	// decode response
	var result UploadNotationResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}

	return &result, nil
}
