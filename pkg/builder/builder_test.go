package builder

import (
	"errors"
	"reflect"
	"testing"
)

type fakeExec struct{ calls [][]string; err error }

func (f *fakeExec) Run(name string, args ...string) error {
	f.calls = append(f.calls, append([]string{name}, args...))
	return f.err
}
func (f *fakeExec) RunOutput(name string, args ...string) (string, error) { return "", nil }

func TestBuild_InvokesGoBuild(t *testing.T) {
	ex := &fakeExec{}
	if err := Build(ex, "./testdata/app", "/tmp/out"); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(ex.calls) != 1 { t.Fatalf("calls=%d", len(ex.calls)) }
	got := ex.calls[0]
	if got[0] != "go" || !reflect.DeepEqual(got[1:3], []string{"build", "-o"}) {
		t.Fatalf("bad invocation: %v", got)
	}
}

func TestBuild_ErrorSurfaced(t *testing.T) {
	ex := &fakeExec{err: errors.New("boom")}
	if err := Build(ex, "./app", "/tmp/out"); err == nil {
		t.Fatalf("expected error")
	}
}
