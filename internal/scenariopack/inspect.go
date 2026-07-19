package scenariopack

import "sort"

type Summary struct {
	PackKey        string         `json:"pack_key"`
	PackVersion    string         `json:"pack_version"`
	TenantTemplate string         `json:"tenant_template"`
	Entities       map[string]int `json:"entities"`
	MasterRecords  int            `json:"master_records"`
	InitialRecords int            `json:"initial_records"`
	Stories        []StorySummary `json:"stories"`
}

type StorySummary struct {
	Key           string `json:"key"`
	Events        int    `json:"events"`
	Assertions    int    `json:"assertions"`
	CorrelationID string `json:"correlation_id"`
}

func Inspect(pack *Pack) Summary {
	summary := Summary{PackKey: pack.Manifest.PackKey, PackVersion: pack.Manifest.PackVersion, TenantTemplate: pack.Manifest.TenantTemplate, Entities: map[string]int{}}
	for _, set := range pack.RecordSets {
		summary.Entities[set.Entity] += len(set.Records)
		summary.MasterRecords += len(set.Records)
	}
	for _, story := range pack.Stories {
		for _, set := range story.Initial.RecordSets {
			summary.Entities[set.Entity] += len(set.Records)
			summary.InitialRecords += len(set.Records)
		}
		summary.Stories = append(summary.Stories, StorySummary{Key: story.Ref.Key, Events: len(story.Events.Events), Assertions: len(story.Expected.Assertions), CorrelationID: story.Events.CorrelationID})
	}
	sort.Slice(summary.Stories, func(i, j int) bool { return summary.Stories[i].Key < summary.Stories[j].Key })
	return summary
}
