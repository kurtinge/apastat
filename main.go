package main

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kurtinge/apastat/collector"
)

func main() {
	apache := collector.NewApacheCollector()

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

	serverStatus, err := apache.GetStats()
	if err != nil {
		fmt.Printf("Unable to fetch serverstatus. Error: %s", err)
		return
	}

	for _, slot := range serverStatus.ServerSlots {
		if slot.Mode == collector.ServerModeWaiting {
			continue
		}
		if slot.Mode == collector.ServerModeOpenSlot {
			continue
		}
		// fmt.Printf(format, slot.ServerSlot, slot.Pid, slot.Mode, slot.SecondsSinceRequest, slot.Client, slot.Vhost, slot.Request)
		r := []*simpletable.Cell{
			{Text: slot.ServerSlot},
			{Text: fmt.Sprintf("%d", slot.Pid)},
			{Align: simpletable.AlignCenter, Text: string(slot.Mode)},
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
