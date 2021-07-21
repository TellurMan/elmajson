package myjson

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseUrl    string
	ApiKeyName string
	ApiKey     string
	HTTPClient *http.Client
}

func NewClient(url, keyName, key string) *Client {
	return &Client{
		BaseUrl:    url,
		ApiKeyName: keyName,
		ApiKey:     key,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type Params map[string]string

func GetJSON(client Client, params Params, resultJSON interface{}) error {
	query := url.Values{}
	query.Add(client.ApiKeyName, client.ApiKey)
	for key, value := range params {
		query.Add(key, value)
	}
	baseUrl, _ := url.Parse(client.BaseUrl)
	baseUrl.RawQuery = query.Encode()

	req, _ := http.NewRequest("GET", baseUrl.String(), nil)

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error http: %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error status: %v", res.StatusCode)
	}

	var buf []byte
	res.Body.Read(buf)

	news := resultJSON

	defer res.Body.Close()

	json_err := json.NewDecoder(res.Body).Decode(&news)

	if json_err != nil {
		return fmt.Errorf("json error: %s", json_err.Error())
	}

	return nil
}
