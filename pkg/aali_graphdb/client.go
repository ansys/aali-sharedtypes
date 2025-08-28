// Copyright (C) 2025 ANSYS, Inc. and/or its affiliates.
// SPDX-License-Identifier: MIT
//
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package aali_graphdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

type Client struct {
	address    string
	logger     *zap.Logger
	httpClient *http.Client
}

func NewClient(address string, httpClient *http.Client) (*Client, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync() //nolint:errcheck
	return &Client{address, logger, httpClient}, nil
}

func DefaultClient(address string) (*Client, error) {
	return NewClient(address, http.DefaultClient)
}

func (client Client) post(u string, body any) (*http.Response, error) {
	jsonReq, err := json.Marshal(body)
	if err != nil {
		return nil, err

	}
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "*/*")
	if err != nil {
		return nil, err
	}

	return client.httpClient.Do(req)
}

func (client Client) GetHealth() (bool, error) {
	url, err := url.JoinPath(client.address, "health")
	if err != nil {
		return false, err
	}
	resp, err := client.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			client.logger.Warn("could not close body")
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return true, nil
}

type getVersionResponse struct {
	Status             string `json:"status"`
	Version            string `json:"version"`
	KuzuVersion        string `json:"kuzu_version"`
	KuzuStorageVersion int    `json:"kuzu_storage_version"`
}

func (client Client) GetVersion() (getVersionResponse, error) {
	url, err := url.JoinPath(client.address, "health")
	if err != nil {
		return getVersionResponse{}, err
	}
	resp, err := client.httpClient.Get(url)
	if err != nil {
		return getVersionResponse{}, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			client.logger.Warn("could not close body")
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return getVersionResponse{}, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	var r getVersionResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return getVersionResponse{}, err
	}

	return r, nil
}

type getDatabasesResponse struct {
	Databases []string `json:"databases"`
}

func (client Client) GetDatabases() ([]string, error) {
	url, err := url.JoinPath(client.address, "databases")
	if err != nil {
		return nil, err
	}
	resp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			client.logger.Warn("could not close body")
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	var r getDatabasesResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.Databases, nil
}

func (client Client) CreateDatabase(name string) error {
	u, err := url.JoinPath(client.address, "databases")
	if err != nil {
		return err
	}

	resp, err := client.post(u, map[string]any{"name": name, "in_memory": false})
	if err != nil {
		return err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			client.logger.Warn("could not close body")
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return nil
}

func (client Client) DeleteDatabase(name string) error {
	u, err := url.JoinPath(client.address, "databases", name)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			client.logger.Warn("could not close body")
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return nil
}

type cypherQueryResponse[T any] struct {
	Result []T `json:"result"`
}

func CypherQueryReadGeneric[T any](client *Client, db string, cypher string, parameters Parameters) ([]T, error) {
	u, err := url.JoinPath(client.address, "databases", db, "read")
	if err != nil {
		return nil, err
	}

	var params map[string]Value
	if parameters != nil {
		params, err = parameters.AsParameters()
		if err != nil {
			return nil, err
		}
	}

	resp, err := client.post(u, map[string]any{"cypher": cypher, "parameters": params})
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			client.logger.Warn("could not close body")
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		body := string(bodyBytes)
		return nil, fmt.Errorf("unexpected status code: %v %q", resp.StatusCode, body)
	}

	var r cypherQueryResponse[T]
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.Result, nil
}

func (client *Client) CypherQueryRead(db string, cypher string, parameters Parameters) ([]map[string]any, error) {
	return CypherQueryReadGeneric[map[string]any](client, db, cypher, parameters)
}

func CypherQueryWriteGeneric[T any](client *Client, db string, cypher string, parameters Parameters) ([]T, error) {
	u, err := url.JoinPath(client.address, "databases", db, "write")
	if err != nil {
		return nil, err
	}

	var params map[string]Value
	if parameters != nil {
		params, err = parameters.AsParameters()
		if err != nil {
			return nil, err
		}
	}

	resp, err := client.post(u, map[string]any{"cypher": cypher, "parameters": params})
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			client.logger.Warn("could not close body")
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		body := string(bodyBytes)
		return nil, fmt.Errorf("unexpected status code: %v %q", resp.StatusCode, body)
	}

	var r cypherQueryResponse[T]
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.Result, nil
}

func (client *Client) CypherQueryWrite(db string, cypher string, parameters Parameters) ([]map[string]any, error) {
	return CypherQueryWriteGeneric[map[string]any](client, db, cypher, parameters)
}

type ParameterMap map[string]Value

func (pm *ParameterMap) UnmarshalJSON(data []byte) error {
	var intermediate map[string]valueUnmarshalHelper
	err := json.Unmarshal(data, &intermediate)
	if err != nil {
		return err
	}

	paramMap := make(map[string]Value, len(intermediate))
	for k, v := range intermediate {
		paramMap[k] = v.Value
	}
	*pm = ParameterMap(paramMap)
	return nil
}

type Parameters interface {
	AsParameters() (map[string]Value, error)
}

func (params ParameterMap) AsParameters() (map[string]Value, error) {
	return params, nil
}
