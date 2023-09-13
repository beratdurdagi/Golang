package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/karalakrepp/Golang/BasicMicroservices/types"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}

}

func (c *Client) FetchPrice(ctx context.Context, ticker string) (*types.ClientResponse, error) {
	endpoint := fmt.Sprintf("%s?ticker=%s", c.endpoint, ticker)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResp := map[string]any{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("service responded with a non 200 status code: %s", errResp["error"])
	}

	priceResp := new(types.ClientResponse)
	if err := json.NewDecoder(resp.Body).Decode(priceResp); err != nil {
		return nil, err
	}
	resp.Body.Close()

	return priceResp, nil
}
