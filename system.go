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
func validIP(str_IP string) bool {
	return net.ParseIP(str_IP) != nil
}

func isPrivateIP(str_IP string) bool {
	ip := net.ParseIP(str_IP)

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

type IP struct {
	Query string
}

var PubIp string

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

// Get preferred outbound ip of this machine
func GetOutboundPrivateIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
