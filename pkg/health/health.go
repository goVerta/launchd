package health

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Wait(host string, port int, timeout time.Duration) error {
	client := &http.Client{ Timeout: 5 * time.Second }
	deadline := time.Now().Add(timeout)
	url := fmt.Sprintf("http://%s:%d/health", host, port)
	for {
		if time.Now().After(deadline) { return errors.New("timeout waiting for health") }
		resp, err := client.Get(url)
		if err == nil {
			b, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 { return nil }
			_ = b
		}
		time.Sleep(500 * time.Millisecond)
	}
}
