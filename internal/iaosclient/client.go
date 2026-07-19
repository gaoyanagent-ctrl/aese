package iaosclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const maxResponseBytes = 4 << 20

// Doer is implemented by *http.Client and keeps the adapter easy to test.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Config struct {
	BaseURL  string
	Token    string
	TenantID string
	HTTP     Doer
}

// Client is a narrow adapter over IAOS' authenticated metadata and dynamic
// entity APIs. It deliberately has no database or NATS escape hatch.
type Client struct {
	baseURL  *url.URL
	token    string
	tenantID string
	http     Doer
}

func New(cfg Config) (*Client, error) {
	base, err := url.Parse(strings.TrimSpace(cfg.BaseURL))
	if err != nil || base.Scheme == "" || base.Host == "" {
		return nil, fmt.Errorf("invalid IAOS base URL %q", cfg.BaseURL)
	}
	if base.Scheme != "http" && base.Scheme != "https" {
		return nil, fmt.Errorf("IAOS base URL must use http or https")
	}
	if base.User != nil {
		return nil, fmt.Errorf("IAOS base URL must not contain user information")
	}
	if base.RawQuery != "" || base.Fragment != "" {
		return nil, fmt.Errorf("IAOS base URL must not contain query or fragment")
	}
	if strings.TrimSpace(cfg.Token) == "" {
		return nil, fmt.Errorf("IAOS bearer token is required")
	}
	doer := cfg.HTTP
	if doer == nil {
		doer = &http.Client{Timeout: 20 * time.Second}
	}
	base.Path = strings.TrimRight(base.Path, "/") + "/"
	return &Client{baseURL: base, token: strings.TrimSpace(cfg.Token), tenantID: strings.TrimSpace(cfg.TenantID), http: doer}, nil
}

type Field struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type Schema struct {
	Entity          string  `json:"entity"`
	Version         string  `json:"version"`
	PhysicalTable   string  `json:"physical_table_name"`
	StorageStrategy string  `json:"storage_strategy"`
	Fields          []Field `json:"fields"`
	Permissions     struct {
		Create bool `json:"create"`
		Update bool `json:"update"`
		Delete bool `json:"delete"`
	} `json:"permissions"`
}

func (c *Client) Schema(ctx context.Context, entity string) (Schema, error) {
	var out Schema
	if err := validateEntity(entity); err != nil {
		return out, err
	}
	err := c.request(ctx, http.MethodGet, "api/v1/metadata/schema/"+url.PathEscape(entity), nil, &out)
	return out, err
}

type UpsertAction string

const (
	ActionCreate    UpsertAction = "create"
	ActionUpdate    UpsertAction = "update"
	ActionUnchanged UpsertAction = "unchanged"
)

type UpsertRequest struct {
	Entity     string
	NaturalKey []string
	Record     map[string]any
}

type UpsertPlan struct {
	Action     UpsertAction   `json:"action"`
	Entity     string         `json:"entity"`
	NaturalKey map[string]any `json:"natural_key"`
	RecordID   string         `json:"record_id,omitempty"`
	Changed    []string       `json:"changed_fields,omitempty"`
	Existing   map[string]any `json:"-"`
}

type UpsertResult struct {
	UpsertPlan
	Applied bool `json:"applied"`
}

