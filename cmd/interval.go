package main

import (
	"strconv"
)

var intervals = map[string]Interval{
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

type Interval struct {
	Unit       string
	LengthSecs int64
}

func intervaltoSecs(s string) int64 {
	l, _ := strconv.Atoi(s[:len(s)-1])
	length := int64(l)
	unit := s[len(s)-1:]
	return length * intervals[unit].LengthSecs
}

func getIntervalUnits(m map[string]Interval) []string {
	unitSlice := make([]string, len(m))

	i := 0
	for k := range m {
		unitSlice[i] = k
		i++
	}
	return unitSlice
}
