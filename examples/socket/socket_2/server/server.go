package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	addr        = "localhost:18000"
	readTimeout = 30 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addrTCP, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.ListenTCP("tcp", addrTCP)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	// Accept Loop connections
	go func(wg *sync.WaitGroup) {
		log.Printf("{{SERVER}} Server started listening on %s\n", addrTCP)

		for {
			conn, err := l.AcceptTCP()
			if err != nil {
				select {
				case <-ctx.Done():
					return // Ignore if during a shutdown signal
				default:
					log.Printf("{{SERVER}} Error accepting connection: %v\n", err)
					continue
				}
			}

			// Add Connection
			wg.Add(1)
			go func(conn *net.TCPConn) {
				defer wg.Done()
				handleConnection(ctx, conn)
			}(conn)

		}

	}(&wg)

	<-ctx.Done() // Block until we receive SIGINT or SIGTERM

	log.Println("Shutting down server...")
	_ = l.Close()

	log.Println("Waiting for connections to disconnect...")
	wg.Wait()
	log.Println("All connections are disconnected, closing...")

	log.Println("Server stopped")

}

func handleConnection(ctx context.Context, conn *net.TCPConn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr()
	log.Printf("{{SERVER}} {{REMOTE=%s}} Handle New connection from %s\n", clientAddr)
	// Optional: set TCP keepalives if available
	if err := conn.SetKeepAlive(true); err != nil {
		log.Printf("{{SERVER}} {{REMOTE=%s}} Error setting keepalive: %v\n", clientAddr, err)
	}
	if err := conn.SetKeepAlivePeriod(2 * time.Minute); err != nil {
		log.Printf("{{SERVER}} {{REMOTE=%s}} Error setting keepalive period: %v\n", clientAddr, err)
	}

	receiver := bufio.NewReader(conn)
	sender := bufio.NewWriter(conn)

	for {
		// Set a deadline for every read so idle clients don't hang forever.
		_ = conn.SetReadDeadline(time.Now().Add(readTimeout))

		line, err := receiver.ReadString('\n')
		if err != nil {
			// Most likely client closed the connection or deadline exceeded.
			// net.Error with Timeout() == true means deadline hit.
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				log.Printf("{{REMOTE=%s}} Read Timeout, error: %v", clientAddr, err)
			} else if errors.Is(err, io.EOF) {
				log.Printf("{{REMOTE=%s}} Connection closed", clientAddr)
			} else {
				log.Printf("{{REMOTE=%s}} Read Errorl, error: %v", clientAddr, err)
			}
			break
		}

		// Trim? For echo we keep the newline.
		log.Printf("{{REMOTE=%s}} recv: %q", clientAddr, line)

		resp := fmt.Sprintf("OK: %s", line)
		if _, err := sender.WriteString(resp); err != nil {
			log.Printf("{{REMOTE=%s}} write error: %v", clientAddr, err)
			return
		}
		if err := sender.Flush(); err != nil {
			log.Printf("{{REMOTE=%s}} flush error: %v", clientAddr, err)
			return
		}

		// Optional: stop processing if the server is shutting down
		select {
		case <-ctx.Done():
			return
		default:
		}

	}

	log.Printf("{{SERVER}} {{REMOTE=%s}} Closing connection...\n", clientAddr)

}
