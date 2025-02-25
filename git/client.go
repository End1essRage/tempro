package git

type Client struct {
	url string
}

func NewClient() *Client {
	return &Client{}
}
