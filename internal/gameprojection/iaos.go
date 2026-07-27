package gameprojection

import (
	"fmt"
	"strings"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
)

var stateChapter = map[string]string{
	"draft": "founder_intent", "incorporation_opened": "founder_office", "incorporation_case_opened": "founder_office",
	"founder_resolution_approved": "founder_office", "capital_commitments_confirmed": "registration",
	"registration_submitted": "registration", "legal_entity_registered": "banking",
	"bank_account_opening_submitted": "banking", "bank_account_opened": "banking",
	"capital_contribution_verified": "banking", "organization_established": "talent_governance",
	"executive_appointment_proposed": "talent_governance", "executive_appointment_accepted": "talent_governance",
	"executive_appointment_approved": "talent_governance", "executive_appointments_accepted": "talent_governance",
	"operating_mandate_active": "opening", "operating_mandates_activated": "opening",
	"initial_budget_approved": "opening", "enterprise_operational_ready": "operating_world",
}

// FromIAOS builds the game view exclusively from IAOS committed trace and
// persistent work items. It never executes a capability or advances a process.
func FromIAOS(trace iaosclient.IncorporationTrace, items []iaosclient.IncorporationWorkItem) (Projection, error) {
	if trace.SchemaVersion != "1.0" || !trace.Verified || trace.CaseCode == "" || trace.State.TenantID == "" {
		return Projection{}, fmt.Errorf("IAOS incorporation trace is incomplete or unverified")
	}
	chapter, ok := stateChapter[trace.State.State]
	if !ok {
		return Projection{}, fmt.Errorf("IAOS state %q has no game chapter", trace.State.State)
	}
	simTime := ""
	correlation := ""
	for _, entry := range trace.Journal {
		if _, err := time.Parse(time.RFC3339, entry.CreatedAt); err == nil {
			simTime = entry.CreatedAt
		}
		if entry.CorrelationID != "" {
			correlation = entry.CorrelationID
		}
	}
	if simTime == "" {
		return Projection{}, fmt.Errorf("IAOS trace has no RFC3339 journal timestamp")
	}
	cursor := int64(len(trace.Journal))
	exchanges := make([]Exchange, 0, len(trace.WorldExchanges))
	for _, item := range trace.WorldExchanges {
		if item.Cursor > cursor {
			cursor = item.Cursor
		}
		exchanges = append(exchanges, Exchange{
			ExchangeID: item.MessageID, Kind: item.Kind, Status: "committed",
			Correlation: item.CorrelationID, EvidenceRef: "world:" + item.MessageID, OccurredAt: item.RecordedAt,
		})
	}
	completed := 0
	workItems := make([]WorkItem, 0, len(items))
	actors := map[string]Actor{
		"founder-principal":      {ActorID: "founder-principal", ActorType: "human", DisplayName: "创始治理者", Position: "董事长", State: "available"},
		"incorporation-agent":    {ActorID: "incorporation-agent", ActorType: "agent", DisplayName: agentLabel("incorporation-agent"), Position: agentLabel("incorporation-agent"), State: "available"},
		"governance-agent":       {ActorID: "governance-agent", ActorType: "agent", DisplayName: agentLabel("governance-agent"), Position: agentLabel("governance-agent"), State: "available"},
		"legal-compliance-agent": {ActorID: "legal-compliance-agent", ActorType: "agent", DisplayName: agentLabel("legal-compliance-agent"), Position: agentLabel("legal-compliance-agent"), State: "available"},
		"finance-agent":          {ActorID: "finance-agent", ActorType: "agent", DisplayName: agentLabel("finance-agent"), Position: agentLabel("finance-agent"), State: "available"},
		"audit-agent":            {ActorID: "audit-agent", ActorType: "agent", DisplayName: agentLabel("audit-agent"), Position: agentLabel("audit-agent"), State: "available"},
	}
	for _, item := range items {
		status := item.Effective
		if status == "" {
			status = item.Status
		}
		if status == "completed" {
			completed++
		}
		ownerType := ownerType(item.Participant)
		workItems = append(workItems, WorkItem{
			WorkItemID: fmt.Sprintf("WI-%02d", item.Sequence), Title: item.Capability,
			Kind: item.TaskType, Status: status, OwnerType: ownerType, OwnerID: item.Participant,
			Capability: item.Capability, RequiresMe: ownerType == "human" && status != "completed" && status != "locked",
			EvidenceRef: fmt.Sprintf("iaos:work-item:%s:%d", trace.CaseCode, item.Sequence),
		})
		if ownerType == "agent" {
			state := "available"
			if status == "ready" || status == "running" {
				state = "working"
			}
			if configured, ok := item.AgentAuth["status"].(string); ok && configured != "" && configured != "active" {
				state = configured
			}
			actors[item.Participant] = Actor{ActorID: item.Participant, ActorType: "agent", DisplayName: agentLabel(item.Participant), Position: agentLabel(item.Participant), State: state, WorkItemID: fmt.Sprintf("WI-%02d", item.Sequence)}
		}
	}
	actorList := make([]Actor, 0, len(actors))
	for _, id := range []string{"founder-principal", "incorporation-agent", "governance-agent", "legal-compliance-agent", "finance-agent", "audit-agent"} {
		if actor, exists := actors[id]; exists {
			actorList = append(actorList, actor)
		}
	}
	progress := 0
	if len(items) > 0 {
		progress = completed * 100 / len(items)
	}
	currency := trace.State.Currency
	if currency == "" {
		currency = "CNY"
	}
	projection := Projection{
		SchemaVersion: SchemaVersion, ProjectionID: fmt.Sprintf("gp-%s-%d", strings.ToLower(trace.CaseCode), cursor),
		TenantID: trace.State.TenantID, CaseCode: trace.CaseCode, WorldRunID: correlation,
		Chapter: chapter, SimTime: simTime, TimeScale: 0, Paused: true, Cursor: cursor,
		Scene:     Scene{SceneID: sceneForChapter(chapter), Mode: "2.5d", Theme: "industrial-warm"},
		Lifecycle: Lifecycle{State: trace.State.State, CurrentStep: currentStep(workItems), Progress: progress},
		Buildings: buildingsForChapter(chapter), Actors: actorList, WorkItems: workItems,
		Resources: Resources{
			FounderCash:      Money{Value: "0", Currency: currency, Scale: 2},
			CompanyCash:      Money{Value: fmt.Sprint(trace.State.ContributionMinor), Currency: currency, Scale: 2},
			CapitalCommitted: Money{Value: fmt.Sprint(trace.State.CommitmentMinor), Currency: currency, Scale: 2},
			CapitalPaid:      Money{Value: fmt.Sprint(trace.State.ContributionMinor), Currency: currency, Scale: 2},
			BudgetAuthorized: Money{Value: fmt.Sprint(trace.State.BudgetMinor), Currency: currency, Scale: 2}, RiskLevel: "low",
		},
		Exchanges: exchanges, Brand: Brand{Status: "selected", CompanyName: trace.State.ProposedName, PrimaryColor: "#167C80"},
		Notifications: []Notification{{NotificationID: fmt.Sprintf("notice-%d", cursor), Severity: "info", Message: currentStep(workItems)}},
		EvidenceRefs:  []EvidenceRef{{Ref: "iaos:trace:" + trace.CaseCode, Kind: "incorporation_trace"}, {Ref: "iaos:work-items:" + trace.CaseCode, Kind: "process_work_items"}},
	}
	if trace.State.State == "enterprise_operational_ready" {
		projection.Brand.Status = "approved"
	}
	for index, discrepancy := range trace.State.Discrepancies {
		projection.Notifications = append(projection.Notifications, Notification{
			NotificationID: fmt.Sprintf("discrepancy-%d", index+1), Severity: "warning",
			Message: discrepancy, ActionRef: "iaos:discrepancy:" + trace.CaseCode,
		})
	}
	if projection.WorldRunID == "" {
		projection.WorldRunID = "case-" + trace.CaseCode
	}
	return projection, projection.Validate()
}

func agentLabel(id string) string {
	return map[string]string{"incorporation-agent": "企业设立专员", "governance-agent": "治理组织专员", "legal-compliance-agent": "法务合规专员", "finance-agent": "财务负责人", "audit-agent": "独立审计专员"}[id]
}

func currentStep(items []WorkItem) string {
	for _, item := range items {
		if item.Status != "completed" && item.Status != "locked" {
			return item.Capability
		}
	}
	return "enterprise_operational_ready"
}
