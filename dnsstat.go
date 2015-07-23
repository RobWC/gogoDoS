package main

import (
	"sort"
	"time"
)

type Info struct {
	Rtt time.Duration
}

type ByRtt []Info

func (this ByRtt) Len() int {
	return len(this)
}

func (this ByRtt) Less(i, j int) bool {
	return this[i].Rtt < this[j].Rtt
}

func (this ByRtt) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type Stats struct {
	InfoCollection []Info
}

//calculate min, max, avg, jitter
func (st *Stats) Calc() (min time.Duration, max time.Duration, avg time.Duration, jitter time.Duration) {

	//calculate average
	var avgCount time.Duration
	avgCount = 0
	var avgTotal time.Duration
	avgTotal = 0
	for info := range st.InfoCollection {
		avgCount = avgCount + 1
		avgTotal = avgTotal + st.InfoCollection[info].Rtt
	}
	avg = avgTotal / avgCount

	//calculate max
	sort.Sort(ByRtt(st.InfoCollection))
	max = st.InfoCollection[len(st.InfoCollection)-1].Rtt
	min = st.InfoCollection[0].Rtt
	jitter = max - min
	return min, max, avg, jitter
}