// PlanUpsert performs only authenticated reads. It is the primitive used by
// dry-run impact reporting and never mutates IAOS.
func (c *Client) PlanUpsert(ctx context.Context, req UpsertRequest) (UpsertPlan, error) {
	if err := validateUpsert(req); err != nil {
		return UpsertPlan{}, err
	}
	key := make(map[string]any, len(req.NaturalKey))
	for _, name := range req.NaturalKey {
		key[name] = req.Record[name]
	}
	records, err := c.findExact(ctx, req.Entity, key)
	if err != nil {
		return UpsertPlan{}, err
	}
	if len(records) > 1 {
		return UpsertPlan{}, fmt.Errorf("%s natural key %v matched %d IAOS records", req.Entity, key, len(records))
	}
	plan := UpsertPlan{Action: ActionCreate, Entity: req.Entity, NaturalKey: key}
	if len(records) == 0 {
		return plan, nil
	}
	plan.Existing = records[0]
	plan.RecordID, _ = records[0]["id"].(string)
	if plan.RecordID == "" {
		return UpsertPlan{}, fmt.Errorf("%s natural key %v matched record without id", req.Entity, key)
	}
	for name, wanted := range req.Record {
		if !equivalentJSON(records[0][name], wanted) {
			plan.Changed = append(plan.Changed, name)
		}
	}
	if len(plan.Changed) == 0 {
		plan.Action = ActionUnchanged
	} else {
		plan.Action = ActionUpdate
	}
	return plan, nil
}

// Upsert applies a previously computable idempotent operation. Callers must
// make their own explicit apply decision; there is intentionally no force flag.
func (c *Client) Upsert(ctx context.Context, req UpsertRequest) (UpsertResult, error) {
	plan, err := c.PlanUpsert(ctx, req)
	if err != nil {
		return UpsertResult{}, err
	}
	result := UpsertResult{UpsertPlan: plan}
	switch plan.Action {
	case ActionUnchanged:
		return result, nil
	case ActionCreate:
		var response struct {
			ID string `json:"id"`
		}
		if err := c.request(ctx, http.MethodPost, "api/v1/entities/"+url.PathEscape(req.Entity), req.Record, &response); err != nil {
			return UpsertResult{}, err
		}
		result.RecordID = response.ID
		result.Applied = true
		return result, nil
	case ActionUpdate:
		path := "api/v1/entities/" + url.PathEscape(req.Entity) + "/" + url.PathEscape(plan.RecordID)
		if err := c.request(ctx, http.MethodPut, path, req.Record, nil); err != nil {
			return UpsertResult{}, err
		}
		result.Applied = true
		return result, nil
	default:
		return UpsertResult{}, fmt.Errorf("unsupported upsert action %q", plan.Action)
	}
}

type DecomposeResult struct {
	Status       string `json:"status"`
	Decomposing  bool   `json:"decomposing"`
	SalesOrderID string `json:"sales_order_id"`
	OrderNo      string `json:"order_no"`
}

type SimulationBusinessObject struct {
	Type string `json:"type"`
	Code string `json:"code"`
}

type SimulationEventRequest struct {
	EventType      string                   `json:"event_type"`
	Source         string                   `json:"source"`
	OccurredAt     string                   `json:"occurred_at"`
	CorrelationID  string                   `json:"correlation_id"`
	CausationID    string                   `json:"causation_id,omitempty"`
	IdempotencyKey string                   `json:"idempotency_key"`
	BusinessObject SimulationBusinessObject `json:"business_object"`
	Payload        map[string]any           `json:"payload"`
}

