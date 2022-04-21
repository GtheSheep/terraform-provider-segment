package segment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type AuthResponseData struct {
	Workspace Workspace `json:"workspace"`
}

type AuthResponse struct {
	Data AuthResponseData `json:"data`
}

func NewClient(account_id *int, token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
		Token:      *token,
	}

	if token != nil {
		url := fmt.Sprintf("%s", HostURL)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req)

		ar := AuthResponse{}
		err = json.Unmarshal(body, &ar)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != 201) {
		return nil, fmt.Errorf("%s url: %s, status: %d, body: %s", req.Method, req.URL, res.StatusCode, body)
	}

	return body, err
}
