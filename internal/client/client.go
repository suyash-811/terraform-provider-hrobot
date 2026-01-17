package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func codeIsInExpected(statusCode int, expectedStatusCodes []int) bool {
	for _, expectedStatusCode := range expectedStatusCodes {
		if statusCode == expectedStatusCode {
			return true
		}
	}
	return false
}

func (c *HetznerRobotClient) doRequest(ctx context.Context, method string, path string, data url.Values, expectedStatusCodes []int) ([]byte, error) {
	fullURL := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimLeft(path, "/"))

	var body io.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	//var req *http.Request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)

	if data != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	response, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if !codeIsInExpected(response.StatusCode, expectedStatusCodes) {
		return nil, fmt.Errorf("The webservice returned an unexpted return code %d with response: %s", response.StatusCode, string(responseBytes))
	}

	return responseBytes, nil
}
