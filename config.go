package main

import (
	"log"
	"net"
	"strings"
)

type ConfigProtocol struct {
	Name string //Name of protocol
	Port uint16 //Port for protocol
	ID   uint16 //Protocol ID
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
	Query     string         //host name to query
}

//Return a new empty configuration
func NewConfig() *Config {
	return &Config{
		Protocol:  "",
		DstPort:   0,
		DstIPs:    make([]string, 0),
		SrcIPs:    make([]string, 0),
		Rate:      0,
		Duration:  0,
		Interface: new(net.Interface),
		Query:     ""}
}

//set the specified protocol
func (cfg *Config) SetProtocol(p string) {
	cfg.Protocol = p
}

//set the destination port
func (cfg *Config) SetDstPort(p uint16) {
	cfg.DstPort = p
}

//set the rate of number of requests per second
func (cfg *Config) SetRate(p uint) {
	cfg.Rate = p
}

//set the total number of seconds to run the test
func (cfg *Config) SetDuration(p uint) {
	cfg.Duration = p
}

func (cfg *Config) SetQuery(p string) {
	cfg.Query = p
}

//set the egress interface by name. If the Interface does not exist then barf and error and exit
func (cfg *Config) SetInterfaceByName(name string) error {
	newInt, err := net.InterfaceByName(name)
	cfg.Interface = newInt
	if err != nil {
		log.Fatalf("Error attempting to resolve interface: %s\n", err)
	}
	return nil
}

//set the destination IPs by the input string we patse them into a slice
func (cfg *Config) SetDstIPsByString(iplist string) {
	cfg.DstIPs = parseIPsToSlice(iplist)
}

//set the source IPs by the input string we parse them into a slice
func (cfg *Config) SetSrcIPsByString(iplist string) {
	cfg.SrcIPs = parseIPsToSlice(iplist)
}

//private function to parse an IP list
func parseIPsToSlice(iplist string) []string {
	return strings.Split(iplist, ",")
}
