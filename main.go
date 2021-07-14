package main

import (
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"strings"
)

const Default_output_file = "map.png"
const Default_city_db = "GeoDB/GeoLite2-City.mmdb"
const Default_asn_db = "GeoDB/GeoLite2-ASN.mmdb"
const exit_node_file = "./exit_nodes/nodes.txt"
const aws_exit_node_file = "./exit_nodes/aws_ip.txt"

var wordPtr = flag.String("ip", "", "a comma delimited string of IPv4 addresses (default to your public IP)")
var markerPtr = flag.Bool("onlymarker", false, "true/false indicating whether to plot just a marker instead of text")
var connPtr = flag.Bool("connect", true, "true/false indicating whether to connect all the markers")
var outFilePtr = flag.String("o", Default_output_file, "output path of static map image")
var maxMindCityFilePtr = flag.String("city-mmdb", Default_city_db, "input path of Geolite2-City.mmdb")
var maxMindASNFilePtr = flag.String("asn-mmdb", Default_asn_db, "input path of Geolite2-ASN.mmdb")

func findIP(db *geoip2.Reader, ipadd string) *geoip2.City {
	ip := net.ParseIP(ipadd)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	return record
}

func OpenMMDB(pth string) *geoip2.Reader {
	db, err := geoip2.Open(pth)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func findASN(db *geoip2.Reader, ipadd string) *geoip2.ASN {
	ip := net.ParseIP(ipadd)
	record, err := db.ASN(ip)
	if err != nil {
		log.Fatal(err)
	}
	return record
}

func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

type Displayable struct {
	City         *geoip2.City
	Asn          *geoip2.ASN
	Exit_node    bool
	Service_type string
	IP_address   string
	Connected_to *[]*Displayable
	Pin_marker   bool
}

func main() {
	flag.Parse()
	Arg_ip := strings.Split(*wordPtr, ",")

	if len(*wordPtr) == 0 {
		Arg_ip = []string{GetPublicIpV4()}
	}

	var allDisplay []*Displayable
	ASN_db := OpenMMDB(*maxMindASNFilePtr)
	defer ASN_db.Close()
	City_db := OpenMMDB(*maxMindCityFilePtr)
	defer City_db.Close()
	exit_node_map := OpenNodesAsMap(exit_node_file)
	aws_ip_lines, err := readLines(aws_exit_node_file)
	if err != nil {
		panic(err)
	}

	var CleanIP []string

	for _, Ip_address := range Arg_ip {
		Ip_address = strings.ToLower(Ip_address)
		if Ip_address == "me" {
			Ip_address = GetPublicIpV4()
		}
		if strings.Contains(Ip_address, "http") {
			panic("Remove protocol from hostnames " + Ip_address)
		}
		if strings.Contains(Ip_address, "/") {
			if strings.Contains(Ip_address, ":") {
				panic("IPv6 Address space for CIDR block is too expansive to locate")
			}
			fmt.Printf("Warning: removing CIDR block from ip notation %s \n", Ip_address)
			IP_sep := strings.Split(Ip_address, "/")
			Ip_address = IP_sep[0]
		}
		if !validIP(Ip_address) {
			addr, err := net.LookupIP(Ip_address)
			if err != nil {
				panic("Unknown IP address or hostname format " + Ip_address)
			}
			Ip_address = addr[0].String()
			if len(addr) > 1 {
				fmt.Printf("Warning: using only one of the host's ip address even though %v were available \n", len(addr))
			}
		}

		if isPrivateIP(Ip_address) {
			panic("Unable to geolocate private IP address " + Ip_address)
		}
		CleanIP = append(CleanIP, Ip_address)
	}

	IP_belongs_to_aws := IsAWSExitNode(&aws_ip_lines, &Arg_ip)
	for idx_raw, Clean_Ip_address := range CleanIP {
		var rawinput string = Arg_ip[idx_raw]
		if rawinput == "me" {
			rawinput = GetPublicIpV4() // it's cached, dw
		}
		var Is_exit bool = false
		var Service_type string = ""

		if !IsIPv6(Clean_Ip_address) {
			Is_exit = IsExitNode(exit_node_map, Clean_Ip_address)
			if Is_exit {
				Service_type = "TOR"
			}
		}
		if !Is_exit {
			v, _ := IP_belongs_to_aws[rawinput]
			Is_exit = v
			if Is_exit {
				Service_type = "AWS"
			}
		}

		ASNrecord := findASN(ASN_db, Clean_Ip_address)
		Cityrecord := findIP(City_db, Clean_Ip_address)

		var myslice *[]*Displayable

		dpl := Displayable{Cityrecord, ASNrecord, Is_exit, Service_type, Clean_Ip_address, myslice, false}
		allDisplay = append(allDisplay, &dpl)

	}
	if *connPtr {
		for _, d := range allDisplay {
			d.Connected_to = &allDisplay
		}
	}

	buildMap(allDisplay, outFilePtr, markerPtr)

}
