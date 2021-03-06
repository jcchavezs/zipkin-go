package zipkin

import "testing"

func TestTraceID(t *testing.T) {

	traceID := TraceID{High: 1, Low: 2}

	if len(traceID.ToHex()) != 32 {
		t.Errorf("Expected zero-padded TraceID to have 32 characters")
	}

	have, err := TraceIDFromHex(traceID.ToHex())
	if err != nil {
		t.Fatalf("Expected traceID got error: %+v", err)
	}
	if traceID.High != have.High || traceID.Low != have.Low {
		t.Errorf("Expected %+v, got %+v", traceID, have)
	}

	traceID = TraceID{High: 0, Low: 2}

	if len(traceID.ToHex()) != 16 {
		t.Errorf("Expected zero-padded TraceID to have 16 characters, got %d", len(traceID.ToHex()))
	}

	have, err = TraceIDFromHex(traceID.ToHex())
	if err != nil {
		t.Fatalf("Expected traceID got error: %+v", err)
	}
	if traceID.High != have.High || traceID.Low != have.Low {
		t.Errorf("Expected %+v, got %+v", traceID, have)
	}

	traceID = TraceID{High: 0, Low: 0}

	if !traceID.Empty() {
		t.Errorf("Expected TraceID to be empty")
	}

	if _, err = TraceIDFromHex("12345678901234zz12345678901234zz"); err == nil {
		t.Errorf("Expected error got nil")
	}

}
