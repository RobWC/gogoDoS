package channelman

import (
	"gogoDoS/protocols/dnsproto/dnsstat"
)

type ChannelMan struct {
	StatsChannel chan *dnsstat.Info
	RunChannel   chan bool
}

func (cm *ChannelMan) Init() {
	cm.StatsChannel = make(chan *dnsstat.Info,5000)
	cm.RunChannel = make(chan bool,5000)
}
