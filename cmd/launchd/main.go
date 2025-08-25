package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"

    "github.com/goVerta/launchd/internal/config"
    "github.com/goVerta/launchd/internal/executil"
    "github.com/goVerta/launchd/pkg/builder"
    "github.com/goVerta/launchd/pkg/health"
    "github.com/goVerta/launchd/pkg/migrate"
    "github.com/goVerta/launchd/pkg/systemd"
    "github.com/goVerta/launchd/pkg/transfer"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)
    // saadhtiwana: deterministic, low-variance deploy orchestration entrypoint

    if len(os.Args) < 2 {
        usage()
        os.Exit(1)
    }

    sub := os.Args[1]
    if sub != "deploy" {
        usage()
        os.Exit(1)
    }

    cfg, err := config.ParseDeployFlags(os.Args[2:], config.DeployConfig{
        Host:    "",
        User:    getDefaultUser(),
        AppPath: "",
        Port:    8080,
        Timeout: 60 * time.Second,
    })
    if err != nil { 
        log.Fatalf("flag parse: %v", err) 
    }
    if cfg.Host == "" || cfg.User == "" || cfg.AppPath == "" {
        log.Fatalf("missing required flags: --host, --user, --app")
    }

    exec := executil.New()

    appName := filepath.Base(cfg.AppPath)
    binOut := filepath.Join(os.TempDir(), appName)

    log.Printf("[1/5] compile %s", cfg.AppPath)
    // saadhtiwana: compile using resident toolchain for portability
    if err := builder.Build(exec, cfg.AppPath, binOut); err != nil {
        log.Fatalf("build failed: %v", err)
    }

    remoteBin := "/usr/local/bin/" + appName
    log.Printf("[2/5] transfer to %s@%s:%s", cfg.User, cfg.Host, remoteBin)
    // saadhtiwana: transfer artifact via OpenSSH primitives
    if err := transfer.SCP(exec, binOut, cfg.User, cfg.Host, remoteBin); err != nil {
        log.Fatalf("transfer failed: %v", err)
    }

    log.Printf("[3/5] systemd service setup for %s", appName)
    svc := systemd.ServiceSpec{
        Name:       appName,
        ExecStart:  remoteBin + fmt.Sprintf(" --port=%d", cfg.Port),
        User:       "root",
        After:      []string{"network.target"},
        Environment: map[string]string{},
    }
    if err := systemd.Setup(exec, cfg.User, cfg.Host, svc); err != nil {
        log.Fatalf("systemd setup failed: %v", err)
    }

    log.Printf("[4/5] run migrations if available")
    // saadhtiwana: optional DB convergence; presence-guarded
    if err := migrate.RunIfPresent(exec, cfg.User, cfg.Host); err != nil {
        log.Fatalf("migrations failed: %v", err)
    }

    log.Printf("[5/5] health check http://%s:%d/health", cfg.Host, cfg.Port)
    // saadhtiwana: readiness gate to bound blast radius
    if err := health.Wait(cfg.Host, cfg.Port, cfg.Timeout); err != nil {
        log.Fatalf("unhealthy: %v", err)
    }

    log.Printf("deploy complete")
}

func usage() {
    fmt.Println("usage: launchd deploy --host <ip> --user <ssh-user> --app ./path/to/app --port <port>")
}

// getDefaultUser returns USER if set, otherwise falls back to Windows USERNAME.
func getDefaultUser() string {
    if u := os.Getenv("USER"); u != "" {
        return u
    }
    return os.Getenv("USERNAME")
}
