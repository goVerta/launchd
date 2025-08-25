package transfer

import (
	"reflect"
	"testing"
)

type fakeExec struct{ calls [][]string; err error }

func (f *fakeExec) Run(name string, args ...string) error {
	f.calls = append(f.calls, append([]string{name}, args...))
	return f.err
}
func (f *fakeExec) RunOutput(name string, args ...string) (string, error) { return "", nil }

func TestSCP_Commands(t *testing.T) {
	ex := &fakeExec{}
	if err := SCP(ex, "/local/bin", "alice", "host", "/usr/local/bin/app"); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(ex.calls) != 2 { t.Fatalf("calls=%d", len(ex.calls)) }
	if ex.calls[0][0] != "bash" { t.Fatalf("want bash pre-mkdir, got %v", ex.calls[0]) }
	if ex.calls[1][0] != "scp" { t.Fatalf("want scp, got %v", ex.calls[1]) }
}
