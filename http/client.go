package http

import (
	"encoding/json"
	gohttp "net/http"
	"time"
)

type Client struct {
	gohttp.Client
}

var DefaultClient = &Client{}

func init() {
	DefaultClient.Timeout = time.Second * 5
}

func (c *Client) GetAsObj(url string, obj interface{}) error {
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(obj)
}
