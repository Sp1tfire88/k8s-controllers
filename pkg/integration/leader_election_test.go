// pkg/integration/leader_election_test.go
package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func TestLeaderElectionAndMetricsPort(t *testing.T) {
	bin := "./build/controller" // или "./controller" если без BUILD_DIR

	// Первый процесс: leader election enabled, метрики на 19091
	cmd1 := exec.Command(bin, "server",
		"--enable-leader-election=true",
		"--metrics-port=19091",
	)
	var out1 bytes.Buffer
	cmd1.Stdout = &out1
	cmd1.Stderr = &out1
	if err := cmd1.Start(); err != nil {
		t.Fatalf("failed to start first controller: %v\nOutput:\n%s", err, out1.String())
	}
	defer func() { _ = cmd1.Process.Kill() }()

	time.Sleep(2 * time.Second) // дать время cmd1 стать лидером

	// Второй процесс: leader election enabled, метрики на 19092
	cmd2 := exec.Command(bin, "server",
		"--enable-leader-election=true",
		"--metrics-port=19092",
	)
	var out2 bytes.Buffer
	cmd2.Stdout = &out2
	cmd2.Stderr = &out2
	if err := cmd2.Start(); err != nil {
		_ = cmd1.Process.Kill()
		t.Fatalf("failed to start second controller: %v\nOutput:\n%s", err, out2.String())
	}
	defer func() { _ = cmd2.Process.Kill() }()

	time.Sleep(4 * time.Second) // дать время для старта и выборов лидера

	// Проверим, что оба сервера метрик доступны
	checkPort := func(port int) error {
		url := fmt.Sprintf("http://localhost:%d/metrics", port)
		client := http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("metrics endpoint not available at %s: %w", url, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			return fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url)
		}
		return nil
	}

	if err := checkPort(19091); err != nil {
		t.Errorf("First controller metrics not available: %v", err)
	}
	if err := checkPort(19092); err != nil {
		t.Errorf("Second controller metrics not available: %v", err)
	}

	// Выведем вывод процессов для дебага
	t.Logf("First controller output:\n%s", out1.String())
	t.Logf("Second controller output:\n%s", out2.String())
}
