package cmd

import "testing"

func TestRootCommand(t *testing.T) {
	cmd := rootCmd

	if cmd.Use != "k8s-controller-tutorial" {
		t.Errorf("Expected command use 'k8s-controller-tutorial', got %s", cmd.Use)
	}
}
