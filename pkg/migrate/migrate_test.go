package migrate

import "testing"

type fakeExec struct{ calls [][]string; err error }

func (f *fakeExec) Run(name string, args ...string) error {
	f.calls = append(f.calls, append([]string{name}, args...))
	return f.err
}
func (f *fakeExec) RunOutput(name string, args ...string) (string, error) { return "", nil }

func TestRunIfPresent_OK(t *testing.T) {
	ex := &fakeExec{}
	if err := RunIfPresent(ex, "alice", "host"); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(ex.calls) != 1 { t.Fatalf("calls=%d", len(ex.calls)) }
}
