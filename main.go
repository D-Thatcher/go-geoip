package main

import (
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"strings"
)

// DefaultOutputFile is the default static map output file location
const DefaultOutputFile = "map.png"

// DefaultCityDb is the default location of the IP to City MaxMind db
const DefaultCityDb = "GeoDB/GeoLite2-City.mmdb"

// DefaultAsnDb is the default location of the IP to ASN MaxMind db
const DefaultAsnDb = "GeoDB/GeoLite2-ASN.mmdb"

// ExitNodeFile is the default location of line-delimited text file containing TOR exit node addresses
const ExitNodeFile = "./exit_nodes/nodes.txt"

// AwsExitNodeFile is the default location of line-delimited text file containing AWS hosted zone subnets
const AwsExitNodeFile = "./exit_nodes/aws_ip.txt"

var wordPtr = flag.String("ip", "", "a comma delimited string of IPv4 addresses (default to your public IP)")
var markerPtr = flag.Bool("onlymarker", false, "true/false indicating whether to plot just a marker instead of text")
var connPtr = flag.Bool("connect", true, "true/false indicating whether to connect all the markers")
var outFilePtr = flag.String("o", DefaultOutputFile, "output path of static map image")
var maxMindCityFilePtr = flag.String("city-mmdb", DefaultCityDb, "input path of Geolite2-City.mmdb")
var maxMindASNFilePtr = flag.String("asn-mmdb", DefaultAsnDb, "input path of Geolite2-ASN.mmdb")

func findIP(db *geoip2.Reader, ipadd string) *geoip2.City {
	ip := net.ParseIP(ipadd)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	return record
}

// OpenMMDB is used to open an arbitrary MaxMind database (*.mmdb), leave the closing to the main func
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

// IsIPv6 is a prelim check to see if an IP might be of the IPv6 format
func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

// Displayable is a struct containing all the info we want to display on the street map. Connected_to lists the associated Displayable objs its linked to on the map
type Displayable struct {
	City        *geoip2.City
	Asn         *geoip2.ASN
	ExitNode    bool
	ServiceType string
	IPAddress   string
	ConnectedTo *[]*Displayable
	PinMarker   bool
}

func main() {
	flag.Parse()
	ArgIP := strings.Split(*wordPtr, ",")

	if len(*wordPtr) == 0 {
		ArgIP = []string{GetPublicIpV4()}
	}

	var allDisplay []*Displayable
	ASNDb := OpenMMDB(*maxMindASNFilePtr)
	defer ASNDb.Close()
	CityDb := OpenMMDB(*maxMindCityFilePtr)
	defer CityDb.Close()
	exitNodeMap := OpenNodesAsMap(ExitNodeFile)
	awsIPLines, err := readLines(AwsExitNodeFile)
	if err != nil {
		panic(err)
	}

	var CleanIP []string

	for _, IPAddressInput := range ArgIP {
		IPAddressInput = strings.ToLower(IPAddressInput)
		if IPAddressInput == "me" {
			IPAddressInput = GetPublicIpV4()
		}
		if strings.Contains(IPAddressInput, "http") {
			panic("Remove protocol from hostnames " + IPAddressInput)
		}
		if strings.Contains(IPAddressInput, "/") {
			if strings.Contains(IPAddressInput, ":") {
				panic("IPv6 Address space for CIDR block is too expansive to locate")
			}
			fmt.Printf("Warning: removing CIDR block from ip notation %s \n", IPAddressInput)
			IPSep := strings.Split(IPAddressInput, "/")
			IPAddressInput = IPSep[0]
		}
		if !validIP(IPAddressInput) {
			addr, err := net.LookupIP(IPAddressInput)
			if err != nil {
				panic("Unknown IP address or hostname format " + IPAddressInput)
			}
			IPAddressInput = addr[0].String()
			if len(addr) > 1 {
				fmt.Printf("Warning: using only one of the host's ip address even though %v were available \n", len(addr))
			}
		}

		if isPrivateIP(IPAddressInput) {
			panic("Unable to geolocate private IP address " + IPAddressInput)
		}
		CleanIP = append(CleanIP, IPAddressInput)
	}

	IPBelongsToAws := IsAWSExitNode(&awsIPLines, &ArgIP)
	for idxRaw, CleanIPAddress := range CleanIP {
		var rawinput string = ArgIP[idxRaw]
		if rawinput == "me" {
			rawinput = GetPublicIpV4() // it's cached, dw
		}
		var IsExit bool = false
		var ServiceType string = ""

		if !IsIPv6(CleanIPAddress) {
			IsExit = IsExitNode(exitNodeMap, CleanIPAddress)
			if IsExit {
				ServiceType = "TOR"
			}
		}
		if !IsExit {
			v, _ := IPBelongsToAws[rawinput]
			IsExit = v
			if IsExit {
				ServiceType = "AWS"
			}
		}

		ASNrecord := findASN(ASNDb, CleanIPAddress)
		Cityrecord := findIP(CityDb, CleanIPAddress)

		var myslice *[]*Displayable

		dpl := Displayable{Cityrecord, ASNrecord, IsExit, ServiceType, CleanIPAddress, myslice, false}
		allDisplay = append(allDisplay, &dpl)

	}
	if *connPtr {
		for _, d := range allDisplay {
			d.ConnectedTo = &allDisplay
		}
	}

	buildMap(allDisplay, outFilePtr, markerPtr)

}
