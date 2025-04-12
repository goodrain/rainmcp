package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 是Rainbond API的客户端
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient 创建一个新的Rainbond API客户端
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get 发送GET请求到指定的API路径
func (c *Client) Get(path string) ([]byte, error) {
	return c.Request("GET", path, nil)
}

// Post 发送POST请求到指定的API路径
func (c *Client) Post(path string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.Request("POST", path, bytes.NewBuffer(jsonData))
}

// Request 发送请求到指定的API路径
func (c *Client) Request(method, path string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API错误: %s, 状态码: %d", string(respBody), resp.StatusCode)
	}
	
	return respBody, nil
}
