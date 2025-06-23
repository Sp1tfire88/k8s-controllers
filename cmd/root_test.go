package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

// captureOutput перехватывает стандартный вывод (stdout)
func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = stdout
	_, _ = io.Copy(&buf, r)

	return buf.String()
}

// Тест на базовое выполнение root-команды
func TestRootCommand_Execution(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{}) // без флагов

	output := captureOutput(func() {
		_ = cmd.Execute()
	})

	if !strings.Contains(output, "Welcome") {
		t.Errorf("Expected output to contain 'Welcome', got %q", output)
	}
}

// Тест с использованием флага --log-level=debug
func TestRootCommand_WithDebugFlag(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"--log-level=debug"})

	output := captureOutput(func() {
		_ = cmd.Execute()
	})

	if !strings.Contains(output, "Welcome") {
		t.Errorf("Expected output to contain 'Welcome', got %q", output)
	}
}

// Проверка, что неизвестный флаг вызывает ошибку
func TestRootCommand_InvalidFlag(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"--non-existent"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for unknown flag, got nil")
	}
}

// Тест GET-запроса к homeHandler
func TestHomeHandler_GET(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")
	ctx.Request.Header.SetMethod("GET")

	homeHandler(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Errorf("Expected 200, got %d", ctx.Response.StatusCode())
	}

	body := string(ctx.Response.Body())
	expected := "Welcome to the FastHTTP server!"
	if body != expected {
		t.Errorf("Expected response %q, got %q", expected, body)
	}
}

// Тест POST-запроса к postHandler
func TestPostHandler_POST(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/post")
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.SetBody([]byte(`{"msg":"hi"}`))

	postHandler(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Errorf("Expected 200, got %d", ctx.Response.StatusCode())
	}

	body := string(ctx.Response.Body())
	expected := "POST received"
	if body != expected {
		t.Errorf("Expected response %q, got %q", expected, body)
	}
}

// Тест GET-запроса к healthHandler
func TestHealthHandler(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/health")
	ctx.Request.Header.SetMethod("GET")

	healthHandler(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Errorf("Expected 200, got %d", ctx.Response.StatusCode())
	}

	body := string(ctx.Response.Body())
	expected := "OK"
	if body != expected {
		t.Errorf("Expected response %q, got %q", expected, body)
	}
}