type SimulationEventResult struct {
	EventID        string `json:"event_id"`
	Subject        string `json:"subject"`
	CorrelationID  string `json:"correlation_id"`
	Duplicate      bool   `json:"duplicate"`
	Committed      bool   `json:"committed"`
	BusinessObject struct {
		Type string `json:"type"`
		Code string `json:"code"`
		ID   string `json:"id"`
	} `json:"business_object"`
	Transition struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"transition"`
}

// IngestSimulationEvent invokes DES-048. Tenant and actor remain owned by the
// authenticated IAOS context; callers cannot provide a NATS subject.
func (c *Client) IngestSimulationEvent(ctx context.Context, request SimulationEventRequest) (SimulationEventResult, error) {
	var out SimulationEventResult
	if strings.TrimSpace(request.EventType) == "" || strings.TrimSpace(request.IdempotencyKey) == "" {
		return out, fmt.Errorf("simulation event_type and idempotency_key are required")
	}
	err := c.request(ctx, http.MethodPost, "api/v1/simulation/events", request, &out)
	return out, err
}

// DecomposeSalesOrder is the governed O2D ingress currently exposed by IAOS.
// IAOS confirms the order and records o2d.order.confirmed in its Outbox.
func (c *Client) DecomposeSalesOrder(ctx context.Context, id string) (DecomposeResult, error) {
	return c.DecomposeSalesOrderTrace(ctx, id, "", "")
}

func (c *Client) DecomposeSalesOrderTrace(ctx context.Context, id, correlationID, idempotencyKey string) (DecomposeResult, error) {
	var out DecomposeResult
	if strings.TrimSpace(id) == "" {
		return out, fmt.Errorf("sales order id is required")
	}
	err := c.requestHeaders(ctx, http.MethodPost, "api/v1/entities/sales_order/"+url.PathEscape(id)+"/decompose", map[string]any{}, &out, map[string]string{"X-Correlation-ID": correlationID, "Idempotency-Key": idempotencyKey})
	return out, err
}

type ScenarioObjectResult struct {
	Index         int            `json:"index"`
	Type          string         `json:"type"`
	Entity        string         `json:"entity"`
	Code          string         `json:"code"`
	Action        string         `json:"action"`
	ObjectID      string         `json:"object_id,omitempty"`
	NaturalKey    map[string]any `json:"natural_key,omitempty"`
	ChangedFields []string       `json:"changed_fields,omitempty"`
	Reason        string         `json:"reason,omitempty"`
}

type ScenarioSummary struct {
	PackKey       string                 `json:"pack_key"`
	PackVersion   string                 `json:"pack_version"`
	ScenarioKey   string                 `json:"scenario_key"`
	RunID         string                 `json:"run_id"`
	TenantID      string                 `json:"tenant_id"`
	DryRun        bool                   `json:"dry_run"`
	Committed     bool                   `json:"committed"`
	InputHash     string                 `json:"input_hash"`
	CorrelationID string                 `json:"correlation_id"`
	Inserted      int                    `json:"inserted"`
	Updated       int                    `json:"updated"`
	NoOp          int                    `json:"no_op"`
	Conflicts     int                    `json:"conflicts"`
	Unsupported   int                    `json:"unsupported"`
	Deleted       int                    `json:"deleted,omitempty"`
	PreservedL1   int                    `json:"preserved_l1,omitempty"`
	Results       []ScenarioObjectResult `json:"results"`
}

// ApplyScenario invokes the DES-047 governed, tenant-scoped scenario endpoint.
// request is deliberately accepted as a wire value so the content projection
// remains owned by AESE's legacyprojection package.
func (c *Client) ApplyScenario(ctx context.Context, request any, correlationID string) (ScenarioSummary, error) {
	var out ScenarioSummary
	err := c.requestHeaders(ctx, http.MethodPost, "api/v1/scenarios/apply", request, &out, map[string]string{"X-Correlation-ID": correlationID})
	return out, err
}

type ScenarioResetRequest struct {
	PackKey     string `json:"pack_key"`
	PackVersion string `json:"pack_version"`
	ScenarioKey string `json:"scenario_key"`
	RunID       string `json:"run_id"`
	DryRun      bool   `json:"dry_run"`
}

type MetadataSchemaRequest struct {
	DisplayName string          `json:"display_name"`
	Fields      json.RawMessage `json:"fields"`
}

func (c *Client) UpsertMetadataSchema(ctx context.Context, entity string, request MetadataSchemaRequest) error {
	if err := validateEntity(entity); err != nil {
		return err
	}
	if strings.TrimSpace(request.DisplayName) == "" || len(request.Fields) == 0 {
		return fmt.Errorf("metadata display_name and fields are required")
	}
	return c.request(ctx, http.MethodPost, "api/v1/metadata/schema/"+url.PathEscape(entity), request, nil)
}

type AIToolManifest struct {
	ToolKey            string          `json:"tool_key"`
	DisplayName        string          `json:"display_name"`
	ToolType           string          `json:"tool_type"`
	SourceRef          string          `json:"source_ref"`
	Description        string          `json:"description,omitempty"`
	InputSchema        json.RawMessage `json:"input_schema"`
	OutputSchema       json.RawMessage `json:"output_schema,omitempty"`
	RiskLevel          string          `json:"risk_level"`
	ConfirmationMode   string          `json:"confirmation_mode"`
	PermissionResource string          `json:"permission_resource,omitempty"`
	DecisionScope      json.RawMessage `json:"decision_scope,omitempty"`
	Examples           json.RawMessage `json:"examples,omitempty"`
	Enabled            bool            `json:"enabled"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
}

