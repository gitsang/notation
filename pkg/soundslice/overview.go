package soundslice

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type Score struct {
	SliceId    string `json:"sliceId,omitempty"`
	Name       string `json:"name,omitempty"`
	Hash       string `json:"hash,omitempty"`
	Embeddable bool   `json:"embeddable,omitempty"`
}

func (c *Client) ListScores() ([]*Score, error) {
	var (
		entryTime = time.Now()
		logger    = slog.New(c.logHandler)
	)
	defer func() {
		logger = logger.With(slog.String("cost", time.Since(entryTime).String()))
		logger.Debug("end")
	}()

	const (
		method = http.MethodGet
		path   = "/"
	)

	req, err := http.NewRequest(method, c.address+path, nil)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "sesn", Value: c.sesn})
	req.Header.Set("Referer", c.address)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.WithStack(fmt.Errorf("response status %s", resp.Status))
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		err = errors.WithStack(err)
		logger = logger.With(slog.Any("err", err))
		return nil, err
	}

	var results []*Score
	doc.Find("div.slice-item.slice-item-full").Each(func(i int, s *goquery.Selection) {
		sliceId, _ := s.Attr("data-score")
		name := s.Find("a.slice-item-title").Text()
		words := strings.Fields(name)
		name = strings.Join(words, " ")
		s.Find("div.slice-item-info").Each(func(j int, info *goquery.Selection) {
			info.Find("span.only10col").Each(func(k int, span *goquery.Selection) {
				results = append(results, &Score{
					SliceId:    sliceId,
					Name:       name,
					Hash:       "",
					Embeddable: span.Find("span.text-muted").Text() == "Embeddable",
				})
			})
		})
	})

	return results, nil
}
