package awgproxy

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPServerServeGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok from upstream"))
	}))
	defer ts.Close()

	server := &HTTPServer{
		dial: func(network, address string) (net.Conn, error) {
			return net.Dial(network, ts.Listener.Addr().String())
		},
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		server.serve(conn)
	}()

	clientConn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = clientConn.Close() }()

	req, err := http.NewRequest(http.MethodGet, "http://"+ts.Listener.Addr().String()+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := req.Write(clientConn); err != nil {
		t.Fatal(err)
	}

	resp, err := http.ReadResponse(bufio.NewReader(clientConn), req)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	if string(body) != "ok from upstream" {
		t.Fatalf("unexpected body: %q", string(body))
	}
}

func TestHTTPServerServeFailureClosesConn(t *testing.T) {
	server := &HTTPServer{}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		server.serve(conn)
	}()

	clientConn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = clientConn.Close() }()

	// Send invalid HTTP request
	if _, err := clientConn.Write([]byte("INVALID\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1)
	_, err = clientConn.Read(buf)
	if err != io.EOF {
		t.Fatalf("expected EOF when server closes connection, got %v", err)
	}
}

// TestHTTPServePlainPostForwarded verifies that a plain (non-CONNECT) HTTP
// request such as POST is forwarded to the upstream peer with its method,
// target and body intact, instead of being rejected as an unsupported method
// (regression test for issue #27).
func TestHTTPServePlainPostForwarded(t *testing.T) {
	// upstream is the fake origin server the proxy dials into.
	upstreamClient, upstream := net.Pipe()

	s := &HTTPServer{
		config: &HTTPConfig{},
		dial: func(network, address string) (net.Conn, error) {
			if want := "example.com:80"; address != want {
				t.Errorf("dial address = %q, want %q", address, want)
			}
			return upstreamClient, nil
		},
	}

	// client talks to the proxy through an in-memory connection.
	clientConn, proxyConn := net.Pipe()

	go s.serve(proxyConn)

	// Client sends a plain HTTP POST with a body, proxy-style (absolute URI).
	go func() {
		_, _ = clientConn.Write([]byte(
			"POST http://example.com/api HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"Content-Length: 5\r\n" +
				"\r\n" +
				"hello"))
	}()

	_ = upstream.SetReadDeadline(time.Now().Add(5 * time.Second))
	forwarded, err := http.ReadRequest(bufio.NewReader(upstream))
	if err != nil {
		t.Fatalf("upstream did not receive forwarded request: %v", err)
	}

	if forwarded.Method != http.MethodPost {
		t.Errorf("forwarded method = %q, want POST", forwarded.Method)
	}
	if forwarded.URL.Path != "/api" {
		t.Errorf("forwarded path = %q, want /api", forwarded.URL.Path)
	}

	body := make([]byte, 5)
	if _, err := upstream.Read(body); err != nil {
		t.Fatalf("reading forwarded body: %v", err)
	}
	if got := string(body); got != "hello" {
		t.Errorf("forwarded body = %q, want %q", got, "hello")
	}

	_ = upstream.Close()
	_ = clientConn.Close()
}

// TestHTTPServeConnectStillTunnels verifies the CONNECT path is unchanged: the
// proxy answers "200 Connection established" and then tunnels bytes.
func TestHTTPServeConnectStillTunnels(t *testing.T) {
	upstreamClient, upstream := net.Pipe()

	s := &HTTPServer{
		config: &HTTPConfig{},
		dial: func(network, address string) (net.Conn, error) {
			if want := "example.com:443"; address != want {
				t.Errorf("dial address = %q, want %q", address, want)
			}
			return upstreamClient, nil
		},
	}

	clientConn, proxyConn := net.Pipe()
	go s.serve(proxyConn)

	go func() {
		_, _ = clientConn.Write([]byte(
			"CONNECT example.com:443 HTTP/1.1\r\n" +
				"Host: example.com:443\r\n" +
				"\r\n"))
	}()

	_ = clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	resp, err := http.ReadResponse(bufio.NewReader(clientConn), nil)
	if err != nil {
		t.Fatalf("reading CONNECT response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("CONNECT status = %d, want 200", resp.StatusCode)
	}

	_ = upstream.Close()
	_ = clientConn.Close()
}

func TestHTTPServePlainPostNotRejected(t *testing.T) {
	// Guard: a POST must not produce a 405 Method Not Allowed back to the client.
	dialed := make(chan struct{}, 1)
	s := &HTTPServer{
		config: &HTTPConfig{},
		dial: func(network, address string) (net.Conn, error) {
			dialed <- struct{}{}
			c, _ := net.Pipe()
			return c, nil
		},
	}

	clientConn, proxyConn := net.Pipe()
	go s.serve(proxyConn)

	go func() {
		_, _ = clientConn.Write([]byte(
			"POST http://example.com/ HTTP/1.1\r\nHost: example.com\r\nContent-Length: 0\r\n\r\n"))
	}()

	select {
	case <-dialed:
		// Reaching dial proves the request was routed to handle(), not rejected.
	case <-time.After(5 * time.Second):
		t.Fatal("POST was not forwarded to upstream (likely rejected as unsupported method)")
	}

	// Sanity: ensure no 405 response leaked back to the client before dial.
	_ = clientConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	buf := make([]byte, 64)
	if n, _ := clientConn.Read(buf); n > 0 {
		if strings.Contains(string(buf[:n]), "405") {
			t.Errorf("client received a 405 response: %q", string(buf[:n]))
		}
	}

	_ = clientConn.Close()
}
