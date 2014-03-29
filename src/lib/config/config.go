package config

import (
	"log"
	"net"
	"strings"
)

type ConfigProtocol struct {
	Name string //Name of protocol
	Port uint16 //Port for protocol
	Id   uint16 //Protocol ID
}

//Struct to hold the config
type Config struct {
	Protocol  string         //specified protocol
	DstPort   uint16         //destination port
	DstIPs    []string       //destination ips
	SrcIPs    []string       //source ips
	Rate      uint           //packet rate
	Duration  uint           //duration of test in seconds
	Interface *net.Interface //Network interface for egress traffic
}

//Return a new empty configuration
func NewConfig() *Config {
	return &Config{Protocol: "",
		DstPort:   0,
		DstIPs:    make([]string, 0),
		SrcIPs:    make([]string, 0),
		Rate:      0,
		Duration:  0,
		Interface: new(net.Interface)}
}

//set the specified protocol
func (cfg *Config) SetProtocol(p string) {
	cfg.Protocol = p
}

//
func (cfg *Config) SetDstPort(p uint16) {
	cfg.DstPort = p
}

func (cfg *Config) SetRate(p uint) {
	cfg.Rate = p
}

func (cfg *Config) SetDuration(p uint) {
	cfg.Duration = p
}

func (cfg *Config) SetInterfaceByName(name string) error {
	newInt, err := net.InterfaceByName(name)
	cfg.Interface = newInt
	if err != nil {
		log.Fatalf("Error attempting to resolve interface: %s\n", err)
	}
	return nil
}

func (cfg *Config) SetDstIPsByString(iplist string) {
	cfg.DstIPs = parseIPsToSlice(iplist)
}

func (cfg *Config) SetSrcIPsByString(iplist string) {
	cfg.SrcIPs = parseIPsToSlice(iplist)
}

func parseIPsToSlice(iplist string) []string {
	return strings.Split(iplist, ",")
}
