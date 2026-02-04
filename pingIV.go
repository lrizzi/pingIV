package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// Convert a Roman numerals string to an integer

func romanToInt(s string) (int, error) {
	// This is an historicall adaptation:
	// Originally the Roman numerals have no concept of numeber 0 you just do 1-1
	// in 725BC "nulla" from latin nothing was introduced N it's ne short version

	if s == "N" || s == "n" || s == "nulla" {
		return 0, nil
	}

	s = strings.ToUpper(s)

	romanMap := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	result := 0
	prevValue := 0

	// Process from right to left
	for i := len(s) - 1; i >= 0; i-- {
		value, ok := romanMap[s[i]]
		if !ok {
			return 0, fmt.Errorf("Invalid Roman numerals: %c", s[i])
		}

		if value < prevValue {
			result -= value
		} else {
			result += value
		}
		prevValue = value
	}

	// Validate range for IPv4 octet
	if result < 0 || result > 255 {
		return 0, fmt.Errorf("Roman numeral converts to %d, which is out of range (0-255)", result)
	}

	return result, nil
}

// This converts a Roman numeral IPv4 address to decimal notation
func romanIPv4ToDecimal(romanIP string) (string, error) {
	parts := strings.Split(romanIP, ".")

	if len(parts) < 1 || len(parts) > 4 {
		return "", fmt.Errorf("invalid IP address format: expected 1-4 octets separated by dots")
	}

	decimalParts := make([]string, len(parts))

	for i, part := range parts {
		num, err := romanToInt(part)
		if err != nil {
			return "", fmt.Errorf("error in octet %d (%s): %v", i+1, part, err)
		}
		decimalParts[i] = strconv.Itoa(num)
	}

	// This a workaround pad with zeros to support 725 BC format
	for len(decimalParts) < 4 {
		decimalParts = append(decimalParts, "0")
	}

	return strings.Join(decimalParts, "."), nil
}

func main() {
	// Define flags
	count := flag.Int("c", 4, "Number of ping packets to send")
	verbose := flag.Bool("v", false, "Verbose output (show conversion)")
	timeout := flag.Duration("t", 5*time.Second, "Timeout for each ping")
	interval := flag.Duration("i", 1*time.Second, "Interval between pings")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS] <Roman IP address>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "This ping utility accepts Roman numerals IPv4 addresses.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s CXXVII.N.N.I                  # Ping 127.0.0.1\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s CXXVII...I                    # Ping 127.0.0.1 (725 BC format)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s CXXVII.nulla.nulla.I          # Ping 127.0.0.1 (725 BC latin format)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -c 10 VIII.VIII.VIII.VIII     # Send 10 pings to 8.8.8.8\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -v CXCII.CLXVIII.I.I          # Verbose mode (show the coversion on top)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFor the compatiblity notation after 725 BC will be automatically padded with zero/s\n")
		fmt.Fprintf(os.Stderr, "\nCC-BY-SA-4.0 | Leonardo Rizzi (XenT)\n")

	}

	flag.Parse()

	// Get remaining arguments
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	romanIP := args[0]

	// Convert Roman IP to decimal
	decimalIP, err := romanIPv4ToDecimal(romanIP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting Roman IP address: %v\n", err)
		os.Exit(1)
	}

	// Show conversion if verbose
	if *verbose {
		fmt.Printf("Converting: %s -> %s\n", romanIP, decimalIP)
		fmt.Println()
	}

	// Create pinger
	pinger, err := probing.NewPinger(decimalIP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating pinger: %v\n", err)
		os.Exit(1)
	}

	// Configure pinger
	pinger.Count = *count
	pinger.Timeout = *timeout
	pinger.Interval = *interval

	// Set up statistics handler
	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
			pkt.Nbytes, romanIP, pkt.Seq, pkt.Rtt, pkt.TTL)
	}

	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, romanIP, pkt.Seq, pkt.Rtt, pkt.TTL)
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", romanIP)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	// Print header
	fmt.Printf("PING %s (%s):\n", romanIP, decimalIP)

	// Run the pinger
	err = pinger.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running ping: %v\n", err)
		os.Exit(1)
	}

	// Exit with error code if some packets get lost
	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		os.Exit(1)
	}
}
