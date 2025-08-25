package transfer

import (
    "fmt"
    "path"

    "github.com/goVerta/launchd/internal/executil"
)

func SCP(exec executil.Executor, localPath, user, host, remotePath string) error {
    if exec == nil { return fmt.Errorf("nil executor") }
    // Copy to a temp location first, then install atomically with proper perms.
    // This avoids fragile remote shell quoting/expansion semantics.
    tmpName := path.Base(remotePath)
    tmpRemote := fmt.Sprintf("/tmp/%s", tmpName)

    // 1) scp binary to /tmp
    if err := exec.Run("scp", "-q", localPath, fmt.Sprintf("%s@%s:%s", user, host, tmpRemote)); err != nil {
        return err
    }
    // 2) sudo install to final path (creates dirs, sets perms, and moves atomically)
    installCmd := fmt.Sprintf("ssh %s@%s 'sudo install -D -m 0755 -o root -g root %s %s'", user, host, tmpRemote, remotePath)
    return exec.Run("bash", "-lc", installCmd)
}
