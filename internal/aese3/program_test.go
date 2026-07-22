package aese3

import (
	"encoding/json"
	"testing"
)

func TestProgramClosesEveryMilestoneDeterministically(t *testing.T) {
	p := BuildProgram()
	if err := Validate(p); err != nil {
		t.Fatal(err)
	}
	first := p.ProgramHash
	for i := 0; i < 100; i++ {
		if got := BuildProgram().ProgramHash; got != first {
			t.Fatalf("unstable hash: %s != %s", got, first)
		}
	}
}

func TestProgramRejectsAutomaticBusinessWrite(t *testing.T) {
	p := BuildProgram()
	p.Milestones[0].BusinessWrites = 1
	if Validate(p) == nil {
		t.Fatal("expected ownership violation")
	}
}

func TestProgramRejectsSkippedTerminal(t *testing.T) {
	p := BuildProgram()
	p.Milestones[4].TerminalReady = false
	if Validate(p) == nil {
		t.Fatal("expected open terminal")
	}
}

func TestStrictParser(t *testing.T) {
	b, err := json.Marshal(BuildProgram())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ParseStrict(b); err != nil {
		t.Fatal(err)
	}
	bad := append(b[:len(b)-1], []byte(`,"unknown":true}`)...)
	if _, err = ParseStrict(bad); err == nil {
		t.Fatal("unknown field accepted")
	}
}
