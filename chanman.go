package main

type ChanMan struct {
	StatsChan chan *Info
	RunChan   chan bool
}

func NewChanMan() *ChanMan {
	return &ChanMan{StatsChan: make(chan *Info, 5000),
		RunChan: make(chan bool, 5000)}
}
