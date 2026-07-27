package creative

import (
	"context"
	"testing"
	"time"
)

func validRequest() FounderIntentRequest {
	return FounderIntentRequest{
		TenantID: "tenant-hctm-genesis", CaseCode: "INC-GAME-001",
		RawIdea:  "创建一家面向新能源汽车客户的高可靠热管理零部件企业",
		Industry: "汽车热管理", Customers: []string{"新能源汽车整车厂"},
		Offerings: []string{"电池冷却板"}, BrandTraits: []string{"可靠", "精密", "长期主义"},
		CapitalMinor: "3000000000", RiskAppetite: "balanced",
	}
}

func TestDeterministicProviderProducesTraceableCandidates(t *testing.T) {
	provider := DeterministicProvider{Now: func() time.Time { return time.Date(2026, 7, 27, 0, 0, 0, 0, time.UTC) }}
	intent, err := provider.AnalyzeIntent(context.Background(), validRequest())
	if err != nil {
		t.Fatal(err)
	}
	names, err := provider.GenerateNames(context.Background(), intent)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 4 || names[0].Status != "candidate" || len(names[0].RiskHints) == 0 {
		t.Fatalf("invalid candidates: %#v", names)
	}
}

func TestIntentKeepsMissingInputsVisible(t *testing.T) {
	req := validRequest()
	req.Customers = nil
	req.Offerings = nil
	intent, err := (DeterministicProvider{}).AnalyzeIntent(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if len(intent.NeedsConfirm) != 2 {
		t.Fatalf("missing confirmation gaps: %#v", intent)
	}
}
