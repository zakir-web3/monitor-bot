package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	modelsEndpoint = "https://models.inference.ai.azure.com/chat/completions"
	modelName      = "gpt-4o-mini"
)

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func interpretRelease(token string, r *Release) (string, error) {
	prompt := fmt.Sprintf(`请用中文解读以下 go-ethereum 版本发布内容，简明扼要，分以下几点：
1. 版本概述
2. 主要新特性
3. 重要变更或不兼容改动
4. 安全修复（如有）
5. 对节点运营者的建议

版本：%s
发布内容：
%s`, r.TagName, r.Body)

	reqBody := chatRequest{
		Model: modelName,
		Messages: []message{
			{
				Role:    "system",
				Content: "你是一个以太坊和区块链技术专家，擅长解读技术发布说明。请用简洁清晰的中文回答，总字数控制在800字以内。",
			},
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", modelsEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github models api returned status %d", resp.StatusCode)
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return chatResp.Choices[0].Message.Content, nil
}
