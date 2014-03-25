package channelman

import (
	"gogoDoS/protocols/dnsproto/dnsstat"
)

type ChannelMan struct {
	StatsChannel chan *dnsstat.Info
	RunChannel   chan bool
}

func (cm *ChannelMan) Init() {
	cm.StatsChannel = make(chan *dnsstat.Info)
	cm.RunChannel = make(chan bool)
}
