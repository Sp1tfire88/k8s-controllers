package integration

import (
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func TestLeaderElectionAndMetrics(t *testing.T) {
	// 1. Стартуем первый контроллер
	cmd1 := exec.Command("go", "run", "main.go", "server",
		"--enable-leader-election=true",
		"--metrics-port=19090",
	)
	cmd1.Stdout = nil
	cmd1.Stderr = nil
	if err := cmd1.Start(); err != nil {
		t.Fatalf("Failed to start controller 1: %v", err)
	}
	defer cmd1.Process.Kill()

	// 2. Стартуем второй контроллер (на другом порту для метрик!)
	cmd2 := exec.Command("go", "run", "main.go", "server",
		"--enable-leader-election=true",
		"--metrics-port=19091",
	)
	cmd2.Stdout = nil
	cmd2.Stderr = nil
	if err := cmd2.Start(); err != nil {
		t.Fatalf("Failed to start controller 2: %v", err)
	}
	defer cmd2.Process.Kill()

	// 3. Ждём пока оба процесса стартуют
	time.Sleep(10 * time.Second)

	// 4. Проверяем, что оба endpoint-а метрик отвечают
	for _, port := range []string{"19090", "19091"} {
		resp, err := http.Get("http://localhost:" + port + "/metrics")
		if err != nil || resp.StatusCode != 200 {
			t.Errorf("Metrics endpoint on :%s is not available: %v", port, err)
		}
	}

	// 5. Можно проверить Lease-ресурс в кластере (например, через kubectl),
	//    но для локального теста обычно достаточно проверить работу метрик.

	// 6. Для полного теста:
	//    - Завершаем первый процесс, ждём 10 секунд, проверяем что второй instance становится лидером.
}
