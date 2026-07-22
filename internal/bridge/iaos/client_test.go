package iaos

import (
	"context"
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestObservationRetryAndCursorRead(t *testing.T) {
	calls := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Error("missing token")
		}
		if r.Method == "POST" {
			calls++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(map[bool]int{true: 200, false: 201}[calls > 1])
			_, _ = w.Write([]byte(`{"message_id":"m1","journal_cursor":41,"accepted":true,"duplicate":` + map[bool]string{true: "true", false: "false"}[calls > 1] + `}`))
			return
		}
		_, _ = w.Write([]byte(`{"items":[],"next_cursor":41,"has_more":false}`))
	}))
	defer s.Close()
	c, _ := New(s.URL, "token")
	o := worldcontract.Observation(worldcontract.Envelope{MessageID: "m1"})
	a, err := c.PostObservation(context.Background(), o)
	if err != nil || a.JournalCursor != 41 {
		t.Fatal(err)
	}
	a, err = c.PostObservation(context.Background(), o)
	if err != nil || !a.Duplicate {
		t.Fatal("retry not duplicate")
	}
	p, err := c.Entries(context.Background(), "run", "main", 0)
	if err != nil || p.NextCursor != 41 {
		t.Fatal(err)
	}
}
