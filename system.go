package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

// https://stackoverflow.com/questions/41240761/check-if-ip-address-is-in-private-network-space
// Courtesy of https://stackoverflow.com/users/961810/brad-peabody and https://stackoverflow.com/users/98050/dougnukem
var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}
func validIP(strIP string) bool {
	return net.ParseIP(strIP) != nil
}

func isPrivateIP(strIP string) bool {
	ip := net.ParseIP(strIP)

	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// IP is a struct of string IP address
type IP struct {
	Query string
}

// PubIp is a cached copy of the public IP for multiple 'me' refs in the input list
var PubIp string

// GetPublicIpV4 is a method used to query an output server to get the true public  IPv4 of client. Seems to be the most reliable method
func GetPublicIpV4() string {
	if len(PubIp) == 0 {
		req, err := http.Get("https://checkip.amazonaws.com")
		if err != nil {
			return err.Error()
		}
		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err.Error()
		}
		PubIp = strings.TrimSpace(string(body))
	}
	return PubIp
}

// GetOutboundPrivateIP is a method to get preferred outbound ip of this machine
func GetOutboundPrivateIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
