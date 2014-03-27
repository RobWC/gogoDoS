package main

import (
	"code.google.com/p/go.net/ipv4"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	posTOS      = 1  // type-of-service
	posTotalLen = 2  // packet total length
	posID       = 4  // identification
	posFragOff  = 6  // fragment offset
	posTTL      = 8  // time-to-live
	posProtocol = 9  // next protocol
	posChecksum = 10 // checksum
	posSrc      = 12 // source address
	posDst      = 16 // destination address
)

const (
	posUDPSrcPort = 1 // Source Port
	posUDPDstPort = 3 // Dest Port
	posLen        = 4
)

type UDPHeader struct {
	SrcPort  uint16 //Source Port
	DstPort  uint16 //Destination Port
	Len      uint16 //length of header
	Checksum uint16 //checksum of header, optional set to 0 for ipv4
}

func (uh *UDPHeader) Marshal() ([]byte, error) {
	if uh == nil {
		//no object used sad face
	}
	//check the header len
	b := make([]byte, 8)
	srcportb := make([]byte, 2)
	dstportb := make([]byte, 2)
	lenb := make([]byte, 2)
	checksumb := make([]byte, 2)
	binary.BigEndian.PutUint16(srcportb, uh.SrcPort)
	binary.BigEndian.PutUint16(dstportb, uh.DstPort)
	binary.BigEndian.PutUint16(lenb, uh.Len)
	binary.BigEndian.PutUint16(checksumb, uh.Checksum)
	fmt.Println(srcportb, uh.SrcPort)
	copy(b[0:2], srcportb)
	copy(b[2:4], dstportb)
	copy(b[4:6], lenb)
	copy(b[6:8], checksumb)

	return b, nil
}

func (uh *UDPHeader) GenChecksum() error {
	var padd uint16
	padd = 0
	var checksum uint16
	checksum = 0
	var word16 uint16

	for i := 0; i < udp_len+padd; i = i + 2 {
		word16 = (())
	}

	return nil
}

func foo() {
	//create new raw packet connection

	//localAddress, _ := net.ResolveIPAddr("ip4", "10.0.1.100")
	//remoteAddress2, _ := net.ResolveIPAddr("ip4", "10.0.1.100")
	localAddress := net.IPv4(1, 2, 3, 4)
	remoteAddress := net.IPv4(10, 0, 1, 100)

	fmt.Println(net.Interfaces())

	//ip on mac, ip4:udp for linux
	con, err := net.ListenPacket("ip4:udp", "0.0.0.0")

	if err != nil {
		fmt.Println(err)
	}

	rawCon, err := ipv4.NewRawConn(con)

	if err != nil {
		println("ERROR")
		fmt.Println(err)
	}

	rawCon.SetTTL(128)

	headers := new(ipv4.Header)
	payload := make([]byte, 42)

	query := []byte{0x0d, 0x35, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x64, 0x61, 0x69, 0x73, 0x79, 0x06, 0x75, 0x62, 0x75, 0x6e, 0x74, 0x75, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}

	headers.Src = localAddress
	headers.Dst = remoteAddress
	headers.Protocol = 17
	headers.Len = 20
	headers.Version = 4
	headers.TTL = 128

	//UDP Header
	uh := new(UDPHeader)
	uh.Len = 42
	uh.SrcPort = 12345
	uh.DstPort = 53
	uh.Checksum = 0
	udpHead, _ := uh.Marshal()
	fmt.Println(udpHead)
	headers.TotalLen = 20 + len(query) + len(udpHead)

	copy(payload[0:8], udpHead)
	copy(payload[8:], query)

	cm := new(ipv4.ControlMessage)

	//cm.Src = net.IPv4(1, 2, 3, 4)
	cm.Dst = net.IPv4(10, 0, 1, 100)
	cm.TTL = 128
	cm.IfIndex = 3

	fmt.Println(headers)
	fmt.Println(cm)
	fmt.Println(payload)

	rawCon.WriteTo(headers, payload, cm)

	fmt.Println("sent")

}

func poo() {
	serverIP, _ := net.ResolveIPAddr("ip4", "10.0.1.100")
	//serverAddr, err := net.ResolveUDPAddr("udp4", "10.0.1.100:33333")
	//udpcon, err := net.DialUDP("udp", nil, serverAddr)
	con1, err := net.ListenPacket("ip4", "0.0.0.0")

	if err != nil {
		fmt.Println(err)
	}
	//udpcon.Write([]byte("food"))

	con := ipv4.NewPacketConn(con1)
	con.SetTTL(30)

	payload := make([]byte, 64)
	payload[0] = 1
	fmt.Println(net.Interfaces())

	cm := new(ipv4.ControlMessage)

	cm.Src = net.IPv4(10, 0, 1, 100)
	cm.Dst = net.IPv4(10, 0, 1, 100)
	cm.TTL = 128
	cm.IfIndex = 4

	fmt.Println(cm)
	con.WriteTo(payload, nil, serverIP)

}

func main() {
	foo()
}
