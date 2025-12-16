package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// version is injected at build time
var version = "dev"

// arrayFlags allows multiple values for a flag
type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var (
		localPort   string
		targetAddr  string
		enableCors  bool
		headers     arrayFlags
		showVersion bool
	)

	flag.StringVar(&localPort, "port", "8080", "Local port to listen on")
	flag.StringVar(&targetAddr, "target", "", "Target server URL (e.g. http://localhost:9000)")
	flag.BoolVar(&enableCors, "cors", false, "Enable default CORS setup")
	flag.Var(&headers, "header", "Custom header to add to response (Key:Value), can be used multiple times")
	flag.BoolVar(&showVersion, "version", false, "Print version and exit")
	flag.Parse()

	if showVersion {
		fmt.Println("fwd-proxy version:", version)
		os.Exit(0)
	}

	if targetAddr == "" {
		fmt.Println("Error: target URL is required")
		flag.Usage()
		os.Exit(1)
	}

	targetURL, err := url.Parse(targetAddr)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Custom Director to ensure paths are preserved/modified if needed.
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Ensure Host header matches target
		req.Host = targetURL.Host
		log.Printf("Forwarding request: %s %s -> %s", req.Method, req.URL.Path, targetAddr)
	}

	// ModifyResponse to add CORS and Custom Headers
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Add Custom Headers
		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				resp.Header.Set(key, value)
			}
		}

		// Add Default CORS Headers if enabled
		if enableCors {
			resp.Header.Set("Access-Control-Allow-Origin", "*")
			resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			resp.Header.Set("Access-Control-Allow-Headers", "*")
		}

		return nil
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If CORS is enabled, handle OPTIONS requests directly
		if enableCors && r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		proxy.ServeHTTP(w, r)
	})

	log.Printf("fwd-proxy %s started", version)
	log.Printf("HTTP Proxy listening on :%s, forwarding to %s", localPort, targetAddr)
	if enableCors {
		log.Println("CORS enabled")
	}
	if len(headers) > 0 {
		log.Printf("Injecting headers: %v", headers)
	}

	if err := http.ListenAndServe(":"+localPort, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
