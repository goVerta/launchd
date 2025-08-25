package builder

import (
	"fmt"
	"path/filepath"

	"github.com/goVerta/launchd/internal/executil"
)

type Builder interface {
	Build(appPath, out string) error
}

type GoBuilder struct{ Exec executil.Executor }

func Build(exec executil.Executor, appPath, out string) error {
	if exec == nil { return fmt.Errorf("nil executor") }
	abs, err := filepath.Abs(appPath)
	if err != nil { return err }
	// saadhtiwana: build using system go toolchain for determinism
	return exec.Run("go", "build", "-o", out, abs)
}
