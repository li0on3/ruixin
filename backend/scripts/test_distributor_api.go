package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"backend/internal/config"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Token string `json:"token"`
		User  struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	} `json:"data"`
}

type DistributorResponse struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		List []struct {
			ID      int     `json:"id"`
			Name    string  `json:"name"`
			Email   string  `json:"email"`
			Status  int     `json:"status"`
			Balance float64 `json:"balance"`
		} `json:"list"`
		Total    int `json:"total"`
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	} `json:"data"`
}

func main() {
	// 加载配置
	if err := config.Init("./configs/config.yaml"); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	baseURL := "http://localhost:8080"
	
	// 1. 测试登录
	fmt.Println("1. 测试管理员登录...")
	token, err := testLogin(baseURL)
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}
	fmt.Printf("登录成功，获取到Token: %s...\n", token[:20])

	// 2. 测试获取分销商列表
	fmt.Println("\n2. 测试获取分销商列表...")
	err = testGetDistributors(baseURL, token)
	if err != nil {
		fmt.Printf("获取分销商列表失败: %v\n", err)
		return
	}

	fmt.Println("\n✅ API测试全部通过！")
}

func testLogin(baseURL string) (string, error) {
	loginReq := LoginRequest{
		Username: "admin",
		Password: "admin123",
	}

	reqBody, _ := json.Marshal(loginReq)
	
	resp, err := http.Post(baseURL+"/api/v1/admin/login", "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if loginResp.Code != 200 {
		return "", fmt.Errorf("登录失败: %s", loginResp.Msg)
	}

	return loginResp.Data.Token, nil
}

func testGetDistributors(baseURL, token string) error {
	req, err := http.NewRequest("GET", baseURL+"/api/v1/admin/distributors?page=1&page_size=100", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容: %s\n", string(body))

	var distResp DistributorResponse
	if err := json.Unmarshal(body, &distResp); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if distResp.Code != 200 {
		return fmt.Errorf("API返回错误: %s", distResp.Msg)
	}

	fmt.Printf("成功获取分销商列表: 总计%d个分销商\n", distResp.Data.Total)
	for _, distributor := range distResp.Data.List {
		fmt.Printf("- ID: %d, Name: %s, Email: %s, Status: %d, Balance: %.2f\n", 
			distributor.ID, distributor.Name, distributor.Email, distributor.Status, distributor.Balance)
	}

	return nil
}