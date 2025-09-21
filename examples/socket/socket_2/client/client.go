package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	serverAddr  = "localhost:18000"
	dealTimeout = 5 * time.Second
)

func main() {
	interactiveMode := len(os.Args) > 1 && os.Args[1] == "-"

	conn, err := net.DialTimeout("tcp", serverAddr, dealTimeout)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	remoteAddr := conn.RemoteAddr().String()
	name := fmt.Sprintf("{{CLIENT=%s}}", localAddr)

	log.Printf("%s Connected to %s\n", name, remoteAddr)

	receiver := bufio.NewReader(conn)
	sender := bufio.NewWriter(conn)

	lines := []string{
		"hello\n",
		"ping\n",
		"goodbye\n",
	}
	// Send default lines
	for _, l := range lines {
		if _, err := sender.WriteString(l); err != nil {
			log.Fatalf("%s send write error: %v", name, err)
		}
		if err := sender.Flush(); err != nil {
			log.Fatalf("%s send flush error: %v", name, err)
		}

		// Read sender response
		resp, err := receiver.ReadString('\n')
		if err != nil {
			log.Fatalf("%s read error: %v", name, err)
		}
		log.Printf("%s Response: %s\n", name, resp)
	}

	if !interactiveMode {
		log.Printf("%s -- Interactive mode disabled, closing...\n", name)
		return
	}

	log.Printf("%s -- Interactive mode Enabled... \n", name)
	log.Println(">>> (type and press Enter, Ctrl+C to exit):")

	stdin := bufio.NewScanner(os.Stdin)
	fmt.Print(">>> ")
	for stdin.Scan() {

		text := stdin.Text() + "\n"
		if _, err := sender.WriteString(text); err != nil {
			log.Fatalf("%s write error: %v\n", name, err)
		}
		if err := sender.Flush(); err != nil {
			log.Fatalf("%s flush error: %v\n", name, err)
		}
		resp, err := receiver.ReadString('\n')
		if err != nil {
			log.Fatalf("%s read error: %v\n", name, err)
		}
		fmt.Printf("- {{SERVER}} %s", resp)
		fmt.Print(">>> ")
	}
}
