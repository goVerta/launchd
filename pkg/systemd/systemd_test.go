package systemd

import "testing"

type fakeExec struct{ calls [][]string; err error }

func (f *fakeExec) Run(name string, args ...string) error {
	f.calls = append(f.calls, append([]string{name}, args...))
	return f.err
}
func (f *fakeExec) RunOutput(name string, args ...string) (string, error) { return "", nil }

func TestSetup_IssuesSSH(t *testing.T) {
	ex := &fakeExec{}
	s := ServiceSpec{Name: "app", ExecStart: "/usr/local/bin/app", User: "root"}
	if err := Setup(ex, "alice", "host", s); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(ex.calls) != 1 { t.Fatalf("calls=%d", len(ex.calls)) }
	if ex.calls[0][0] != "bash" { t.Fatalf("expected bash -lc, got %v", ex.calls[0]) }
}
