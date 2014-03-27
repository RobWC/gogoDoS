package dnsraw

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
	//calcumete
	return nil
}

func (uh *UDPHeader) SetSrcPort(port uint16) {
	un.SrcPort = port
}

func (uh *UDPHeader) SetDstPort(port uint16) {
	un.DstPort = port
}

func (uh *UDPHeader) SetLen(l uint16) {
	un.Len = l
}

func (uh *UDPHeader) SetChecksum(cs uint16) {
	un.Checksum = cs
}
