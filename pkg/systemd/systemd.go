package systemd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/goVerta/launchd/internal/executil"
)

type ServiceSpec struct {
	Name        string
	ExecStart   string
	User        string
	After       []string
	Environment map[string]string
}

func unitContent(s ServiceSpec) string {
	var env []string
	for k, v := range s.Environment { env = append(env, fmt.Sprintf("Environment=%s=%s", k, v)) }
	sort.Strings(env)
	lines := []string{
		"[Unit]",
		fmt.Sprintf("Description=%s service", s.Name),
		fmt.Sprintf("After=%s", strings.Join(append(s.After, "network.target"), " ")),
		"[Service]",
		"Type=simple",
		fmt.Sprintf("User=%s", s.User),
		fmt.Sprintf("ExecStart=%s", s.ExecStart),
		"Restart=always",
		"RestartSec=2",
	}
	lines = append(lines, env...)
	lines = append(lines,
		"[Install]",
		"WantedBy=multi-user.target",
	)
	return strings.Join(lines, "\n") + "\n"
}

func Setup(exec executil.Executor, user, host string, spec ServiceSpec) error {
	if exec == nil { return fmt.Errorf("nil executor") }
	unit := unitContent(spec)
	path := fmt.Sprintf("/etc/systemd/system/%s.service", spec.Name)
	// saadhtiwana: write unit via sudo tee for root-ownership, then enable+restart
	cmd := fmt.Sprintf("ssh %s@%s 'echo %q | sudo tee %s >/dev/null && sudo systemctl daemon-reload && sudo systemctl enable %s && sudo systemctl restart %s'", user, host, unit, path, spec.Name, spec.Name)
	return exec.Run("bash", "-lc", cmd)
}
