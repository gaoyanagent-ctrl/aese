package scenarioknowledge

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Manifest struct {
	SchemaVersion  string    `json:"schema_version"`
	EditionKey     string    `json:"edition_key"`
	EditionVersion string    `json:"edition_version"`
	IAOSEdition    string    `json:"iaos_edition"`
	IAOSProcess    string    `json:"iaos_process"`
	Language       string    `json:"language"`
	ContentHash    string    `json:"content_hash"`
	Articles       []Article `json:"articles"`
	Nodes          []Node    `json:"nodes"`
}
type Article struct {
	ArticleID        string `json:"article_id"`
	Title            string `json:"title"`
	SourceRef        string `json:"source_ref"`
	AppliesToVersion string `json:"applies_to_version"`
}
type Node struct {
	Sequence      int      `json:"sequence"`
	Capability    string   `json:"capability"`
	TaskType      string   `json:"task_type"`
	Actor         string   `json:"actor"`
	Gate          string   `json:"gate,omitempty"`
	ArticleID     string   `json:"article_id"`
	Purpose       string   `json:"purpose"`
	Inputs        []string `json:"inputs"`
	Outputs       []string `json:"outputs"`
	EvidenceTypes []string `json:"evidence_types"`
	IAOSMenus     []string `json:"iaos_menus"`
	WorldActions  []string `json:"world_actions"`
}

var allowedCapabilities = map[string]bool{
	"incorporation.case.open": true, "founder.resolution.prepare": true, "founder.resolution.approve": true, "capital.commitment.record": true,
	"registration.package.validate": true, "registration.submit": true, "registration.observation.commit": true, "bank.account.opening.submit": true,
	"bank.account.observation.commit": true, "capital.contribution.verify": true, "organization.establish": true, "executive.appointment.propose": true,
	"executive.appointment.acceptance.commit": true, "executive.appointment.approve": true, "operating.mandate.grant": true, "initial.budget.prepare": true,
	"initial.budget.approve": true, "enterprise.readiness.evaluate": true,
}

func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(strings.NewReader(string(data)))
	d.DisallowUnknownFields()
	var m Manifest
	if err := d.Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
func (m Manifest) Digest() string {
	m.ContentHash = ""
	sort.Slice(m.Articles, func(i, j int) bool { return m.Articles[i].ArticleID < m.Articles[j].ArticleID })
	sort.Slice(m.Nodes, func(i, j int) bool { return m.Nodes[i].Sequence < m.Nodes[j].Sequence })
	raw, _ := json.Marshal(m)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
func (m Manifest) Validate() []error {
	errs := []error{}
	add := func(f string, v any) { errs = append(errs, fmt.Errorf("%s: %v", f, v)) }
	if m.SchemaVersion != "1.0" {
		add("schema_version", "must be 1.0")
	}
	if m.EditionKey == "" || m.EditionVersion == "" {
		add("edition", "key and version required")
	}
	if m.IAOSProcess != "enterprise.incorporation.lifecycle.v1" {
		add("iaos_process", "unsupported process")
	}
	if m.ContentHash != m.Digest() {
		add("content_hash", "does not match canonical manifest")
	}
	articleIDs := map[string]bool{}
	for _, a := range m.Articles {
		if a.ArticleID == "" || a.SourceRef == "" {
			add("articles", a.ArticleID+" missing source")
		}
		if articleIDs[a.ArticleID] {
			add("articles", a.ArticleID+" duplicated")
		}
		articleIDs[a.ArticleID] = true
	}
	if len(m.Nodes) != 18 {
		add("nodes", fmt.Sprintf("got %d, want 18", len(m.Nodes)))
	}
	seen := map[string]bool{}
	for i, n := range m.Nodes {
		if n.Sequence != i+1 {
			add("nodes.sequence", fmt.Sprintf("got %d at index %d", n.Sequence, i))
		}
		if !allowedCapabilities[n.Capability] {
			add("nodes.capability", n.Capability)
		}
		if seen[n.Capability] {
			add("nodes.capability", n.Capability+" duplicated")
		}
		seen[n.Capability] = true
		if !articleIDs[n.ArticleID] {
			add("nodes.article_id", n.ArticleID+" missing")
		}
		if n.Purpose == "" || len(n.Inputs) == 0 || len(n.Outputs) == 0 || len(n.EvidenceTypes) == 0 || len(n.IAOSMenus) == 0 {
			add("nodes.contract", n.Capability+" incomplete")
		}
	}
	return errs
}
