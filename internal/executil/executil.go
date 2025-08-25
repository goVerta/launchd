package executil

import (
	"bytes"
	"os/exec"
)

type Executor interface {
	Run(name string, args ...string) error
	RunOutput(name string, args ...string) (string, error)
}

type Default struct{}

func New() *Default { return &Default{} }

func (d *Default) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func (d *Default) RunOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
