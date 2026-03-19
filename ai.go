package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func callAI(ctx context.Context, token, sysPrompt, userPrompt string) (string, error) {
	reqBody := chatRequest{
		Model: modelName,
		Messages: []message{
			{Role: "system", Content: sysPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", modelsEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("github models api status %d: %s", resp.StatusCode, respBody)
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

func interpretRelease(ctx context.Context, token string, repo string, r *Release) (string, error) {
	prompt := fmt.Sprintf(userPromptTmpl, repo, r.TagName, r.Body)
	return callAI(ctx, token, systemPrompt, prompt)
}

func interpretReleaseDeep(ctx context.Context, token string, repo string, r *Release) (string, error) {
	prompt := fmt.Sprintf(deepUserPromptTmplZH, repo, r.TagName, r.Body)
	return callAI(ctx, token, deepSystemPromptZH, prompt)
}

func interpretReleaseDeepEN(ctx context.Context, token string, repo string, r *Release) (string, error) {
	prompt := fmt.Sprintf(deepUserPromptTmplEN, repo, r.TagName, r.Body)
	return callAI(ctx, token, deepSystemPromptEN, prompt)
}
