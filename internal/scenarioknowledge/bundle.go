package scenarioknowledge

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// KnowledgeArticle mirrors the IAOS product knowledge wire contract. AESE
// owns scenario prose; IAOS remains the governed installer and resolver.
type KnowledgeArticle struct {
	ArticleID        string          `json:"article_id"`
	Title            string          `json:"title"`
	Summary          string          `json:"summary"`
	ContentType      string          `json:"content_type"`
	Audience         json.RawMessage `json:"audience"`
	Module           string          `json:"module"`
	MenuCode         string          `json:"menu_code"`
	Route            string          `json:"route"`
	AppliesToVersion string          `json:"applies_to_version"`
	Status           string          `json:"status"`
	Language         string          `json:"language"`
	Purpose          string          `json:"purpose"`
	Prerequisites    json.RawMessage `json:"prerequisites"`
	Steps            json.RawMessage `json:"steps"`
	Validation       json.RawMessage `json:"validation"`
	Recovery         json.RawMessage `json:"recovery"`
	RelatedAssets    json.RawMessage `json:"related_assets"`
	Keywords         json.RawMessage `json:"keywords"`
	SourceRefs       json.RawMessage `json:"source_refs"`
	BodyMarkdown     string          `json:"body_markdown"`
	Owner            string          `json:"owner"`
	SourceLayer      string          `json:"source_layer"`
	EditionKey       string          `json:"edition_key,omitempty"`
	ArticleVersion   int             `json:"article_version"`
	ReviewedAt       time.Time       `json:"reviewed_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type PackageAsset struct {
	Kind, Key   string
	Version     string `json:"version,omitempty"`
	Ownership   string `json:"ownership"`
	Description string `json:"description,omitempty"`
}

func (a PackageAsset) MarshalJSON() ([]byte, error) {
	type wire struct {
		Kind        string `json:"kind"`
		Key         string `json:"key"`
		Version     string `json:"version,omitempty"`
		Ownership   string `json:"ownership"`
		Description string `json:"description,omitempty"`
	}
	return json.Marshal(wire{a.Kind, a.Key, a.Version, a.Ownership, a.Description})
}

type PackageManifest struct {
	SchemaVersion string         `json:"schema_version"`
	PackageKey    string         `json:"package_key"`
	DisplayName   string         `json:"display_name"`
	Version       string         `json:"version"`
	Description   string         `json:"description"`
	DependsOn     []string       `json:"depends_on"`
	Assets        []PackageAsset `json:"assets"`
	ContentHash   string         `json:"content_hash"`
}
type PackageRef struct {
	PackageKey  string `json:"package_key"`
	Version     string `json:"version"`
	ContentHash string `json:"content_hash"`
}
type Edition struct {
	SchemaVersion string       `json:"schema_version"`
	EditionKey    string       `json:"edition_key"`
	DisplayName   string       `json:"display_name"`
	Version       string       `json:"version"`
	Description   string       `json:"description"`
	Packages      []PackageRef `json:"packages"`
	ContentHash   string       `json:"content_hash"`
}
type Bundle struct {
	SchemaVersion string             `json:"schema_version"`
	Edition       Edition            `json:"edition"`
	Packages      []PackageManifest  `json:"packages"`
	Articles      []KnowledgeArticle `json:"articles"`
}

func jsonDigest(value any) string {
	data, _ := json.Marshal(value)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
func signPackage(m PackageManifest) PackageManifest {
	m.SchemaVersion = "1.0"
	sort.Slice(m.Assets, func(i, j int) bool {
		if m.Assets[i].Kind == m.Assets[j].Kind {
			return m.Assets[i].Key < m.Assets[j].Key
		}
		return m.Assets[i].Kind < m.Assets[j].Kind
	})
	m.ContentHash = ""
	m.ContentHash = jsonDigest(m)
	return m
}
func signEdition(e Edition) Edition {
	e.SchemaVersion = "1.0"
	e.ContentHash = ""
	e.ContentHash = jsonDigest(e)
	return e
}
func rawJSON(value any) json.RawMessage { data, _ := json.Marshal(value); return data }

func findRepoRoot(manifestPath string) (string, error) {
	dir, _ := filepath.Abs(filepath.Dir(manifestPath))
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("cannot resolve AESE repository root")
		}
		dir = parent
	}
}
func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func CompileBundle(manifestPath string) (Bundle, error) {
	m, err := Load(manifestPath)
	if err != nil {
		return Bundle{}, err
	}
	if issues := m.Validate(); len(issues) > 0 {
		return Bundle{}, fmt.Errorf("invalid knowledge manifest: %v", issues)
	}
	root, err := findRepoRoot(manifestPath)
	if err != nil {
		return Bundle{}, err
	}
	articles := make([]KnowledgeArticle, 0, len(m.Articles))
	for _, source := range m.Articles {
		body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(source.SourceRef)))
		if err != nil {
			return Bundle{}, fmt.Errorf("read article %s: %w", source.ArticleID, err)
		}
		steps := []string{}
		capabilities := []string{}
		menus := []string{}
		worldActions := []string{}
		for _, node := range m.Nodes {
			steps = append(steps, fmt.Sprintf("%02d · %s · %s · %s", node.Sequence, node.Capability, node.Actor, node.Purpose))
			capabilities = append(capabilities, node.Capability)
			menus = append(menus, node.IAOSMenus...)
			worldActions = append(worldActions, node.WorldActions...)
		}
		articles = append(articles, KnowledgeArticle{ArticleID: source.ArticleID, Title: source.Title, Summary: "华辰 M9 企业设立 18 个节点的角色、输入、输出、审批门、World 动作和 IAOS 证据入口。", ContentType: "scenario", Audience: rawJSON([]string{"founder", "tester", "implementer"}), Module: "incorporation", MenuCode: "enterprise_lifecycle", Route: "/enterprise_lifecycle", AppliesToVersion: source.AppliesToVersion, Status: "active", Language: m.Language, Purpose: "让实施人员和业务用户逐节点理解并验收 AESE 与 IAOS 的人机协作闭环。", Prerequisites: rawJSON([]string{"目标租户已安装 " + m.IAOSEdition, "企业设立流程已发布", "用户具备知识读取权限"}), Steps: rawJSON(steps), Validation: rawJSON([]string{"18 个节点顺序完整", "每个节点可穿透到实际 IAOS 证据", "World 节点不由 IAOS 伪造"}), Recovery: rawJSON([]string{"版本不匹配时停止安装并升级 IAOS 基础 Edition", "缺失运行证据时保持失败关闭并检查节点全链"}), RelatedAssets: rawJSON(map[string][]string{"processes": {m.IAOSProcess}, "capabilities": capabilities, "menus": uniqueSorted(menus), "world_actions": uniqueSorted(worldActions)}), Keywords: rawJSON([]string{"M9", "企业设立", "HCTM", "AESE", "人工审批", "Agent 任务", "World Observation"}), SourceRefs: rawJSON([]string{source.SourceRef, "scenario-packs/hctm/knowledge/m9-incorporation.json"}), BodyMarkdown: string(body), Owner: "AESE HCTM Scenario Governance", SourceLayer: "industry", EditionKey: m.EditionKey, ArticleVersion: 1})
	}
	assets := make([]PackageAsset, 0, len(articles))
	for _, article := range articles {
		assets = append(assets, PackageAsset{Kind: "knowledge_article", Key: article.ArticleID, Version: jsonDigest(article), Ownership: "industry", Description: article.Title})
	}
	pkg := signPackage(PackageManifest{PackageKey: m.EditionKey, DisplayName: "华辰 M9 企业设立知识包", Version: m.EditionVersion, Description: "HCTM 行业场景说明和 18 节点 IAOS 证据映射。", DependsOn: []string{m.IAOSEdition}, Assets: assets})
	edition := signEdition(Edition{EditionKey: m.EditionKey, DisplayName: "华辰 M9 企业设立知识 Edition", Version: m.EditionVersion, Description: "按租户安装的 HCTM 场景知识，不包含业务实例数据。", Packages: []PackageRef{{PackageKey: pkg.PackageKey, Version: pkg.Version, ContentHash: pkg.ContentHash}}})
	return Bundle{SchemaVersion: "1.0", Edition: edition, Packages: []PackageManifest{pkg}, Articles: articles}, nil
}

func BundleJSON(bundle Bundle) ([]byte, error) { return json.MarshalIndent(bundle, "", "  ") }
