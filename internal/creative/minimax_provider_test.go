package creative

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMiniMaxDefaultTimeoutPreservesRetryWindow(t *testing.T) {
	provider, err := NewMiniMaxProvider(MiniMaxConfig{BaseURL: "https://api.minimax.chat/v1", APIKey: "secret", Model: "MiniMax-M3"})
	if err != nil {
		t.Fatal(err)
	}
	if provider.client.Timeout > 40*time.Second {
		t.Fatalf("single attempt timeout %s leaves no bounded retry window", provider.client.Timeout)
	}
}

func TestMiniMaxCompleteJSONDisablesM3Thinking(t *testing.T) {
	var thinkingType string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Thinking struct {
				Type string `json:"type"`
			} `json:"thinking"`
		}
		if err := decodeJSON(r, &request); err != nil {
			t.Fatal(err)
		}
		thinkingType = request.Thinking.Type
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"ok\":true}"},"finish_reason":"stop"}]}`))
	}))
	defer upstream.Close()
	provider, err := NewMiniMaxProvider(MiniMaxConfig{BaseURL: upstream.URL, APIKey: "secret", Model: "MiniMax-M3", Client: upstream.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := provider.CompleteJSON(context.Background(), "system", "user", 0.25, 4096); err != nil {
		t.Fatal(err)
	}
	if thinkingType != "disabled" {
		t.Fatalf("thinking.type=%q, want disabled for strict structured generation", thinkingType)
	}
}

func TestMiniMaxProviderGeneratesModelCandidates(t *testing.T) {
	var requestedModel string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("unexpected authorization header %q", got)
		}
		var request struct {
			Model string `json:"model"`
		}
		if err := decodeJSON(r, &request); err != nil {
			t.Fatal(err)
		}
		requestedModel = request.Model
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"<think>internal reasoning must not be parsed as business output</think>\n{\"proposals\":[{\"chinese_name\":\"澜衡热控科技有限公司\",\"english_name\":\"LanHeng Thermal Control\",\"short_name\":\"澜衡热控\",\"rationale\":\"体现流体热管理与可靠平衡\",\"slogan\":\"让每一度热量都有秩序\",\"keywords\":[\"可靠\",\"热管理\"],\"primary_color\":\"#176B68\"},{\"chinese_name\":\"矩川热管理有限公司\",\"english_name\":\"MatrixRiver Thermal\",\"short_name\":\"矩川热管\",\"rationale\":\"突出工程矩阵与流动能力\",\"slogan\":\"精密驱动可靠交付\",\"keywords\":[\"精密\",\"工程\"],\"primary_color\":\"#315C8C\"},{\"chinese_name\":\"熵澈工业科技有限公司\",\"english_name\":\"EntropyClear Industries\",\"short_name\":\"熵澈工业\",\"rationale\":\"表达复杂热系统的清晰治理\",\"slogan\":\"看清热流，掌控未来\",\"keywords\":[\"清晰\",\"长期主义\"],\"primary_color\":\"#8A5635\"},{\"chinese_name\":\"曜冷系统科技有限公司\",\"english_name\":\"Lumicool Systems\",\"short_name\":\"曜冷系统\",\"rationale\":\"强调新能源冷却系统能力\",\"slogan\":\"为新能源守住温度边界\",\"keywords\":[\"新能源\",\"可靠\"],\"primary_color\":\"#3C6478\"}]}"}}]}`))
	}))
	defer upstream.Close()

	provider, err := NewMiniMaxProvider(MiniMaxConfig{
		BaseURL: upstream.URL + "/v1",
		APIKey:  "secret",
		Model:   "MiniMax-M3",
		Client:  upstream.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	intent, err := (DeterministicProvider{}).AnalyzeIntent(context.Background(), validRequest())
	if err != nil {
		t.Fatal(err)
	}
	proposals, err := provider.GenerateNames(context.Background(), intent)
	if err != nil {
		t.Fatal(err)
	}
	if requestedModel != "MiniMax-M3" {
		t.Fatalf("model = %q", requestedModel)
	}
	if len(proposals) != 4 || proposals[0].ChineseName != "澜衡热控科技有限公司" {
		t.Fatalf("unexpected proposals: %#v", proposals)
	}
	if proposals[0].ProposalID == "" || proposals[0].Status != "candidate" {
		t.Fatalf("missing governed candidate fields: %#v", proposals[0])
	}
}

func TestMiniMaxProviderDoesNotSilentlyFallback(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "quota exceeded", http.StatusTooManyRequests)
	}))
	defer upstream.Close()
	provider, err := NewMiniMaxProvider(MiniMaxConfig{
		BaseURL: upstream.URL,
		APIKey:  "secret",
		Model:   "MiniMax-M3",
		Client:  upstream.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	intent, _ := (DeterministicProvider{}).AnalyzeIntent(context.Background(), validRequest())
	_, err = provider.GenerateNames(context.Background(), intent)
	if err == nil || !strings.Contains(err.Error(), "429") {
		t.Fatalf("expected explicit provider failure, got %v", err)
	}
}

func TestMiniMaxProviderRetriesTransientBadGateway(t *testing.T) {
	calls := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls == 1 {
			http.Error(w, "temporary upstream failure", http.StatusBadGateway)
			return
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"proposals\":[{\"chinese_name\":\"澜衡热控科技有限公司\",\"english_name\":\"LanHeng Thermal Control\",\"short_name\":\"澜衡热控\",\"rationale\":\"体现可靠热管理\",\"slogan\":\"让热量有秩序\",\"keywords\":[\"可靠\"],\"primary_color\":\"#176B68\"},{\"chinese_name\":\"矩川热管理有限公司\",\"english_name\":\"MatrixRiver Thermal\",\"short_name\":\"矩川热管\",\"rationale\":\"突出工程能力\",\"slogan\":\"精密驱动交付\",\"keywords\":[\"工程\"],\"primary_color\":\"#315C8C\"},{\"chinese_name\":\"熵澈工业科技有限公司\",\"english_name\":\"EntropyClear Industries\",\"short_name\":\"熵澈工业\",\"rationale\":\"治理复杂热系统\",\"slogan\":\"看清热流\",\"keywords\":[\"清晰\"],\"primary_color\":\"#8A5635\"},{\"chinese_name\":\"曜冷系统科技有限公司\",\"english_name\":\"Lumicool Systems\",\"short_name\":\"曜冷系统\",\"rationale\":\"新能源冷却能力\",\"slogan\":\"守住温度边界\",\"keywords\":[\"新能源\"],\"primary_color\":\"#3C6478\"}]}"}}]}`))
	}))
	defer upstream.Close()
	provider, err := NewMiniMaxProvider(MiniMaxConfig{
		BaseURL: upstream.URL, APIKey: "secret", Model: "MiniMax-M3", Client: upstream.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	intent, _ := (DeterministicProvider{}).AnalyzeIntent(context.Background(), validRequest())
	proposals, err := provider.GenerateNames(context.Background(), intent)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 || len(proposals) != 4 {
		t.Fatalf("calls=%d proposals=%d", calls, len(proposals))
	}
}

func TestMiniMaxCompleteJSONRejectsLengthTruncation(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"options\":[{\"title\":\"被截断\""},"finish_reason":"length"}],"usage":{"completion_tokens":8192,"total_tokens":8200}}`))
	}))
	defer upstream.Close()
	provider, err := NewMiniMaxProvider(MiniMaxConfig{
		BaseURL: upstream.URL, APIKey: "secret", Model: "MiniMax-M3", Client: upstream.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, err = provider.CompleteJSON(context.Background(), "system", "user", 0.5, 8192)
	if err == nil || !strings.Contains(err.Error(), "truncated") {
		t.Fatalf("expected explicit truncation error, got %v", err)
	}
}

func TestMiniMaxProviderRepairsTruncatedJSONOnce(t *testing.T) {
	calls := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls == 1 {
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]any{"content": `{"proposals":[{"chinese_name":"截断`}}}})
			return
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"proposals\":[{\"chinese_name\":\"澜衡热控科技有限公司\",\"english_name\":\"LanHeng Thermal Control\",\"short_name\":\"澜衡热控\",\"rationale\":\"体现可靠热管理\",\"slogan\":\"让热量有秩序\",\"keywords\":[\"可靠\"],\"primary_color\":\"#176B68\"},{\"chinese_name\":\"矩川热管理有限公司\",\"english_name\":\"MatrixRiver Thermal\",\"short_name\":\"矩川热管\",\"rationale\":\"突出工程能力\",\"slogan\":\"精密驱动交付\",\"keywords\":[\"工程\"],\"primary_color\":\"#315C8C\"},{\"chinese_name\":\"熵澈工业科技有限公司\",\"english_name\":\"EntropyClear Industries\",\"short_name\":\"熵澈工业\",\"rationale\":\"治理复杂热系统\",\"slogan\":\"看清热流\",\"keywords\":[\"清晰\"],\"primary_color\":\"#8A5635\"},{\"chinese_name\":\"曜冷系统科技有限公司\",\"english_name\":\"Lumicool Systems\",\"short_name\":\"曜冷系统\",\"rationale\":\"新能源冷却能力\",\"slogan\":\"守住温度边界\",\"keywords\":[\"新能源\"],\"primary_color\":\"#3C6478\"}]}"}}]}`))
	}))
	defer upstream.Close()
	provider, err := NewMiniMaxProvider(MiniMaxConfig{
		BaseURL: upstream.URL, APIKey: "secret", Model: "MiniMax-M3", Client: upstream.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	intent, _ := (DeterministicProvider{}).AnalyzeIntent(context.Background(), validRequest())
	proposals, err := provider.GenerateNames(context.Background(), intent)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 || len(proposals) != 4 {
		t.Fatalf("calls=%d proposals=%d", calls, len(proposals))
	}
}

func decodeJSON(r *http.Request, target any) error {
	return json.NewDecoder(r.Body).Decode(target)
}
