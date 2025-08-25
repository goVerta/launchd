package config

import (
	"flag"
	"time"
)

type DeployConfig struct {
	Host    string
	User    string
	AppPath string
	Port    int
	Timeout time.Duration
}

func ParseDeployFlags(args []string, defaults DeployConfig) (DeployConfig, error) {
	fs := flag.NewFlagSet("deploy", flag.ContinueOnError)
	host := fs.String("host", defaults.Host, "target host")
	user := fs.String("user", defaults.User, "ssh user")
	appPath := fs.String("app", defaults.AppPath, "path to Go app (main package directory)")
	port := fs.Int("port", defaults.Port, "health port")
	timeout := fs.Duration("timeout", defaults.Timeout, "health timeout")
	if err := fs.Parse(args); err != nil { return DeployConfig{}, err }
	return DeployConfig{Host: *host, User: *user, AppPath: *appPath, Port: *port, Timeout: *timeout}, nil
}
