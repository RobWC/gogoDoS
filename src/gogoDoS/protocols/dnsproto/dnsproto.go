package dnsproto

import (
	"github.com/miekg/dns"
	"gogoDoS/channels/channelman"
	"gogoDoS/protocols/dnsproto/dnsstat"
	"log"
	"net"
	"sync"
)

func DnsQuery(wg *sync.WaitGroup, config *dns.ClientConfig, queryCounter int, cm *channelman.ChannelMan) {
	defer wg.Done()
	dnsClient := new(dns.Client)
	message := new(dns.Msg)
	message.SetQuestion(dns.Fqdn("cnn.com"), dns.TypeA)
	message.RecursionDesired = true

	response, rtt, err := dnsClient.Exchange(message, net.JoinHostPort(config.Servers[0], "53"))

	if err != nil {
		log.Println(err)
		cm.RunChannel <- true
	}

	if response == nil {
	} else {
		if response.Rcode != dns.RcodeSuccess {
			log.Println(" query fail")
		}
		var stat = new(dnsstat.Info)
		stat.Rtt = rtt
		cm.StatsChannel <- stat
		cm.RunChannel <- true
	}
}
