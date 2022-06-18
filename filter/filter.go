package filter

import (
	"sort"

	"github.com/kurtinge/apastat/collector"
)

type SortingAndFilterOptions struct {
	ShowAllSlots bool
	SortBy       SortingField
}

type SortingField string

const (
	SortingFieldSrvSlot     SortingField = "srv"
	SortingFieldRequestTime SortingField = "request_time"
)

func FilterAndSortSlots(slots []collector.Slot, filterOptions SortingAndFilterOptions) []collector.Slot {
	var filteredSlots []collector.Slot
	for _, slot := range slots {
		if !filterOptions.ShowAllSlots && slot.Mode == collector.ServerModeOpenSlot {
			continue
		}
		filteredSlots = append(filteredSlots, slot)
	}

	if filterOptions.SortBy == SortingFieldRequestTime {
		sort.SliceStable(filteredSlots, func(i, j int) bool {
			return filteredSlots[i].SecondsSinceRequest > filteredSlots[j].SecondsSinceRequest
		})
	}

	return filteredSlots
}
