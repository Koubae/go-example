package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.LUTC | log.Lmicroseconds | log.Lshortfile)

	// simpleOneUseIoReadAll()
	// simpleTwoUseBufioScanner()
	// simpleThreeJSONBody()
	customClientOne()

}

func simpleOneUseIoReadAll() {
	fmt.Println(
		`
=====================================================
	simpleOneUseIoReadAll
=====================================================`,
	)
	response, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %s\n", err.Error())
		}
	}()

	body, err := io.ReadAll(response.Body)
	fmt.Println(string(body))
}

func simpleTwoUseBufioScanner() {
	fmt.Println(
		`
=====================================================
	simpleTwoUseBufioScanner
=====================================================`,
	)
	response, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %s\n", err.Error())
		}
	}()

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		fmt.Println("=====================")
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func simpleThreeJSONBody() {
	fmt.Println(
		`
=====================================================
	simpleOneUseIoReadAll
=====================================================`,
	)
	response, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("Error closing Response body: %s\n", err.Error())
		}
	}()

	// var responseJSON struct {
	// 	Data json.RawMessage
	// }
	// var responseJSON []byte
	var responseJSON map[string]interface{}
	// var responseJSON map[string]any
	// var responseJSON any

	decoder := json.NewDecoder(response.Body)
	decoder.UseNumber() // keep numbers as json.Number (avoid float64 rounding)

	if err := decoder.Decode(&responseJSON); err != nil {
		log.Fatal(err)
	}
	for k, v := range responseJSON {

		m, ok := v.(map[string]interface{})
		if ok {
			fmt.Printf("%s:\n", k)
			for k2, v2 := range m {
				fmt.Printf("\tKey: %s, Value: %v\n", k2, v2)
			}
		} else {
			fmt.Printf("Key: %s, Value: %v\n", k, v)
		}

	}
}

// Generated with help of ChatGTP 5
func customClientOne() {
	fmt.Println(
		`
=====================================================
	customClientOne
=====================================================`,
	)
	client := NewHTTPClient()
	defer Shutdown(client)

	// Create context with Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/get", nil)
	if err != nil {
		log.Fatal(err)
	}

	// lazy way to generate a request ID
	requestID := strconv.FormatUint((uint64)(time.Now().UTC().UnixNano()), 36)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "go-example/1.0")
	req.Header.Set("X-Request-ID", requestID)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode > 299 {
		log.Fatal(resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()

	var body any
	if err := decoder.Decode(&body); err != nil {
		log.Fatal(err)
	}

	pretty, err := json.MarshalIndent(body, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", pretty)
}

func NewHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,  // connection timeout
		KeepAlive: 30 * time.Second, // TCP keep-alives
	}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,

		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,             // enable HTTP/2 if server supports it
		MaxIdleConns:          200,              // total idle conns across all hosts
		MaxIdleConnsPerHost:   20,               // per-host idle conns
		IdleConnTimeout:       90 * time.Second, // how long idle conns stay open
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second, // time to receive headers after request write

		// Optional: tighten TLS; tweak as needed
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},

		// DisableCompression: true, // only if you want raw bodies (usually keep compression)
	}
	return &http.Client{
		Transport: transport,

		// Global request timeout (hard cap). Prefer per-request context timeouts;
		// If you set this, it applies to the whole exchange (including body read).
		// Timeout: 20 * time.Second,

		// Optional: limit redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse // or return an error to stop earlier
			}
			return nil
		},
	}

}

func Shutdown(c *http.Client) {
	if tr, ok := c.Transport.(*http.Transport); ok {
		tr.CloseIdleConnections()
	}
}
