package dnsraw

import (
	"encoding/binary"
	"math/rand"
	"time"
)

const (
	IPHeaderLen  = 20
	UDPHeaderLen = 42
	IPProtoUDP   = 17
	IPProtoTCP   = 6
	IPProtoICMP  = 1
)

//The UDP header structure is the header that will come after the IP header in an IP packet.
//It contains the Source Port, Destination Port, Length, and UDP Checksum
type UDPHeader struct {
	SrcPort  uint16 //Source Port
	DstPort  uint16 //Destination Port
	Len      uint16 //length of header
	Checksum uint16 //checksum of header, optional set to 0 for ipv4
}

//The marshal function packs the UDP header into a raw, network friendly (BigEndian), and returns it as a byte array.
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
	copy(b[0:2], srcportb)
	copy(b[2:4], dstportb)
	copy(b[4:6], lenb)
	copy(b[6:8], checksumb)

	return b, nil
}

//BUG(Rob): This doesnt work yet
//Generates the proper checksum for the UDP header
func (uh *UDPHeader) GenChecksum() error {
	//calcumete
	uh.Checksum = 0
	return nil
}

//Sets the source port of the packet
func (uh *UDPHeader) SetSrcPort(port uint16) {
	uh.SrcPort = port
}

//Sets the source port to a random port greater than 1025
func (uh *UDPHeader) GenRandomSrcPort() {
	rand.Seed(time.Now().UnixNano() * time.Now().UnixNano())
	uh.SrcPort = uint16(rand.Intn(65535-1025) + 1025)
}

func (uh *UDPHeader) SetDstPort(port uint16) {
	uh.DstPort = port
}

func (uh *UDPHeader) SetLen(l uint16) {
	uh.Len = l
}

func (uh *UDPHeader) SetChecksum(cs uint16) {
	uh.Checksum = cs
}
