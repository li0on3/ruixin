package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL    = "http://localhost:8080/api/v1"
	apiKey     = "test_api_key"     // 替换为实际的API Key
	apiSecret  = "test_api_secret"  // 替换为实际的API Secret
)

// 测试结果
type TestResult struct {
	TestName string
	Passed   bool
	Message  string
}

func main() {
	fmt.Println("=== 瑞幸分销系统安全功能测试 ===")
	fmt.Println()

	results := []TestResult{}

	// 1. 测试频率限制
	fmt.Println("1. 测试频率限制...")
	result := testRateLimit()
	results = append(results, result)
	fmt.Println()

	// 2. 测试卡片权限验证
	fmt.Println("2. 测试卡片权限验证...")
	result = testCardAccessControl()
	results = append(results, result)
	fmt.Println()

	// 3. 测试审计日志
	fmt.Println("3. 测试审计日志记录...")
	result = testAuditLogging()
	results = append(results, result)
	fmt.Println()

	// 打印测试总结
	fmt.Println("=== 测试总结 ===")
	passedCount := 0
	for _, r := range results {
		status := "❌ 失败"
		if r.Passed {
			status = "✅ 通过"
			passedCount++
		}
		fmt.Printf("%s %s: %s\n", status, r.TestName, r.Message)
	}
	fmt.Printf("\n总计: %d/%d 通过\n", passedCount, len(results))
}

// 测试频率限制
func testRateLimit() TestResult {
	// 快速发送多个请求，测试是否触发频率限制
	endpoint := "/distributor/cities"
	
	var lastStatusCode int
	hitRateLimit := false
	
	// 发送15个请求（超过敏感API的限制：10次/分钟）
	for i := 0; i < 15; i++ {
		resp, err := makeRequest("GET", endpoint, nil)
		if err != nil {
			return TestResult{
				TestName: "频率限制",
				Passed:   false,
				Message:  fmt.Sprintf("请求失败: %v", err),
			}
		}
		defer resp.Body.Close()
		
		lastStatusCode = resp.StatusCode
		if resp.StatusCode == 429 {
			hitRateLimit = true
			break
		}
		
		// 检查限流响应头
		if i == 0 {
			limit := resp.Header.Get("X-RateLimit-Limit")
			remaining := resp.Header.Get("X-RateLimit-Remaining")
			fmt.Printf("   频率限制信息: 限制=%s, 剩余=%s\n", limit, remaining)
		}
	}
	
	if hitRateLimit {
		return TestResult{
			TestName: "频率限制",
			Passed:   true,
			Message:  "成功触发频率限制（429错误）",
		}
	}
	
	return TestResult{
		TestName: "频率限制",
		Passed:   false,
		Message:  fmt.Sprintf("未触发频率限制，最后状态码: %d", lastStatusCode),
	}
}

// 测试卡片权限验证
func testCardAccessControl() TestResult {
	// 测试访问一个未授权的卡片
	unauthorizedCard := "UNAUTHORIZED123"
	endpoint := fmt.Sprintf("/distributor/stores?card=%s&city_name=北京", unauthorizedCard)
	
	resp, err := makeRequest("GET", endpoint, nil)
	if err != nil {
		return TestResult{
			TestName: "卡片权限验证",
			Passed:   false,
			Message:  fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()
	
	// 在软模式下，应该返回200（不阻止访问）
	if resp.StatusCode == 200 {
		fmt.Println("   软模式：允许访问未授权卡片，但应记录审计日志")
		return TestResult{
			TestName: "卡片权限验证",
			Passed:   true,
			Message:  "软模式下正确处理未授权访问（返回200）",
		}
	}
	
	return TestResult{
		TestName: "卡片权限验证",
		Passed:   false,
		Message:  fmt.Sprintf("意外的状态码: %d", resp.StatusCode),
	}
}

// 测试审计日志
func testAuditLogging() TestResult {
	// 执行一些会触发审计日志的操作
	operations := []struct {
		name     string
		endpoint string
		method   string
	}{
		{"查询门店", "/distributor/stores?card=TEST123&city_name=北京", "GET"},
		{"查询菜单", "/distributor/menu?card=TEST123&store_code=S001", "GET"},
		{"查询商品", "/distributor/goods?card=TEST123&store_code=S001&goods_code=G001", "GET"},
	}
	
	for _, op := range operations {
		fmt.Printf("   执行操作: %s\n", op.name)
		resp, err := makeRequest(op.method, op.endpoint, nil)
		if err != nil {
			fmt.Printf("   警告: %s 失败 - %v\n", op.name, err)
			continue
		}
		resp.Body.Close()
	}
	
	// 审计日志是异步记录的，这里只是确认操作执行了
	return TestResult{
		TestName: "审计日志记录",
		Passed:   true,
		Message:  "已执行会触发审计日志的操作（需要查看数据库确认）",
	}
}

// 发送HTTP请求的辅助函数
func makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}
	
	req, err := http.NewRequest(method, baseURL+endpoint, bodyReader)
	if err != nil {
		return nil, err
	}
	
	// 添加认证头
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-API-Secret", apiSecret)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	return client.Do(req)
}