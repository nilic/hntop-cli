package htclient

import (
	"strconv"
)

var intervals = map[string]interval{
	"h": {
		Unit:       "hour",
		LengthSecs: 60 * 60,
	},
	"d": {
		Unit:       "day",
		LengthSecs: 60 * 60 * 24,
	},
	"w": {
		Unit:       "week",
		LengthSecs: 60 * 60 * 24 * 7,
	},
	"m": {
		Unit:       "month",
		LengthSecs: 60 * 60 * 24 * 30,
	},
	"y": {
		Unit:       "year",
		LengthSecs: 60 * 60 * 24 * 365,
	},
}

type interval struct {
	Unit       string
	LengthSecs int64
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
	return length * intervals[unit].LengthSecs
}
