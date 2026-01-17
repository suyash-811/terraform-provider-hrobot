package client

import (
	"context"
	"net/http"
	"strings"
	"time"
)

type HetznerRobotClient struct {
	baseURL  string
	username string
	password string
	Client   *http.Client
}

func NewHetznerRobotClient(baseURL, username, password string) *HetznerRobotClient {
	return &HetznerRobotClient{
		baseURL:  strings.TrimRight(baseURL, "/"),
		username: username,
		password: password,
		Client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *HetznerRobotClient) doRequest(ctx context.Context, method string, uri string, data)