package main_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestProxyIntegration(t *testing.T) {
	// 1. Build the binary
	tmpDir, err := os.MkdirTemp("", "proxy-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	proxyBin := filepath.Join(tmpDir, "proxy")
	buildCmd := exec.Command("go", "build", "-o", proxyBin, ".")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\n%s", err, out)
	}

	// 2. Start Target HTTP Server
	targetListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start target: %v", err)
	}
	targetPort := targetListener.Addr().(*net.TCPAddr).Port

	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "bar")
	})

	targetServer := &http.Server{Handler: mux}
	go targetServer.Serve(targetListener)
	defer targetServer.Close()

	// 3. Start Proxy
	proxyListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	proxyPort := proxyListener.Addr().(*net.TCPAddr).Port
	proxyListener.Close()

	targetURL := fmt.Sprintf("http://127.0.0.1:%d", targetPort)
	// Run with CORS enabled and a custom header
	proxyCmd := exec.Command(proxyBin,
		"-port", fmt.Sprintf("%d", proxyPort),
		"-target", targetURL,
		"-cors",
		"-header", "X-Custom-Foo: CustomBar",
	)
	proxyCmd.Stdout = os.Stdout
	proxyCmd.Stderr = os.Stderr
	if err := proxyCmd.Start(); err != nil {
		t.Fatalf("Failed to start proxy: %v", err)
	}
	defer func() {
		proxyCmd.Process.Kill()
	}()

	// Wait for start with retry
	var conn net.Conn
	for i := 0; i < 20; i++ {
		conn, err = net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", proxyPort))
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("Failed to connect to proxy after retries: %v", err)
	}

	// 4. Test Request
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/foo", proxyPort))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "bar" {
		t.Errorf("Expected 'bar', got %q", string(body))
	}

	// Verify CORS Headers
	if val := resp.Header.Get("Access-Control-Allow-Origin"); val != "*" {
		t.Errorf("Expected CORS header *, got %q", val)
	}

	// Verify Custom Header
	if val := resp.Header.Get("X-Custom-Foo"); val != "CustomBar" {
		t.Errorf("Expected custom header CustomBar, got %q", val)
	}

	// Verify OPTIONS request (Preflight)
	client := &http.Client{}
	req, _ := http.NewRequest("OPTIONS", fmt.Sprintf("http://127.0.0.1:%d/foo", proxyPort), nil)
	respOpt, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer respOpt.Body.Close()

	if respOpt.StatusCode != http.StatusNoContent {
		t.Errorf("Expected OPTIONS 204 No Content, got %d", respOpt.StatusCode)
	}
	if val := respOpt.Header.Get("Access-Control-Allow-Methods"); val == "" {
		t.Error("Expected CORS methods header in OPTIONS response")
	}
}
