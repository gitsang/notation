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
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

func (c *Client) CreateNotation(_ context.Context) (sliceId string, err error) {
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

func (c *Client) UploadNotation(_ context.Context, sliceId string, fh *os.File) (*UploadNotationResponse, error) {
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

	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	defer writer.Close()
	formWriter, err := writer.CreateFormFile("score", filepath.Base(fh.Name()))
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}
	_, err = io.Copy(formWriter, fh)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}
	writer.WriteField("type", "application/octet-stream")

	req, err := http.NewRequest(method, c.address+path, reqBody)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "sesn", Value: c.sesn})
	req.Header.Set("Referer", c.address)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		err = errors.WithStack(fmt.Errorf("response status %s: %s", resp.Status, string(respBody)))
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}

	var result UploadNotationResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteNotation(_ context.Context, sliceId string) error {
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
		path   = "/api/v1/slices/delete-multiple/"
	)

	formData := url.Values{}
	formData.Set("ids", sliceId)

	req, err := http.NewRequest(method, c.address+path, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return err
	}
	req.AddCookie(&http.Cookie{Name: "sesn", Value: c.sesn})
	req.Header.Set("Referer", c.address)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.WithStack(fmt.Errorf("response status %s", resp.Status))
		logger = logger.With(slog.Any("err", err))
		return err
	}

	return nil
}
