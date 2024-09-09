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
		fmt.Println("Usage: traceroute <IP_ADDRESS>")
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

	const maxHops = 30
	const timeout = time.Second * 3
	const packetSize = 52

	for ttl := 1; ttl <= maxHops; ttl++ {
		if err := conn.IPv4PacketConn().SetTTL(ttl); err != nil {
			fmt.Printf("Failed to set TTL: %v\n", err)
			return
		}

		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  ttl,
				Data: make([]byte, packetSize),
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

		conn.SetReadDeadline(time.Now().Add(timeout))

		reply := make([]byte, 1500)
		n, peer, err := conn.ReadFrom(reply)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("%d * * * Request timed out\n", ttl)
			continue
		}

		rm, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			fmt.Printf("Failed to parse ICMP reply: %v\n", err)
			return
		}

		switch rm.Type {
		case ipv4.ICMPTypeTimeExceeded:
			fmt.Printf("%d %s %.2fms\n", ttl, peer, float64(duration.Milliseconds()))
		case ipv4.ICMPTypeEchoReply:
			fmt.Printf("%d %s %.2fms (Reached destination)\n", ttl, peer, float64(duration.Milliseconds()))
			return
		default:
			fmt.Printf("%d * * * Unexpected reply\n", ttl)
		}
	}
}
