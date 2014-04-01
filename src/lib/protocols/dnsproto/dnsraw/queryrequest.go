package dnsraw

import (
	"encoding/binary"
	"strings"
)

const (
	ARecordType     = 0x0001
	NSRecordType    = 0x0002
	CNAMERecordType = 0x0005
	MXRecordType    = 0x000F
)

const (
	INClass = 0x0001
)

type QueryRequest struct {
	Name  string
	Type  uint16
	Class uint16
}

func NewQueryRequest() *QueryRequest {
	return &QueryRequest{Name: "",
		Type:  0,
		Class: 0}
}

func (qr *QueryRequest) SetType(s string) {
	s = strings.ToLower(s)
	switch {
	case s == "a":
		qr.Type = ARecordType
	case s == "ns":
		qr.Type = NSRecordType
	case s == "cname":
		qr.Type = CNAMERecordType
	case s == "mx":
		qr.Type = MXRecordType
	}
}

func (qr *QueryRequest) SetName(s string) {
	qr.Name = s
}

func (qr *QueryRequest) SetClass(i uint16) {
	qr.Class = i
}

func (qr *QueryRequest) SetClassDefault() {
	qr.Class = INClass
}

func (qr *QueryRequest) Marshal() []byte {
	//return byte array

	//break apart the name by period
	//count eatch segment
	nameBytes := make([]byte, 0)

	splitName := strings.Split(qr.Name, ".")
	for name := range splitName {
		nameBytes = append(nameBytes, byte(len(splitName[name])))
		nameBytes = append(nameBytes, []byte(splitName[name])...)
	}

	nameBytes = append(nameBytes, 0)
	//nameLen := len(nameBytes)
	b := make([]byte, 0)
	typeb := make([]byte, 2)
	classb := make([]byte, 2)
	binary.BigEndian.PutUint16(typeb, qr.Type)
	binary.BigEndian.PutUint16(classb, qr.Class)
	b = append(b, nameBytes...)
	b = append(b, typeb...)
	b = append(b, classb...)
	return b
}
