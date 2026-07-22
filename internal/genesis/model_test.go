package genesis

import "testing"

func TestTraceAndConservation(t *testing.T) {
	if err := ValidateTrace(BuildTrace()); err != nil {
		t.Fatal(err)
	}
	if err := ValidateConservation(Conservation{Opening: "100.00", Inbound: "20.00", Consumed: "70.00", Loss: "5.00", Closing: "45.00"}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateConservation(Conservation{Opening: "1", Inbound: "0", Consumed: "0", Loss: "0", Closing: "2"}); err == nil {
		t.Fatal("imbalance accepted")
	}
}
