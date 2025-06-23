package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

// func TestRootCommand(t *testing.T) {
// 	cmd := rootCmd

//		if cmd.Use != "k8s-controller-tutorial" {
//			t.Errorf("Expected command use 'k8s-controller-tutorial', got %s", cmd.Use)
//		}
//	}
func TestRootCommand_Execution(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{}) // без флагов и аргументов

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Welcome") {
		t.Errorf("Expected output to contain 'Welcome', got %q", output)
	}
}

func TestRootCommand_WithDebugFlag(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--log-level=debug"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error with debug flag, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Welcome") {
		t.Errorf("Expected output to contain 'Welcome', got %q", output)
	}
}

func TestRootCommand_InvalidFlag(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"--non-existent"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for unknown flag, got nil")
	}
}

func TestHandler_GET(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.SetMethod("GET")

	homeHandler(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Errorf("Expected 200, got %d", ctx.Response.StatusCode())
	}

	body := string(ctx.Response.Body())
	if !strings.Contains(body, "Hello") {
		t.Errorf("Unexpected response body: %q", body)
	}
}
