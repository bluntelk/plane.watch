package ws_client

var DefaultClient *Client

func init() {
	DefaultClient = NewClient("https://plane.watch")
}

func Connect() error {
	return DefaultClient.Connect()
}
