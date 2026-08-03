package creative

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const namingPromptVersion = "genesis-naming-v1"

var hexColor = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type MiniMaxConfig struct {
	BaseURL string
	APIKey  string
	Model   string
	Client  *http.Client
}

type MiniMaxProvider struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

func NewMiniMaxProvider(cfg MiniMaxConfig) (*MiniMaxProvider, error) {
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if _, err := url.ParseRequestURI(base); err != nil || base == "" {
		return nil, fmt.Errorf("invalid MiniMax base URL")
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("MiniMax API key is required")
	}
	if strings.TrimSpace(cfg.Model) == "" {
		return nil, fmt.Errorf("MiniMax model is required")
	}
	client := cfg.Client
	if client == nil {
		client = &http.Client{Timeout: 75 * time.Second}
	}
	return &MiniMaxProvider{
		baseURL: base,
		apiKey:  strings.TrimSpace(cfg.APIKey),
		model:   strings.TrimSpace(cfg.Model),
		client:  client,
	}, nil
}

func (p *MiniMaxProvider) ProviderStatus() ProviderStatus {
	host := ""
	if parsed, err := url.Parse(p.baseURL); err == nil {
		host = parsed.Hostname()
	}
	return ProviderStatus{State: "connected", Provider: "MiniMax", Model: p.model, BaseURLHost: host, PromptVersion: namingPromptVersion}
}

func (p *MiniMaxProvider) AnalyzeIntent(ctx context.Context, req FounderIntentRequest) (FounderIntent, error) {
	return (DeterministicProvider{}).AnalyzeIntent(ctx, req)
}

func (p *MiniMaxProvider) GenerateNames(ctx context.Context, intent FounderIntent) ([]NamingProposal, error) {
	if err := intent.Validate(); err != nil {
		return nil, err
	}
	intentJSON, _ := json.Marshal(intent)
	prompt := `你是企业品牌战略顾问。根据创始人输入生成4个明显不同的公司身份候选。
只输出一个JSON对象，不要Markdown、解释或思考过程。格式：
{"proposals":[{"chinese_name":"","english_name":"","short_name":"","rationale":"","slogan":"","keywords":[""],"primary_color":"#RRGGBB"}]}
要求中文全称符合“字号+行业/业务特征+有限公司”的候选表达；名称不是现实工商核名结果；避免复用示例名。
prompt_version=` + namingPromptVersion + `
founder_intent=` + string(intentJSON)
	messages := []map[string]string{
		{"role": "system", "content": "Return strict JSON only. Never include chain-of-thought."},
		{"role": "user", "content": prompt},
	}
	content, err := p.complete(ctx, messages, 0.8, 8192)
	if err != nil {
		return nil, err
	}
	result, decodeErr := decodeNamingResult(content)
	if decodeErr != nil {
		repairMessages := []map[string]string{
			{"role": "system", "content": "Return one complete strict JSON object only. Do not include markdown or reasoning."},
			{"role": "user", "content": "重新生成完整的公司身份候选 JSON。上一次输出不完整或非法。必须严格符合原格式并返回4个候选。\n原始创业意图：" + string(intentJSON)},
		}
		repaired, repairErr := p.complete(ctx, repairMessages, 0.3, 8192)
		if repairErr != nil {
			return nil, fmt.Errorf("MiniMax naming repair failed after invalid response: %w", repairErr)
		}
		result, decodeErr = decodeNamingResult(repaired)
		if decodeErr != nil {
			return nil, fmt.Errorf("MiniMax returned invalid naming JSON after one repair: %w", decodeErr)
		}
	}
	if len(result.Proposals) < 4 || len(result.Proposals) > 6 {
		return nil, fmt.Errorf("MiniMax must return 4 to 6 naming proposals")
	}
	seen := map[string]bool{}
	for i := range result.Proposals {
		item := &result.Proposals[i]
		item.ChineseName = strings.TrimSpace(item.ChineseName)
		item.EnglishName = strings.TrimSpace(item.EnglishName)
		item.ShortName = strings.TrimSpace(item.ShortName)
		item.Rationale = strings.TrimSpace(item.Rationale)
		item.Slogan = strings.TrimSpace(item.Slogan)
		if item.ChineseName == "" || item.EnglishName == "" || item.ShortName == "" || item.Rationale == "" || item.Slogan == "" || len(item.Keywords) == 0 || !hexColor.MatchString(item.Primary) {
			return nil, fmt.Errorf("MiniMax proposal %d is incomplete", i+1)
		}
		if seen[item.ChineseName] {
			return nil, fmt.Errorf("MiniMax returned duplicate company names")
		}
		seen[item.ChineseName] = true
		item.ProposalID = fmt.Sprintf("name-%s-ai-%02d", stableSlug(intent.CaseCode), i+1)
		item.RiskHints = []string{"AI 生成候选；尚未完成现实工商核名与商标检索"}
		item.Status = "candidate"
	}
	return result.Proposals, nil
}

