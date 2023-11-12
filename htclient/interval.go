package htclient

import (
	"strconv"
)

var intervals = map[string]interval{
	"h": {
		unit:       "hour",
		lengthSecs: 60 * 60,
	},
	"d": {
		unit:       "day",
		lengthSecs: 60 * 60 * 24,
	},
	"w": {
		unit:       "week",
		lengthSecs: 60 * 60 * 24 * 7,
	},
	"m": {
		unit:       "month",
		lengthSecs: 60 * 60 * 24 * 30,
	},
	"y": {
		unit:       "year",
		lengthSecs: 60 * 60 * 24 * 365,
	},
}

type interval struct {
	unit       string
	lengthSecs int64
}

func IntervalUnits() []string {
	unitSlice := make([]string, len(intervals))

	i := 0
	for k := range intervals {
		unitSlice[i] = k
		i++
	}
	return unitSlice
}

func intervaltoSecs(s string) int64 {
	l, _ := strconv.Atoi(s[:len(s)-1])
	length := int64(l)
	unit := s[len(s)-1:]
	return length * intervals[unit].lengthSecs
}
