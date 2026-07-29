package creative

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Provider interface {
	AnalyzeIntent(context.Context, FounderIntentRequest) (FounderIntent, error)
	GenerateNames(context.Context, FounderIntent) ([]NamingProposal, error)
}

type ProviderStatus struct {
	State         string `json:"state"`
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	BaseURLHost   string `json:"base_url_host,omitempty"`
	PromptVersion string `json:"prompt_version"`
}

type StatusProvider interface{ ProviderStatus() ProviderStatus }

type GenerationEvidence struct {
	RequestID    string
	FinishReason string
	TokenUsage   map[string]int
}

type generationEvidenceKey struct{}

func WithGenerationEvidence(ctx context.Context, sink func(GenerationEvidence)) context.Context {
	return context.WithValue(ctx, generationEvidenceKey{}, sink)
}

func recordGenerationEvidence(ctx context.Context, evidence GenerationEvidence) {
	if sink, ok := ctx.Value(generationEvidenceKey{}).(func(GenerationEvidence)); ok && sink != nil {
		sink(evidence)
	}
}

func (DeterministicProvider) ProviderStatus() ProviderStatus {
	return ProviderStatus{State: "fallback", Provider: "deterministic", Model: "built-in", PromptVersion: namingPromptVersion}
}

type FounderIntentRequest struct {
	TenantID     string   `json:"tenant_id"`
	CaseCode     string   `json:"case_code"`
	RawIdea      string   `json:"raw_idea"`
	Industry     string   `json:"industry"`
	Customers    []string `json:"customers"`
	Offerings    []string `json:"offerings"`
	BrandTraits  []string `json:"brand_traits"`
	CapitalMinor string   `json:"capital_minor"`
	RiskAppetite string   `json:"risk_appetite"`
}

// DeterministicProvider is the offline and test-safe baseline. Its outputs
// are explicit candidates and never masquerade as an external model call.
type DeterministicProvider struct{ Now func() time.Time }

func (p DeterministicProvider) AnalyzeIntent(_ context.Context, req FounderIntentRequest) (FounderIntent, error) {
	now := time.Now().UTC()
	if p.Now != nil {
		now = p.Now().UTC()
	}
	intent := FounderIntent{
		SchemaVersion: SchemaVersion, IntentID: "intent-" + stableSlug(req.CaseCode),
		TenantID: req.TenantID, CaseCode: req.CaseCode, RawIdea: strings.TrimSpace(req.RawIdea),
		Industry: req.Industry, Customers: req.Customers, Offerings: req.Offerings,
		BrandTraits: req.BrandTraits, CapitalMinor: req.CapitalMinor,
		RiskAppetite: req.RiskAppetite, CreatedAt: now.Format(time.RFC3339),
	}
	if len(intent.Customers) == 0 {
		intent.NeedsConfirm = append(intent.NeedsConfirm, "目标客户尚未确认")
	}
	if len(intent.Offerings) == 0 {
		intent.NeedsConfirm = append(intent.NeedsConfirm, "核心产品或服务尚未确认")
	}
	if err := intent.Validate(); err != nil {
		return FounderIntent{}, err
	}
	return intent, nil
}

func (DeterministicProvider) GenerateNames(ctx context.Context, intent FounderIntent) ([]NamingProposal, error) {
	if err := intent.Validate(); err != nil {
		return nil, err
	}
	industry := strings.TrimSuffix(strings.TrimSpace(intent.Industry), "行业")
	if industry == "" {
		industry = "工业科技"
	}
	base := []struct{ zh, en, short, promise, color string }{
		{"澄流" + industry + "有限公司", "ClearFlow Industrial Technologies", "澄流科技", "让复杂工业系统稳定而高效地流动", "#167C80"},
		{"衡创" + industry + "有限公司", "Equilibrium Innovation Group", "衡创集团", "用工程平衡连接创新与可靠交付", "#315C8C"},
		{"启域" + industry + "有限公司", "Genesis Domain Systems", "启域系统", "从零建立可持续进化的产业能力", "#9A5B2F"},
		{"铸远" + industry + "有限公司", "Forge Horizon Industries", "铸远工业", "以制造根基创造长期价值", "#50677A"},
	}
	out := make([]NamingProposal, 0, len(base))
	for i, candidate := range base {
		out = append(out, NamingProposal{
			ProposalID:  fmt.Sprintf("name-%s-%02d", stableSlug(intent.CaseCode), i+1),
			ChineseName: candidate.zh, EnglishName: candidate.en, ShortName: candidate.short,
			Rationale: candidate.promise, Slogan: candidate.promise,
			Keywords: append([]string{}, intent.BrandTraits...), Primary: candidate.color,
			RiskHints: []string{"虚构世界名称可用；尚未完成现实工商核名与商标检索"},
			Status:    "candidate",
		})
	}
	recordGenerationEvidence(ctx, GenerationEvidence{
		RequestID: "deterministic-" + stableSlug(intent.CaseCode), FinishReason: "deterministic",
		TokenUsage: map[string]int{"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
	})
	return out, nil
}

func stableSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer(" ", "-", "_", "-", "/", "-").Replace(value)
	if value == "" {
		return "unnamed"
	}
	return value
}
