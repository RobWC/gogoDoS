package dnsraw

import (
	"code.google.com/p/go.net/ipv4"
	"lib/channels/chanman"
	"lib/config"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

//RawDNS data struct. Contains L3, L4 Headers, payload and control message for specifying the egress interface
type RawDNS struct {
	IPHeaders     *ipv4.Header
	UDPHeader     *UDPHeader
	LocalAddress  net.IP
	RemoteAddress net.IP
	Payload       []byte
	CtrlMsg       *ipv4.ControlMessage
}

//return an initalized RawDNS struct
func NewRawDNS() *RawDNS {
	return &RawDNS{IPHeaders: new(ipv4.Header),
		UDPHeader:     new(UDPHeader),
		LocalAddress:  net.IPv4(0, 0, 0, 0),
		RemoteAddress: net.IPv4(0, 0, 0, 0),
		Payload:       make([]byte, 0),
		CtrlMsg:       new(ipv4.ControlMessage)}
}

//set destination port
func (rdns *RawDNS) SetDestPort(port uint16) {
	rdns.UDPHeader.SetDstPort(port)
}

//set local or source address
func (rdns *RawDNS) SetLocalAddress(ip string) {
	parsedIP := strings.Split(ip, ".")
	ip0, _ := strconv.Atoi(parsedIP[0])
	ip1, _ := strconv.Atoi(parsedIP[1])
	ip2, _ := strconv.Atoi(parsedIP[2])
	ip3, _ := strconv.Atoi(parsedIP[3])
	rdns.LocalAddress = net.IPv4(byte(ip0), byte(ip1), byte(ip2), byte(ip3))
}

//set remote address of
func (rdns *RawDNS) SetRemoteAddress(ip string) {
	parsedIP := strings.Split(ip, ".")
	ip0, _ := strconv.Atoi(parsedIP[0])
	ip1, _ := strconv.Atoi(parsedIP[1])
	ip2, _ := strconv.Atoi(parsedIP[2])
	ip3, _ := strconv.Atoi(parsedIP[3])
	rdns.RemoteAddress = net.IPv4(byte(ip0), byte(ip1), byte(ip2), byte(ip3))
	rdns.CtrlMsg.Dst = rdns.RemoteAddress
}

func (rdns *RawDNS) DnsQuery(wg *sync.WaitGroup, config *config.Config, cm *chanman.ChanMan) {
	defer wg.Done()

	//set the IP headers
	rdns.IPHeaders.Src = rdns.LocalAddress
	rdns.IPHeaders.Dst = rdns.RemoteAddress
	rdns.IPHeaders.Protocol = IPProtoUDP
	rdns.IPHeaders.Len = IPHeaderLen
	rdns.IPHeaders.Version = 4
	rdns.IPHeaders.TTL = 128

	//set the UDP headers
	rdns.UDPHeader.SetLen(UDPHeaderLen)
	rdns.UDPHeader.GenRandomSrcPort()
	rdns.UDPHeader.SetChecksum(0)
	udpHead, _ := rdns.UDPHeader.Marshal()

	//set the control message
	rdns.CtrlMsg.TTL = 128
	rdns.CtrlMsg.IfIndex = config.Interface.Index

	//set the query
	query := NewQuery()
	query.SetRequest(config.Query, "A")
	queryb := query.Marshal()

	//ip on mac, ip4:udp for linux
	con, err := net.ListenPacket("ip4:udp", "0.0.0.0")
	if err != nil {
		log.Fatalln(err)
	}

	//new raw packet connection
	rawCon, err := ipv4.NewRawConn(con)
	if err != nil {
		log.Fatalln(err)
	}

	//set query
	//query := []byte{0x0d, 0x35, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x64, 0x61, 0x69, 0x73, 0x79, 0x06, 0x75, 0x62, 0x75, 0x6e, 0x74, 0x75, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}

	//set final payload
	rdns.Payload = append(rdns.Payload, udpHead)
	rdns.Payload = append(rdns.Payload, queryb)

	//set packet length
	rdns.IPHeaders.TotalLen = 20 + len(queryb) + len(udpHead)

	rawCon.WriteTo(rdns.IPHeaders, rdns.Payload, rdns.CtrlMsg)
	cm.RunChan <- true
}
