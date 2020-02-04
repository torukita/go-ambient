package ambient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Config struct {
	ChannelId string
	ReadKey   string
	WriteKey  string
}

type Client struct {
	config Config
	client *http.Client
}

func NewClient(c Config) (*Client, error) {
	return &Client{
		config: c,
		client: &http.Client{},
	}, nil
}

func (c *Client) CreateData(ctx context.Context, b []byte) error {
	url := fmt.Sprintf("http://ambidata.io/api/v2/channels/%s/dataarray", c.config.ChannelId)
	req, err := postRequest(ctx, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer discardBody(resp)
	if resp.StatusCode != 200 {
		return fmt.Errorf("response code: %d\n", resp.StatusCode)
	}
	return nil
}

func postRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return req, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	return req, err
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func discardBody(resp *http.Response) error {
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
	return nil
}

/* データ取得
http.NewRequest("GET", fmt.Sprintf("http://ambidata.io/api/v2/channels/%s/data?readKey=%s&n=%d", c.config.ChannelId, c.config.ReadKey, n), nil)
*/

/* プロパティ取得
http.NewRequest("GET", fmt.Sprintf("http://ambidata.io/api/v2/channels/%s?readKey=%s", c.config.ChannelId, c.config.ReadKey), nil)
*/
