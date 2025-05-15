// diameter_stats.go
package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// For Counters handling Starts
// PeerCounter holds the count and start time for a specific message type and peer.
type PeerCounter struct {
	count int
	//startTime time.Time
	count_old int
	//timestamps []time.Time
	mu sync.Mutex
}

// DiameterStats holds counters for each message type and peer IP.
type DiameterStats struct {
	counters map[string]*PeerCounter
	mu       sync.Mutex
}

// NewDiameterStats creates and initializes a new DiameterStats.
func NewDiameterStats() *DiameterStats {
	return &DiameterStats{
		counters: make(map[string]*PeerCounter),
	}
}

// IncrementReceived increments the counter for a specific message type and peer IP.
// don't send "_" parameter
func (ds *DiameterStats) IncrementReceived(msgName string, peerIP string, additionalIndex string) {
	//log.Printf("KEY: %s_%s_%s", msgName, peerIP, additionalIndex)
	//key := fmt.Sprintf("%s_%s_%s", msgName, peerIP, additionalIndex)
	key := fmt.Sprintf("%s_%s_%s", msgName, "PEERIP", additionalIndex)
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, ok := ds.counters[key]; !ok {
		ds.counters[key] = &PeerCounter{}
	}
	counter := ds.counters[key]
	counter.mu.Lock()
	defer counter.mu.Unlock()
	counter.count++
	//now := time.Now()
	//counter.timestamps = append(counter.timestamps, now)

}

// printMetrics prints the DIAMETER statistics every second in a column format with auto-adjusted column widths.
func printMetrics(stats *DiameterStats, ReportHeading string) {

	HeadingLen := len(ReportHeading) - 20
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("\033[H\033[2J") // Clear the terminal screen
		fmt.Printf("%+*s+\n", HeadingLen, strings.Repeat("+", HeadingLen))
		fmt.Printf("DIAMETER SERVER Statistics %s\n", now)
		fmt.Printf("%s\n", ReportHeading)
		fmt.Printf("%+*s+\n", HeadingLen, strings.Repeat("+", HeadingLen))

		// Determine maximum lengths for columns
		maxNameLen := len("Command-Name")
		maxPeerLen := len("Peer-Info")
		maxIndexLen := len("Additional-Index")
		maxCountLen := len("Total_Count")
		maxTPSLen := len("TPS")

		stats.mu.Lock()
		for key, counter := range stats.counters {
			parts := strings.Split(key, "_")
			if len(parts) == 3 {
				msgName := parts[0]
				peerInfo := parts[1]
				additionalIndex := parts[2]
				//tpsStr := fmt.Sprintf("%.2f", counter.calculateTPS(time.Second))
				tpsStr := fmt.Sprintf("%d", counter.count_old)
				countStr := fmt.Sprintf("%d", counter.count)

				if len(msgName) > maxNameLen {
					maxNameLen = len(msgName)
				}
				if len(peerInfo) > maxPeerLen {
					maxPeerLen = len(peerInfo)
				}
				if len(additionalIndex) > maxIndexLen {
					maxIndexLen = len(additionalIndex)
				}
				if len(countStr) > maxCountLen {
					maxCountLen = len(countStr)
				}
				if len(tpsStr) > maxTPSLen {
					maxTPSLen = len(tpsStr)
				}
			}
		}
		stats.mu.Unlock()

		// Adjust lengths for a bit of padding
		padding := 2
		nameLen := maxNameLen + padding
		peerLen := maxPeerLen + padding
		indexLen := maxIndexLen + padding
		countLen := maxCountLen + padding
		tpsLen := maxTPSLen + padding

		fmt.Printf("|%-*s-|-%-*s-|-%-*s-|-%-*s-|-%-*s|\n",
			nameLen, strings.Repeat("-", nameLen),
			peerLen, strings.Repeat("-", peerLen),
			indexLen, strings.Repeat("-", indexLen),
			countLen, strings.Repeat("-", countLen),
			tpsLen, strings.Repeat("-", tpsLen))

		fmt.Printf("|%-*s | %-*s | %-*s | %-*s | %-*s|\n",
			nameLen, "Command-Name",
			peerLen, "Peer-Info",
			indexLen, "Additional-Index",
			countLen, "Total_Count",
			tpsLen, "TPS")

		fmt.Printf("|%-*s-|-%-*s-|-%-*s-|-%-*s-|-%-*s|\n",
			nameLen, strings.Repeat("-", nameLen),
			peerLen, strings.Repeat("-", peerLen),
			indexLen, strings.Repeat("-", indexLen),
			countLen, strings.Repeat("-", countLen),
			tpsLen, strings.Repeat("-", tpsLen))

		stats.mu.Lock()

		// Sort keys for consistent order
		var sortedKeys []string
		for key := range stats.counters {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Strings(sortedKeys)

		for _, key := range sortedKeys {
			counter := stats.counters[key]
			parts := strings.Split(key, "_")
			if len(parts) == 3 {
				msgName := parts[0]
				peerInfo := parts[1]
				additionalIndex := parts[2]
				tps := int(counter.count - counter.count_old)
				//tps := counter.calculateTPS(time.Second)
				fmt.Printf("|%-*s | %-*s | %-*s | %-*d | %-*d|\n",
					nameLen, msgName,
					peerLen, peerInfo,
					indexLen, additionalIndex,
					countLen, counter.count,
					tpsLen, tps)
			}
			counter.count_old = counter.count
		}

		stats.mu.Unlock()

		fmt.Printf("|%-*s-|-%-*s-|-%-*s-|-%-*s-|-%-*s|\n",
			nameLen, strings.Repeat("-", nameLen),
			peerLen, strings.Repeat("-", peerLen),
			indexLen, strings.Repeat("-", indexLen),
			countLen, strings.Repeat("-", countLen),
			tpsLen, strings.Repeat("-", tpsLen))
	}
}

// For Counters handling Ends
