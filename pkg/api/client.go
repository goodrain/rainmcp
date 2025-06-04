package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rainmcp/pkg/logger"
	"time"
)

// Client 是Rainbond API的客户端
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient 创建一个新的Rainbond API客户端
func NewClient(baseURL string) *Client {
	logger.Debug("创建新的API客户端: BaseURL=%s", baseURL)

	// 检查baseURL是否为空
	if baseURL == "" {
		logger.Warn("BaseURL为空")
		baseURL = "https://rainbond-api.example.com" // 设置一个默认值以避免空指针
	}

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get 发送GET请求到指定的API路径
func (c *Client) Get(path string) ([]byte, error) {
	logger.Debug("发送GET请求到: %s%s", c.BaseURL, path)
	return c.Request("GET", path, nil)
}

// Post 发送POST请求到指定的API路径
func (c *Client) Post(path string, data interface{}) ([]byte, error) {
	logger.Debug("发送POST请求到: %s%s", c.BaseURL, path)

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("序列化POST数据失败: %v", err)
		return nil, err
	}

	logger.Debug("POST请求数据: %s", string(jsonData))
	return c.Request("POST", path, bytes.NewBuffer(jsonData))
}

// Put 发送PUT请求到指定的API路径
func (c *Client) Put(path string, data interface{}) ([]byte, error) {
	logger.Debug("发送PUT请求到: %s%s", c.BaseURL, path)

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("序列化PUT数据失败: %v", err)
		return nil, fmt.Errorf("序列化PUT数据失败: %v", err)
	}

	return c.Request("PUT", path, bytes.NewBuffer(jsonData))
}

// Delete 发送DELETE请求到指定的API路径
func (c *Client) Delete(path string) ([]byte, error) {
	logger.Debug("发送DELETE请求到: %s%s", c.BaseURL, path)
	return c.Request("DELETE", path, nil)
}

// Request 发送请求到指定的API路径
func (c *Client) Request(method, path string, body io.Reader) ([]byte, error) {
	// 验证客户端是否正确初始化
	if c.BaseURL == "" {
		logger.Error("BaseURL为空")
		return nil, fmt.Errorf("API客户端未正确初始化，BaseURL为空")
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	logger.Debug("发送 %s 请求到: %s", method, url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logger.Error("创建请求失败: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		logger.Debug("添加授权头")
		req.Header.Set("Authorization", c.Token)
	} else {
		logger.Warn("未设置访问令牌")
	}

	logger.Debug("发送请求: %s %s", method, url)
	logger.Debug("请求头: %+v", req.Header)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		logger.Error("请求失败: %v", err)
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	logger.Debug("收到响应: 状态码=%d", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取响应体失败: %v", err)
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	logger.Debug("响应体: %s", string(respBody))

	if resp.StatusCode >= 400 {
		logger.Error("请求失败: 状态码=%d, 响应=%s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("API错误: %s, 状态码: %d", string(respBody), resp.StatusCode)
	}

	logger.Debug("请求成功: 状态码=%d", resp.StatusCode)
	return respBody, nil
}
