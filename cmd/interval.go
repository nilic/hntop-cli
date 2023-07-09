package main

import (
	"fmt"
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

// Intervals{
// 	Intervals: []Interval{
// 		Interval{
// 			UnitShort:     "h",
// 			UnitLong:      "hour",
// 			LengthSeconds: 60 * 60,
// 		},
// 		Interval{
// 			UnitShort:     "d",
// 			UnitLong:      "day",
// 			LengthSeconds: 60 * 60 * 24,
// 		},
// 		Interval{
// 			UnitShort:     "w",
// 			UnitLong:      "week",
// 			LengthSeconds: 60 * 60 * 24 * 7,
// 		},
// 		Interval{
// 			UnitShort:     "m",
// 			UnitLong:      "month",
// 			LengthSeconds: 60 * 60 * 24 * 30,
// 		},
// 		Interval{
// 			UnitShort:     "y",
// 			UnitLong:      "year",
// 			LengthSeconds: 60 * 60 * 24 * 365,
// 		},
// 	},
// }

type Interval struct {
	Unit       string
	LengthSecs int64
}

// type Intervals struct {
// 	Intervals []Interval
// }

// func (i *Intervals) getShortUnits() []string {
// 	units := make([]string, len(i.Intervals))

// 	for c, v := range i.Intervals {
// 		units[c] = v.UnitShort
// 	}

// 	return units
// }

// func (i *Intervals) printUnits() string {
// 	var units string

// 	for _, v := range i.Intervals {
// 		units += fmt.Sprintf("\"%s\" (%s), ", v.UnitShort, v.UnitLong)
// 	}

// 	return units[:len(units)-2]
// }

func intervaltoSecs(s string) int64 {
	l, _ := strconv.Atoi(s[:len(s)-1])
	length := int64(l)
	unit := s[len(s)-1:]
	return length * intervals[unit].LengthSecs
	// switch unit := s[len(s)-1:]; unit {
	// case "h":
	// 	return length * intervals["hour"]
	// case "d":
	// 	return length * intervals["day"]
	// case "w":
	// 	return length * intervals["week"]
	// case "m":
	// 	return length * intervals["month"]
	// case "y":
	// 	return length * intervals["year"]
	// default:
	// 	return intervals["week"]
	// }

	// return intervals["week"]
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

func printUnits(m map[string]Interval) string {
	var units string

	for k, v := range m {
		units += fmt.Sprintf("\"%s\" - %s, ", k, v.Unit)
	}

	return units[:len(units)-2]
}
