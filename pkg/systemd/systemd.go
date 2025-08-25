package systemd

import (
	"fmt"
	"sort"
	"strings"
	"os"

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
	// Write unit to a local temporary file to preserve exact newlines/content
	tmpFile, err := os.CreateTemp("", spec.Name+"-*.service")
	if err != nil { return err }
	tmpName := tmpFile.Name()
	if _, err := tmpFile.WriteString(unit); err != nil { _ = tmpFile.Close(); _ = os.Remove(tmpName); return err }
	if err := tmpFile.Close(); err != nil { _ = os.Remove(tmpName); return err }
	defer os.Remove(tmpName)

	remoteTmp := fmt.Sprintf("/tmp/%s.service", spec.Name)
	finalPath := fmt.Sprintf("/etc/systemd/system/%s.service", spec.Name)

	// 1) Copy unit to remote tmp
	if err := exec.Run("scp", "-q", tmpName, fmt.Sprintf("%s@%s:%s", user, host, remoteTmp)); err != nil { return err }
	// 2) Install unit with proper ownership and permissions, then reload+enable+restart
	cmd := fmt.Sprintf("ssh %s@%s 'sudo install -D -m 0644 -o root -g root %s %s && sudo systemctl daemon-reload && sudo systemctl enable %s && sudo systemctl restart %s'", user, host, remoteTmp, finalPath, spec.Name, spec.Name)
	return exec.Run("bash", "-lc", cmd)
}
