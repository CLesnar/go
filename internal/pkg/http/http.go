package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
)

type HttpI interface {
	GetRequest(criteria interface{}, parameters interface{}) interface{}
}

type Http struct {
}

func ModelMap(model interface{}) (map[string]string, error) {
	modelJson, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	modelMap := map[string]string{}
	if err := json.Unmarshal(modelJson, &modelMap); err != nil {
		return nil, err
	}
	return modelMap, nil
}

func AddQueryParams(parameters interface{}, query url.Values) error {
	if query == nil {
		return errors.New("query param cannot be nil")
	}
	paramMap, err := ModelMap(parameters)
	if err != nil {
		return err
	}
	for k, v := range paramMap {
		val := fmt.Sprintf("%v", v) // convert interface{} to string
		query.Add(k, val)
	}
	return nil
}

func (h *Http) GetRequest(ctx context.Context, url string, criteria interface{}, parameters interface{}, headers map[string]string) ([]byte, error) {
	var body io.Reader
	if criteria != nil {
		criteriaBytes, err := json.Marshal(criteria)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal criteria body: %v", err)
		}
		body = bytes.NewBuffer(criteriaBytes)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	q := req.URL.Query()
	if err := AddQueryParams(parameters, q); err != nil {
		return nil, fmt.Errorf("failed to add path parameters: %v", err)
	}
	req.URL.RawQuery = req.URL.Query().Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("faile to do request: %v", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return respBody, fmt.Errorf("failed to read response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("response status not okay '200', recieved '%v'", resp.StatusCode)
	}
	return respBody, nil
}