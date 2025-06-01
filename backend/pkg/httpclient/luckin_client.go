package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/internal/config"
)

type LuckinClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

func NewLuckinClient(cfg *config.LuckinConfig) *LuckinClient {
	return &LuckinClient{
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
		baseURL: cfg.BaseURL,
		headers: cfg.Headers,
	}
}

func (c *LuckinClient) doRequest(method, path string, body interface{}, referrer string) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	
	// 为queryByCard添加更完整的浏览器头部
	if path == "/api/api/v2/queryByCard" {
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Sec-Ch-Ua", `"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`)
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	}
	
	if referrer != "" {
		req.Header.Set("Referer", referrer)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}

// CheckGoodsCard 检查优惠卡状态
func (c *LuckinClient) CheckGoodsCard(card string) (*CheckGoodsCardResponse, error) {
	payload := map[string]string{
		"card":  card,
		"brand": "lk",
	}

	respBody, err := c.doRequest("POST", "/api/api/v2/checkGoodsCard", payload, fmt.Sprintf("%s/?card=%s", c.baseURL, card))
	if err != nil {
		return nil, err
	}

	var resp CheckGoodsCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// CityByCard 获取城市列表
func (c *LuckinClient) CityByCard(req *CityByCardRequest) (*CityByCardResponse, error) {
	respBody, err := c.doRequest("POST", "/api/api/v2/cityByCard", req, "https://lkcoffe.cn/city")
	if err != nil {
		return nil, err
	}

	var resp CityByCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// StoreByCard 根据卡片搜索门店
func (c *LuckinClient) StoreByCard(req *StoreByCardRequest) (*StoreByCardResponse, error) {
	respBody, err := c.doRequest("POST", "/api/api/v2/storeByCard", req, "")
	if err != nil {
		return nil, err
	}

	var resp StoreByCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// MenuByCard 获取菜单
func (c *LuckinClient) MenuByCard(req *MenuByCardRequest) (*MenuByCardResponse, error) {
	respBody, err := c.doRequest("POST", "/api/api/v2/menuByCard", req, "")
	if err != nil {
		return nil, err
	}

	var resp MenuByCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// GoodsByCard 获取商品详情
func (c *LuckinClient) GoodsByCard(req *GoodsByCardRequest) (*GoodsByCardResponse, error) {
	respBody, err := c.doRequest("POST", "/api/api/v2/goodsByCard", req, "")
	if err != nil {
		return nil, err
	}

	var resp GoodsByCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// CheckByCard 检查订单
func (c *LuckinClient) CheckByCard(req *CheckByCardRequest) (*CheckByCardResponse, error) {
	// 打印请求参数用于调试
	reqJSON, _ := json.Marshal(req)
	fmt.Printf("CheckByCard request: %s\n", string(reqJSON))
	
	respBody, err := c.doRequest("POST", "/api/api/v2/checkByCard", req, "")
	if err != nil {
		return nil, err
	}

	// 打印响应用于调试
	fmt.Printf("CheckByCard response: %s\n", string(respBody))

	var resp CheckByCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// OrderByCard 下单
func (c *LuckinClient) OrderByCard(req *OrderByCardRequest) (*OrderByCardResponse, error) {
	respBody, err := c.doRequest("POST", "/api/api/v2/orderByCard", req, "")
	if err != nil {
		return nil, err
	}

	var resp OrderByCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// QueryByCard 查询订单
func (c *LuckinClient) QueryByCard(req *QueryByCardRequest) (*QueryByCardResponse, error) {
	// 构建类似浏览器的referrer
	referrer := fmt.Sprintf("https://lkcoffe.cn/queryOrder?orderNo=%s&storeName=加载中&storeAddress=加载中", req.OrderNo)
	
	respBody, err := c.doRequest("POST", "/api/api/v2/queryByCard", req, referrer)
	if err != nil {
		return nil, err
	}

	// 添加原始响应日志用于调试
	fmt.Printf("QueryByCard raw response: %s\n", string(respBody))

	var resp QueryByCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}
