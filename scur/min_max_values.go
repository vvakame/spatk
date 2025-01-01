package scur

import (
	"strings"
	"time"
)

func StringMinValue() string {
	return ""
}

func StringMaxValue() string {
	return strings.Repeat("\uffff", 100)
}

func TimestampMinValue() time.Time {
	// TIMESTAMP("0001-01-01 00:00:00Z")
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TimestampMaxValue() time.Time {
	// TIMESTAMP("9999-12-31 23:59:59.999999999Z")
	return time.Date(10000, 1, 1, 0, 0, 0, -1, time.UTC)
}

func UUIDMinValue() string {
	// uuid.Nil
	return "00000000-0000-0000-0000-000000000000"
}

func UUIDMaxValue() string {
	// uuid.Max
	return "ffffffff-ffff-ffff-ffff-ffffffffffff"
}
