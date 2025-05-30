package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"yadro.com/course/update/core"
)

type Client struct {
	log    *slog.Logger
	client http.Client
	url    string
}

const uriPath = "/info.0.json"

func NewClient(url string, timeout time.Duration, log *slog.Logger) (*Client, error) {
	if url == "" {
		return nil, fmt.Errorf("empty base url specified")
	}
	return &Client{
		client: http.Client{Timeout: timeout},
		log:    log,
		url:    url,
	}, nil
}

func fetchComic(url string) (core.XKCDInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return core.XKCDInfo{}, fmt.Errorf("failed to make a comic request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return core.XKCDInfo{}, core.ErrNotFound
	}
	info := struct {
		ID         int    `json:"num"`
		URL        string `json:"img"`
		Title      string `json:"title"`
		SafeTitle  string `json:"safe_title"`
		Transcript string `json:"transcript"`
		Alt        string `json:"alt"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return core.XKCDInfo{}, err
	}

	return core.XKCDInfo{
		ID:    info.ID,
		URL:   info.URL,
		Title: info.Title,
		Description: strings.Join([]string{
			info.Title, info.SafeTitle, info.Transcript, info.Alt},
			" "),
	}, nil
}

func (c Client) Get(ctx context.Context, id int) (core.XKCDInfo, error) {
	resp, err := fetchComic(c.url + fmt.Sprintf("/%d/%s", id, uriPath))
	if err != nil {
		c.log.Error("Could not fetch comic", "error", err)
		return resp, err
	}
	return resp, nil
}

func (c Client) LastID(ctx context.Context) (int, error) {
	resp, err := fetchComic(c.url + fmt.Sprintf("/%s", uriPath))
	if err != nil {
		c.log.Error("Could not fetch comic", "error", err)
		return 0, err
	}
	return resp.ID, nil
}
