package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ping <IP_ADDRESS>")
		return
	}

	target := os.Args[1]

	ipAddr, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		fmt.Printf("Failed to resolve IP address: %v\n", err)
		return
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Printf("Failed to listen for ICMP packets: %v\n", err)
		return
	}
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("HELLO-PING"),
		},
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		fmt.Printf("Failed to marshal ICMP message: %v\n", err)
		return
	}

	start := time.Now()
	_, err = conn.WriteTo(msgBytes, ipAddr)
	if err != nil {
		fmt.Printf("Failed to send ICMP request: %v\n", err)
		return
	}

	reply := make([]byte, 1500)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second)) // Timeout after 3 seconds

	n, _, err := conn.ReadFrom(reply)
	if err != nil {
		fmt.Printf("Failed to receive ICMP reply: %v\n", err)
		return
	}
	duration := time.Since(start)

	rm, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		fmt.Printf("Failed to parse ICMP reply: %v\n", err)
		return
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("Ping to %s: seq=%d time=%v\n", target, 1, duration)
	default:
		fmt.Printf("Unexpected reply type: %v\n", rm)
	}
}
