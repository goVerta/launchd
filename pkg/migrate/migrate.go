package migrate

import (
	"fmt"

	"github.com/goVerta/launchd/internal/executil"
)

func RunIfPresent(exec executil.Executor, user, host string) error {
	if exec == nil { return fmt.Errorf("nil executor") }
	// saadhtiwana: run only if migrate is available to avoid hard-fail in MVP
	cmd := fmt.Sprintf("ssh %s@%s 'command -v migrate >/dev/null 2>&1 && migrate up || true'", user, host)
	return exec.Run("bash", "-lc", cmd)
}
