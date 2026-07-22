// Package knowledge projects actor-scoped observations without exposing full World State.
package knowledge

import (
	"fmt"
	"sort"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

type Store struct {
	records map[string][]worldcontract.Knowledge
}

func New() *Store { return &Store{records: map[string][]worldcontract.Knowledge{}} }
func actorKey(ref worldcontract.StableRef) string {
	return ref.Namespace + ":" + ref.Type + ":" + ref.Code
}
func (s *Store) Learn(record worldcontract.Knowledge) error {
	if err := record.Validate(); err != nil {
		return err
	}
	key := actorKey(record.ActorRef)
	for _, existing := range s.records[key] {
		if existing.KnowledgeID == record.KnowledgeID {
			return fmt.Errorf("duplicate knowledge_id %s", record.KnowledgeID)
		}
	}
	s.records[key] = append(s.records[key], record)
	sort.Slice(s.records[key], func(i, j int) bool { return s.records[key][i].ObservedAt < s.records[key][j].ObservedAt })
	return nil
}
func (s *Store) Visible(actor worldcontract.StableRef, at time.Time) []worldcontract.Knowledge {
	var out []worldcontract.Knowledge
	for _, record := range s.records[actorKey(actor)] {
		observed, _ := time.Parse(time.RFC3339, record.ObservedAt)
		if !observed.After(at) {
			out = append(out, record)
		}
	}
	return out
}
