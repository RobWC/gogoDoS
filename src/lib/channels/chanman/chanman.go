package chanman

import (
	"lib/protocols/dnsproto/dnsstat"
)

type ChanMan struct {
	StatsChan chan *dnsstat.Info
	RunChan   chan bool
}

func NewChanMan() *ChanMan {
	return &ChanMan{StatsChan: make(chan *dnsstat.Info, 5000),
		RunChan: make(chan bool, 5000)}
}
