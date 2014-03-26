package main

import (
	"code.google.com/p/go.net/ipv4"
	"fmt"
	"net"
)

func foo() {
	//create new raw packet connection

	//localAddress, _ := net.ResolveIPAddr("ip4", "10.0.1.100")
	//remoteAddress2, _ := net.ResolveIPAddr("ip4", "10.0.1.100")
	localAddress := net.IPv4(10, 0, 1, 100)
	remoteAddress := net.IPv4(10, 0, 1, 100)

	fmt.Println(net.InterfaceByName("wlan0"))

	con, err := net.ListenPacket("ip4:udp", "0.0.0.0")

	if err != nil {
		fmt.Println(err)
	}

	rawCon, err := ipv4.NewRawConn(con)

	if err != nil {
		println("ERROR")
		fmt.Println(err)
	}

	rawCon.SetTTL(30)

	headers := new(ipv4.Header)
	payload := make([]byte, 64)

	headers.Src = localAddress
	headers.Dst = remoteAddress
	headers.Protocol = 17
	headers.Len = 20
	headers.Version = 4
	headers.TotalLen = 84

	payload[0] = 1

	fmt.Println(headers)

	rawCon.WriteTo(headers, payload, nil)

}

func main() {
	//serverIP, _ := net.ResolveIPAddr("ip", "10.0.1.100")
	serverAddr, err := net.ResolveUDPAddr("udp4", "10.0.1.100:33333")
	udpcon, err := net.DialUDP("udp", nil, serverAddr)
	//con1, err := net.ListenPacket("ip4:udp", "0.0.0.0")

	if err != nil {
		fmt.Println(err)
	}
	//udpcon.Write([]byte("food"))

	con := ipv4.NewPacketConn(udpcon)
	con.SetTTL(30)

	payload := make([]byte, 64)
	payload[0] = 1
	fmt.Println(net.InterfaceByName("wlan0"))

	cm := new(ipv4.ControlMessage)

	cm.Src = net.IPv4(10, 0, 1, 100)
	cm.Dst = net.IPv4(10, 0, 1, 100)
	cm.TTL = 128
	cm.IfIndex = 3

	fmt.Println(cm)
	con.WriteTo(payload, cm, serverAddr)
}
