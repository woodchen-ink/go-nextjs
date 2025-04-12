package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-nextjs/config"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// VPSConfig 表示从描述中提取的VPS配置信息
type VPSConfig struct {
	CPU       string `json:"cpu"`
	RAM       string `json:"ram"`
	Disk      string `json:"disk"`
	Bandwidth string `json:"bandwidth"`
	IP        string `json:"ip"`
	Location  string `json:"location"`
	Remark    string `json:"remark"`
}

// Message 表示OpenAI API的消息格式
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 表示OpenAI聊天请求
type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

// ChatResponse 表示OpenAI聊天响应
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// ParseVPSDescription 使用AI分析VPS描述，提取配置信息
func ParseVPSDescription(description string) (*VPSConfig, error) {
	// 检查API配置
	if config.AIAPIKey == "" {
		return nil, fmt.Errorf("未配置AI API密钥")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 准备对AI的提示
	systemPrompt := `你是一个专门提取VPS配置信息的AI助手。请分析提供的VPS描述文本，提取以下信息：
CPU核心数、内存大小、硬盘容量及类型、带宽/流量、IP地址信息、服务器位置和其他需要注意的备注事项。
请以JSON格式返回结果，包含以下字段：cpu, ram, disk, bandwidth, ip, location, remark。
对于无法确定的字段，请使用空字符串。

提取的基本配置信息应该尽量简洁、标准化，例如：
- CPU: "2 Core" 而不是 "2x Intel CPU"
- RAM: "4GB RAM" 而不是 "4GB Memory"
- Disk: "50GB SSD" 而不是 "50 Gigabytes Solid State Drive"

对于带宽信息，请使用以下格式：
- 如果有明确的流量限制和带宽速度，使用"流量@带宽速度"格式，例如："500GB Traffic@1Gbps port"
- 如果只有流量限制，例如："500GB Traffic"
- 如果是无限流量，则使用"Unlimited Traffic@带宽速度"或者简单的"Unlimited Traffic"
- 端口速度单位统一使用Gbps或Mbps

对于IP信息，请标准化为：
- 如果只有IPv4，例如："1 IPv4"或"2 IPv4"
- 如果同时有IPv4和IPv6，例如："IPv4 + IPv6"或"2 IPv4 + IPv6"
- 如果有特殊情况，请清晰描述，如："1 专用IPv4 + IPv6 子网"

对于位置信息，请保留完整的数据中心信息：
- 包括数据中心代号，如"Singapore DC1"、"Hong Kong DC2"
- 保留国家/地区代码，如"Tokyo, JP"、"Hanoi, VN"
- 如果有多个位置，请使用逗号分隔，保留原始信息的完整性
- 例如："Singapore DC1, Hong Kong DC2, Tokyo JP, Hanoi VN"

对于备注信息，请重点关注：
1. 促销优惠内容和条件
2. 特殊功能或限制（如备份、快照、DDoS防护、私有网络等）
3. 操作系统或面板相关信息
4. 支持的虚拟化技术
5. 任何可能影响用户体验的重要说明
6. 进行汉化

请将所有不属于CPU、内存、硬盘、带宽、IP、位置这六个基本字段的重要信息整理到备注(remark)字段中。
只返回JSON数据，不要有其他文字。`

	// 分析描述
	userPrompt := fmt.Sprintf("VPS描述: %s", description)
	log.Printf("AI分析VPS描述: %s", userPrompt)

	// 准备API请求
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	chatReq := ChatRequest{
		Model:     config.AIModel,
		Messages:  messages,
		MaxTokens: 500,
	}

	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		return nil, fmt.Errorf("序列化AI请求失败: %v", err)
	}

	// 创建API请求
	endpoint := fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(config.AIURL, "/"))
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建AI请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.AIAPIKey))

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("AI请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI服务返回非成功状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取AI响应失败: %v", err)
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("AI响应中没有选择项")
	}

	// 提取AI返回的JSON
	content := chatResp.Choices[0].Message.Content
	// 清理可能的markdown格式
	content = cleanJSONFromMarkdown(content)

	var config VPSConfig
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		return nil, fmt.Errorf("解析AI返回的JSON失败: %v", err)
	}

	return &config, nil
}

// cleanJSONFromMarkdown 从可能包含markdown格式的内容中提取JSON
func cleanJSONFromMarkdown(content string) string {
	// 处理可能的代码块
	if strings.Contains(content, "```json") && strings.Contains(content, "```") {
		parts := strings.Split(content, "```")
		for i, part := range parts {
			if strings.HasPrefix(part, "json") || i > 0 && i < len(parts)-1 {
				return strings.TrimPrefix(part, "json")
			}
		}
	}

	// 如果没有markdown格式但内容以{开始以}结束，直接返回
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
		return content
	}

	return content
}

// OptimizeTitle 使用AI优化VPS标题，进行汉化和长度优化
func OptimizeTitle(title string) (string, error) {
	// 检查API配置
	if config.AIAPIKey == "" {
		return title, fmt.Errorf("未配置AI API密钥")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 准备对AI的提示
	systemPrompt := `你是一个专门优化VPS标题的AI助手。请对输入的VPS标题进行如下优化：
1. 如果标题中包含英文描述，尝试将其汉化为更易于中文用户理解的形式
2. 如果标题过长（超过30个字符），进行适当缩短，但保留关键信息
3. 保留原标题中的规格信息，如CPU核心数、内存大小、硬盘容量等
4. 保留原标题中的特殊优惠或促销信息
5. 保留品牌名称，不要汉化品牌名
6. 如果原标题已经简洁且为中文，则无需更改

请直接返回优化后的标题，不要包含任何解释或额外文字。如果标题已经符合要求或无法优化，则返回原标题。`

	// 分析标题
	userPrompt := fmt.Sprintf("原标题: %s", title)

	// 准备API请求
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	chatReq := ChatRequest{
		Model:     config.AIModel,
		Messages:  messages,
		MaxTokens: 100,
	}

	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		return title, fmt.Errorf("序列化AI请求失败: %v", err)
	}

	// 创建API请求
	endpoint := fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(config.AIURL, "/"))
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return title, fmt.Errorf("创建AI请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.AIAPIKey))

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return title, fmt.Errorf("AI请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return title, fmt.Errorf("AI服务返回非成功状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return title, fmt.Errorf("读取AI响应失败: %v", err)
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return title, fmt.Errorf("解析AI响应失败: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return title, fmt.Errorf("AI响应中没有选择项")
	}

	// 提取AI返回的优化标题
	optimizedTitle := strings.TrimSpace(chatResp.Choices[0].Message.Content)

	// 如果AI返回空标题或与原标题相同，使用原标题
	if optimizedTitle == "" || optimizedTitle == title {
		return title, nil
	}

	return optimizedTitle, nil
}

// OptimizeTitleAsync 异步优化多个VPS标题
func OptimizeTitleAsync(titles []string, callback func(int, string)) {
	// 创建工作池
	const maxWorkers = 300
	workChan := make(chan int, len(titles))
	var wg sync.WaitGroup

	// 启动工作协程
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range workChan {
				optimizedTitle, err := OptimizeTitle(titles[idx])
				if err != nil {
					log.Printf("优化标题 %s 失败: %v", titles[idx], err)
					continue
				}

				// 回调函数处理优化后的标题
				if callback != nil {
					callback(idx, optimizedTitle)
				}
			}
		}()
	}

	// 发送工作
	for i := range titles {
		workChan <- i
	}
	close(workChan)

	// 等待所有工作完成
	wg.Wait()
}
