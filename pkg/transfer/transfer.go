package transfer

import (
	"fmt"

	"github.com/goVerta/launchd/internal/executil"
)

func SCP(exec executil.Executor, localPath, user, host, remotePath string) error {
	if exec == nil { return fmt.Errorf("nil executor") }
	// saadhtiwana: ensure remote dir then copy atomically via scp
	mk := fmt.Sprintf("ssh %s@%s 'sudo mkdir -p $(dirname %s) && sudo chown %s:$(id -gn %s) $(dirname %s)'", user, host, remotePath, user, user, remotePath)
	if err := exec.Run("bash", "-lc", mk); err != nil { return err }
	return exec.Run("scp", "-q", localPath, fmt.Sprintf("%s@%s:%s", user, host, remotePath))
}
