package main

import (
    "github.com/miekg/dns"
    "strings"
    "log"
    "flag"
    "net"
//    "time"
)

var protoFlag = flag.String("p","","Specify protocol to use for DoS")
var destPort = flag.Int("P",53,"Specify destination port for DoS")
var destIPs = flag.String("d","","Specify an single host or a list of hosts seperated by comma (example: 1.2.3.4 or 1.2.3.4,2.3.4.5)")
var rateFlag = flag.Int("r",1,"Specify the amount of protocol requests per second")
var durationFlag = flag.Int("D",60,"Specify the total duration of the test")

func main() {
    flag.Parse()
    log.Printf("Starting DoS against %s on Port %d via protocol %s at rate %d/s", *destIPs, *destPort, *protoFlag, *rateFlag)
    config := new(dns.ClientConfig)
    config.Servers = strings.Split(*destIPs,":")
    config.Port = string(*destPort)
    config.Ndots = 1
    config.Timeout = 5
    config.Attempts = 1

    for i := 0; i < *rateFlag; i++ {
        dnsClient := new(dns.Client)
        message := new(dns.Msg)
        message.SetQuestion(dns.Fqdn("cnn.com"), dns.TypeA)
        message.RecursionDesired = true

        response, _, err := dnsClient.Exchange(message, net.JoinHostPort(config.Servers[0], "53"))

        if err != nil {
            log.Fatalln(err)
        }

        if response == nil {
            log.Println(" *** invalid ***")
        } else {
            if response.Rcode != dns.RcodeSuccess {
                log.Println(" query fail")
            }
            log.Println(response.MsgHdr.Id)
        }
    }

}