// CompleteJSON exposes the same configured provider to other candidate-only
// AESE planners without leaking credentials or accepting arbitrary endpoints.
func (p *MiniMaxProvider) CompleteJSON(ctx context.Context, system, user string, temperature float64, maxTokens int) (string, string, map[string]int, error) {
	evidence := GenerationEvidence{TokenUsage: map[string]int{}}
	generationCtx := WithGenerationEvidence(ctx, func(value GenerationEvidence) { evidence = value })
	content, err := p.complete(generationCtx, []map[string]string{{"role": "system", "content": system}, {"role": "user", "content": user}}, temperature, maxTokens)
	return content, evidence.RequestID, evidence.TokenUsage, err
}

func (p *MiniMaxProvider) complete(ctx context.Context, messages []map[string]string, temperature float64, maxTokens int) (string, error) {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		content, transient, err := p.completeOnce(ctx, messages, temperature, maxTokens)
		if err == nil {
			return content, nil
		}
		lastErr = err
		if !transient || attempt == 1 {
			break
		}
		timer := time.NewTimer(180 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return "", ctx.Err()
		case <-timer.C:
		}
	}
	return "", lastErr
}

func (p *MiniMaxProvider) completeOnce(ctx context.Context, messages []map[string]string, temperature float64, maxTokens int) (string, bool, error) {
	requestBody := map[string]any{
		"model": p.model, "messages": messages,
		"temperature": temperature, "max_tokens": maxTokens,
	}
	encoded, _ := json.Marshal(requestBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(encoded))
	if err != nil {
		return "", false, err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return "", true, fmt.Errorf("MiniMax request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", true, fmt.Errorf("MiniMax response read failed: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		transient := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
		return "", transient, fmt.Errorf("MiniMax returned HTTP %d", resp.StatusCode)
	}
	var completion struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(body, &completion); err != nil || len(completion.Choices) == 0 {
		return "", false, fmt.Errorf("MiniMax returned an invalid completion envelope")
	}
	content := strings.TrimSpace(completion.Choices[0].Message.Content)
	finishReason := strings.ToLower(strings.TrimSpace(completion.Choices[0].FinishReason))
	recordGenerationEvidence(ctx, GenerationEvidence{
		RequestID:    firstNonEmpty(resp.Header.Get("x-request-id"), resp.Header.Get("request-id")),
		FinishReason: completion.Choices[0].FinishReason,
		TokenUsage: map[string]int{
			"prompt_tokens": completion.Usage.PromptTokens, "completion_tokens": completion.Usage.CompletionTokens,
			"total_tokens": completion.Usage.TotalTokens,
		},
	})
	if finishReason == "length" || finishReason == "max_tokens" {
		return "", false, fmt.Errorf("MiniMax completion truncated (finish_reason=%s)", finishReason)
	}
	if content == "" {
		return "", false, fmt.Errorf("MiniMax returned empty completion content")
	}
	if strings.HasPrefix(content, "<think>") {
		end := strings.Index(content, "</think>")
		if end < 0 {
			return "", false, fmt.Errorf("MiniMax returned unterminated reasoning")
		}
		content = strings.TrimSpace(content[end+len("</think>"):])
	}
	if strings.Contains(content, "<think>") || strings.Contains(content, "```") {
		return "", false, fmt.Errorf("MiniMax returned non-JSON reasoning or markdown")
	}
	return content, false, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

type namingResult struct {
	Proposals []NamingProposal `json:"proposals"`
}

func decodeNamingResult(content string) (namingResult, error) {
	var result namingResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return namingResult{}, fmt.Errorf("invalid naming JSON: %w", err)
	}
	return result, nil
}
