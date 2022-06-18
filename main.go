package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/kurtinge/apastat/collector"
	"github.com/kurtinge/apastat/filter"
)

const VERSION = "0.1.0"

func main() {
	filterOptions := filter.SortingAndFilterOptions{
		ShowAllSlots: false,
		SortBy:       filter.SortingFieldSrvSlot,
	}

	host := flag.String("host", "localhost", "host to fetch server status from")
	sortByRequestTime := flag.Bool("s", false, "Sort by request time")
	showVersion := flag.Bool("v", false, "Show version and exit")
	flag.BoolVar(&filterOptions.ShowAllSlots, "a", false, "Show all slots")
	flag.Parse()

	if *showVersion {
		fmt.Printf("apastat %s\n", VERSION)
		os.Exit(0)
	}

	if *sortByRequestTime {
		filterOptions.SortBy = filter.SortingFieldRequestTime
	}

	tab := simpletable.New()
	tab.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: "Srv"},
			{Align: simpletable.AlignLeft, Text: "PID"},
			{Align: simpletable.AlignCenter, Text: "Mode"},
			{Align: simpletable.AlignRight, Text: "SS"},
			{Align: simpletable.AlignCenter, Text: "Client"},
			{Align: simpletable.AlignLeft, Text: "VHost"},
			{Align: simpletable.AlignLeft, Text: "Request"},
		},
	}

	apache := collector.NewApacheCollector(host)
	serverStatus, err := apache.GetStats()
	if err != nil {
		fmt.Printf("Unable to fetch serverstatus. Error: %s", err)
		return
	}

	serverStatusSlots := filter.FilterAndSortSlots(serverStatus.ServerSlots, filterOptions)
	for _, slot := range serverStatusSlots {
		r := []*simpletable.Cell{
			{Text: slot.ServerSlot},
			{Text: fmt.Sprintf("%d", slot.Pid)},
			{Align: simpletable.AlignCenter, Text: txtSlotMode(slot.Mode)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", slot.SecondsSinceRequest)},
			{Text: slot.Client},
			{Text: slot.Vhost},
			{Text: slot.Request},
		}
		tab.Body.Cells = append(tab.Body.Cells, r)
	}

	tab.SetStyle(simpletable.StyleCompactLite)
	tab.Println()
}

// fmt.Println(" _ = Waiting for connection")
// "_" Waiting for Connection, "S" Starting up, "R" Reading Request,
// "W" Sending Reply, "K" Keepalive (read), "D" DNS Lookup,
// "C" Closing connection, "L" Logging, "G" Gracefully finishing,
// "I" Idle cleanup of worker, "." Open slot with no current process

func txtSlotMode(mode collector.ServerMode) string {
	switch mode {
	case collector.ServerModeWaiting:
		return "Waiting"
	case collector.ServerModeStartingUp:
		return "Starting up"
	case collector.ServerModeReadingRequest:
		return "Reading"
	case collector.ServerModeSendingReply:
		return "Sending"
	case collector.ServerModeKeepalive:
		return "Keepalive"
	case collector.ServerModeDNSLookup:
		return "DNS Lookup"
	case collector.ServerModeClosingConnection:
		return "Closing"
	case collector.ServerModeLogging:
		return "Logging"
	case collector.ServerModeGracefullyFinishing:
		return "Gracefully finishing"
	case collector.ServerModeIdle:
		return "Idle"
	case collector.ServerModeOpenSlot:
		return "Open"
	}
	return string(mode)
}
