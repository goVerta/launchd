package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/goVerta/launchd/pkg/health"
)

func TestHealthWait_Succeeds(t *testing.T) {
	h := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" { w.WriteHeader(404); return }
		w.WriteHeader(200)
	}))
	defer h.Close()

	u, err := url.Parse(h.URL)
	if err != nil { t.Fatalf("parse url: %v", err) }
	hostPort := strings.TrimPrefix(u.Host, "[")
	hostPort = strings.TrimSuffix(hostPort, "]")
	var host string
	var port int
	// host:port form already; split last colon
	parts := strings.Split(hostPort, ":")
	if len(parts) < 2 { t.Fatalf("bad host:port: %q", hostPort) }
	host = strings.Join(parts[:len(parts)-1], ":")
	if host == "" { host = "127.0.0.1" }
	// parse port
	var p int
	_, err = fmt.Sscanf(parts[len(parts)-1], "%d", &p)
	if err != nil { t.Fatalf("parse port: %v", err) }
	port = p

	if err := health.Wait(host, port, 2*time.Second); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestHealthWait_TimesOut(t *testing.T) {
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	u, err := url.Parse(s.URL)
	if err != nil { t.Fatalf("parse url: %v", err) }
	hostPort := strings.TrimPrefix(u.Host, "[")
	hostPort = strings.TrimSuffix(hostPort, "]")
	var host string
	var port int
	parts := strings.Split(hostPort, ":")
	if len(parts) < 2 { t.Fatalf("bad host:port: %q", hostPort) }
	host = strings.Join(parts[:len(parts)-1], ":")
	if host == "" { host = "127.0.0.1" }
	var p int
	_, err = fmt.Sscanf(parts[len(parts)-1], "%d", &p)
	if err != nil { t.Fatalf("parse port: %v", err) }
	port = p

	deadline := 1500 * time.Millisecond
	if err := health.Wait(host, port, deadline); err == nil {
		t.Fatalf("expected timeout")
	}
}
