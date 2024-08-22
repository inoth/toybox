package wssvr

import "net/http"

type Client struct {
	ID string
}

func (c *Client) Close() {}

func NewClient(w http.ResponseWriter, r *http.Request) *Client {
	return nil
}
