package iconikio

import (
	"io"
	"log"
	"net/http"
)

// GetResponseBody is a helper function for making API calls.
func GetResponseBody(method, uri string, params io.Reader, c *Client) (*http.Response, []byte, error) {
	log.Println(uri)

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, params)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("App-ID", c.cfg.AppID)
	req.Header.Add("Auth-Token", c.cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return res, resBody, nil
}
