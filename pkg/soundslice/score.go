package soundslice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type ScoreData struct {
	Slug string `json:"slug"`
}

func (c *Client) GetScoreId(sliceId string) (string, error) {
	var (
		entryTime = time.Now()
		logger    = slog.New(c.logHandler)
	)
	defer func() {
		logger = logger.With(slog.String("cost", time.Since(entryTime).String()))
		logger.Debug("end")
	}()

	var (
		method = http.MethodGet
		path   = fmt.Sprintf("/slices/%s/edit/scoredata/", sliceId)
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

	var data ScoreData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return "", err
	}

	return data.Slug, nil
}

func (c *Client) EnableEmbed(scoreId string) error {
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
		path   = fmt.Sprintf("/api/v1/scores/%s/", scoreId)
	)

	formData := url.Values{}
	formData.Set("embed_status", "4")

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

func (c *Client) DisableEmbed(scoreId string) error {
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
		path   = fmt.Sprintf("/api/v1/scores/%s/", scoreId)
	)

	formData := url.Values{}
	formData.Set("embed_status", "1")

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