type AIToolSummary struct {
	ToolKey string `json:"tool_key"`
}

func (c *Client) ListAITools(ctx context.Context, includeDisabled bool) ([]AIToolSummary, error) {
	path := "api/v1/ai/tools?limit=200"
	if includeDisabled {
		path += "&include_disabled=true"
	}
	var response struct {
		Items []AIToolSummary `json:"items"`
	}
	if err := c.request(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (c *Client) CreateAITool(ctx context.Context, manifest AIToolManifest) error {
	return c.request(ctx, http.MethodPost, "api/v1/ai/tools", manifest, nil)
}

func (c *Client) UpdateAITool(ctx context.Context, manifest AIToolManifest) error {
	patch := map[string]any{
		"display_name": manifest.DisplayName, "description": manifest.Description,
		"input_schema": manifest.InputSchema, "risk_level": manifest.RiskLevel,
		"confirmation_mode":   manifest.ConfirmationMode,
		"permission_resource": manifest.PermissionResource,
		"examples":            manifest.Examples, "metadata": manifest.Metadata,
	}
	return c.request(ctx, http.MethodPatch, "api/v1/ai/tools/"+url.PathEscape(manifest.ToolKey), patch, nil)
}

func (c *Client) EnableAITool(ctx context.Context, toolKey string) error {
	return c.request(ctx, http.MethodPost, "api/v1/ai/tools/"+url.PathEscape(toolKey)+"/enable", map[string]any{}, nil)
}

type AIToolCallResult struct {
	CallID       string          `json:"call_id"`
	Status       string          `json:"status"`
	Output       json.RawMessage `json:"output"`
	ExecutionRef json.RawMessage `json:"execution_ref"`
}

func (c *Client) CallAITool(ctx context.Context, toolKey, correlationID, sessionID string, input any) (AIToolCallResult, error) {
	var out AIToolCallResult
	body := map[string]any{"input": input, "session_id": sessionID}
	err := c.requestHeaders(ctx, http.MethodPost, "api/v1/ai/tools/"+url.PathEscape(toolKey)+"/call", body, &out, map[string]string{"X-Correlation-ID": correlationID})
	if err == nil && (out.CallID == "" || out.Status != "succeeded") {
		return out, fmt.Errorf("AI tool %s returned incomplete success", toolKey)
	}
	return out, err
}

func (c *Client) ResetScenario(ctx context.Context, request ScenarioResetRequest, correlationID string) (ScenarioSummary, error) {
	var out ScenarioSummary
	err := c.requestHeaders(ctx, http.MethodPost, "api/v1/scenarios/reset", request, &out, map[string]string{"X-Correlation-ID": correlationID})
	return out, err
}

// FindExact returns tenant-scoped rows matching all supplied logical fields.
func (c *Client) FindExact(ctx context.Context, entity string, match map[string]any) ([]map[string]any, error) {
	if err := validateEntity(entity); err != nil {
		return nil, err
	}
	if len(match) == 0 {
		return nil, fmt.Errorf("at least one match field is required")
	}
	return c.findExact(ctx, entity, match)
}

func (c *Client) findExact(ctx context.Context, entity string, match map[string]any) ([]map[string]any, error) {
	var matches []map[string]any
	for page := 1; ; page++ {
		query := url.Values{"page": {strconv.Itoa(page)}, "page_size": {"100"}}
		var response struct {
			Total int              `json:"total"`
			Data  []map[string]any `json:"data"`
		}
		path := "api/v1/entities/" + url.PathEscape(entity) + "/records?" + query.Encode()
		if err := c.request(ctx, http.MethodGet, path, nil, &response); err != nil {
			return nil, err
		}
		for _, record := range response.Data {
			if recordMatches(record, match) {
				matches = append(matches, record)
			}
		}
		if page*100 >= response.Total || len(response.Data) == 0 {
			return matches, nil
		}
	}
}

type APIError struct {
	Method     string
	Path       string
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("IAOS %s %s returned %d: %s", e.Method, e.Path, e.StatusCode, e.Message)
}

func (c *Client) request(ctx context.Context, method, path string, body any, out any) error {
	return c.requestHeaders(ctx, method, path, body, out, nil)
}

func (c *Client) requestHeaders(ctx context.Context, method, path string, body any, out any, headers map[string]string) error {
	rel, err := url.Parse(path)
	if err != nil {
		return err
	}
	endpoint := c.baseURL.ResolveReference(rel)
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode IAOS request: %w", err)
		}
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.tenantID != "" {
		req.Header.Set("X-IAOS-Tenant-Id", c.tenantID)
	}
	for name, value := range headers {
		if strings.TrimSpace(value) != "" {
			req.Header.Set(name, value)
		}
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("IAOS %s %s: %w", method, endpoint.EscapedPath(), err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return fmt.Errorf("read IAOS response: %w", err)
	}
	if len(data) > maxResponseBytes {
		return fmt.Errorf("IAOS response exceeds %d bytes", maxResponseBytes)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message := strings.TrimSpace(string(data))
		var envelope struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(data, &envelope) == nil && envelope.Error != "" {
			message = envelope.Error
		}
		return &APIError{Method: method, Path: endpoint.EscapedPath(), StatusCode: resp.StatusCode, Message: message}
	}
	if out == nil || len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("decode IAOS response for %s: %w", endpoint.EscapedPath(), err)
	}
	return nil
}

func validateEntity(entity string) error {
	if entity == "" {
		return fmt.Errorf("entity is required")
	}
	for i, r := range entity {
		if !((r >= 'a' && r <= 'z') || r == '_' || (i > 0 && r >= '0' && r <= '9')) {
			return fmt.Errorf("invalid entity %q", entity)
		}
	}
	return nil
}

func validateUpsert(req UpsertRequest) error {
	if err := validateEntity(req.Entity); err != nil {
		return err
	}
	if len(req.NaturalKey) == 0 {
		return fmt.Errorf("%s natural key is required", req.Entity)
	}
	if req.Record == nil {
		return fmt.Errorf("%s record is required", req.Entity)
	}
	for _, name := range req.NaturalKey {
		if err := validateEntity(name); err != nil {
			return fmt.Errorf("%s natural key: %w", req.Entity, err)
		}
		if value, ok := req.Record[name]; !ok || value == nil || value == "" {
			return fmt.Errorf("%s record is missing natural key field %q", req.Entity, name)
		}
	}
	return nil
}

func recordMatches(record, match map[string]any) bool {
	for key, value := range match {
		if !equivalentJSON(record[key], value) {
			return false
		}
	}
	return true
}

func equivalentJSON(a, b any) bool {
	if reflect.DeepEqual(a, b) {
		return true
	}
	left, lerr := json.Marshal(a)
	right, rerr := json.Marshal(b)
	if lerr != nil || rerr != nil {
		return false
	}
	var lv, rv any
	ld := json.NewDecoder(bytes.NewReader(left))
	ld.UseNumber()
	rd := json.NewDecoder(bytes.NewReader(right))
	rd.UseNumber()
	if ld.Decode(&lv) != nil || rd.Decode(&rv) != nil {
		return false
	}
	return fmt.Sprint(lv) == fmt.Sprint(rv)
}

func IsStatus(err error, status int) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == status
}
