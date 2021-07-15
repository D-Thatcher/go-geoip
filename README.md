[![PkgGoDev](https://pkg.go.dev/badge/github.com/D-Thatcher/go-geoip)](https://pkg.go.dev/github.com/D-Thatcher/go-geoip)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/D-Thatcher/go-geoip)
[![Go Report Card](https://goreportcard.com/badge/github.com/D-Thatcher/go-geoip?update=no-cache)](https://goreportcard.com/badge/github.com/D-Thatcher/go-geoip?update=no-cache)

# go-geoip

### A Golang CLI tool for geolocating IP addresses and hosts and displaying them on a static street map. With support for identifying cloud provider addresses and TOR exit nodes


## Overview
```
go-geoip -o=mymap.png -ip=158.69.63.54,2605:8d80:648:eab7:290e:5c6f:f83:f189
```

![Connected AWS IP, TOR exit node and a Vancouver IP](https://github.com/D-Thatcher/go-geoip/blob/master/doc/assets/asset1.png)


### Supports:
* IPv4 & IPv4 address format
* CIDR blocks
* AWS address lookup
* TOR exit node lookup
* Domain resolution 


### Getting started:
* Install [Go](https://golang.org/doc/install)
* Download the free Geolite2 City and ASN MaxMind [databases](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data?lang=en)
* Run `go install github.com/D-Thatcher/go-geoip`


### Using the CLI tool

    Usage of go-geoip:
        -asn-mmdb    input path of Geolite2-ASN.mmdb (default "GeoDB/GeoLite2-ASN.mmdb")
        -city-mmdb   input path of Geolite2-City.mmdb (default "GeoDB/GeoLite2-City.mmdb")
        -connect     true/false indicating whether to connect all the markers (default true)
        -ip          a comma delimited string of IP addresses or domains (default to your public IPv4)
        -o           output path of static map image (default "map.png")
        -onlymarker  true/false indicating whether to plot just a marker instead of text (default false)

    Help Options:
      -h, --help                      Show this help message



## More examples

#### Example of:
* Domain resolution (google.com, slack.com)
* Referencing your own IP address by using the keyword `me`

```
go-geoip -o=asset2.png -ip=me,google.com,slack.com
```

![Connected Google IP, Slack IP and a Vancouver IP](https://github.com/D-Thatcher/go-geoip/blob/master/doc/assets/asset2.png)


#### Using just markers:

```
go-geoip -onlymarker=true -o=asset3.png -ip=52.93.126.199
```

![Connected AWS IP, TOR exit node and a Vancouver IP](https://github.com/D-Thatcher/go-geoip/blob/master/doc/assets/asset3.png)


#### Example of

* Address lookup in cloud provider CIDR block address space. Note we store the CIDR block in our `exit_node` data `2406:da70:4000::/40` and check if the IP `2406:da70:4000::1:3f51` is in the subnet
* Displaying both an AWS hosted zone, and a TOR exit node



```
go-geoip -o=asset4.png -ip=2406:da70:4000::1:3f51,199.249.230.186
```

![Connected AWS IPv6 CIDR block, TOR exit node](https://github.com/D-Thatcher/go-geoip/blob/master/doc/assets/asset4.png)



#### Example of

* Domain resolution (github.com)
* Referencing your own IPv4
* Custom Geolite2 database paths
* Not connecting the markers


```
go-geoip -o=asset5.png -connect=false -city-mmdb=GeoDB/GeoLite2-City.mmdb -asn-mmdb=GeoDB/GeoLite2-ASN.mmdb -ip=me,github.com
```

![Disconnected github.com and Vancouver IP](https://github.com/D-Thatcher/go-geoip/blob/master/doc/assets/asset5.png)

### Exit node data
* TOR exit nodes are listed in `exit_nodes/nodes.txt`. You can update the list by running 
    ```
    bash scripts/get_latest_exit_nodes.sh
    ```
* Cloud provider hosted address spaces are listed (currently only AWS IPv4 & IPv6 address blocks) `exit_nodes/aws_ip.txt`. You can update the list by running 
    ```
    bash scripts/get_cloud_host.sh
    ```


### Acknowledgements
* [OpenStreetMap](https://www.openstreetmap.org/copyright)
* [go-staticmaps](https://github.com/flopp/go-staticmaps)
* Free databases provided by [MaxMind](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data?lang=en)
* Parsing MaxMind databases with [geoip2-golang](https://github.com/oschwald/geoip2-golang)
