package dnsstat

import (
	"time"
)

type Info struct {
	Rtt time.Duration
}

type Stats struct {
	InfoCollection []Info
}
