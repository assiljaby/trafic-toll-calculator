package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/assiljaby/trafic-toll-calculator/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) AggregateInvoice(distance types.Distance) error {
	b, err := json.Marshal(distance)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("this service responded with a non 200 status: %s", res.Status)
	}

	return nil
}
