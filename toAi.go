package goToAi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DeepSeek API配置
const (
	apiKey    = "sk-7e3ff6b8f5704484921554b9e9b22f73"
	modelName = "deepseek-v4-flash"
	baseURL   = "https://api.deepseek.com/v1/chat/completions"
)

// 请求体结构
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 响应体结构
type ResponseBody struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ai请求
// AiRequest 首字母大写表示公开方法
func AiRequest() {
	fmt.Println("=== 使用标准net/http库访问DeepSeek API ===")
	fmt.Println()

	response, err := callDeepSeekWithHTTP("你好")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}

	fmt.Println("用户: 你好")
	fmt.Println("DeepSeek:", response)
	fmt.Println()
	fmt.Println("=== API使用统计 ===")
	fmt.Printf("模型: %s\n", modelName)
}

// 使用标准net/http库调用DeepSeek API
func callDeepSeekWithHTTP(message string) (string, error) {
	// 构建请求体
	requestBody := RequestBody{
		Model: modelName,
		Messages: []Message{
			{Role: "user", Content: message},
		},
	}

	// 将请求体序列化为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应体
	var responseBody ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取回复内容
	if len(responseBody.Choices) > 0 {
		return responseBody.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("响应中没有找到回复内容")
}
