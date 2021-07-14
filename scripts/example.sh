
# ./go-geoip -connect=true -onlymarker=false -o=asset1.png -ip=64.252.108.0/24,me,2406:da70:4000::1:8264
# Since 2406:da70:4000::1:8264 is a member of the CIDR block subnet 2406:da70:4000::/40 listed in our exit nodes, it'll be labelled as such on the map. We also included an entire block 64.252.108.0/24 that will resolve to 64.252.108.0 when geolocating. Passing in an entire block of IPv6 is not available as it'll likely resolve imprecisely 

# Using an IP of me will make a network request to a server that will respond with your public IPv4 address. This will then be geolocated and displayed on the map 
# ./go-geoip -connect=true -onlymarker=false -o=asset1.png -ip=me

# TOR exit nodes are also labelled. You can update the list by running scripts/get_latest_exit_nodes.sh
# ./go-geoip -connect=true -onlymarker=false -o=asset1.png -ip=me,209.141.50.178,2406:da70:4000::1:8264

#asset1 ./go-geoip -o=mymap.png -ip=158.69.63.54,2605:8d80:648:eab7:290e:5c6f:f83:f189
#asset2 ./go-geoip -o=asset2.png -ip=me,google.com,slack.com
#asset3 ./go-geoip -onlymarker=true -o=asset3.png -ip=52.93.126.199
# ./go-geoip -connect=true -onlymarker=false -o=mymap.png -ip=205.206.163.40,2406:da70:4000::1:3bb2,193.218.118.182,2406:db70:4000::1:3f52,amazon.com
#asset4 ./go-geoip -o=asset4.png -ip=2406:da70:4000::1:3f51,199.249.230.186
#asset5 ./go-geoip -o=asset5.png -connect=false -city-mmdb=GeoDB/GeoLite2-City.mmdb -asn-mmdb=GeoDB/GeoLite2-ASN.mmdb -ip=me,github.com
./go-geoip -o=asset1.png -ip=158.69.63.54,2605:8d80:648:eab7:290e:5c6f:f83:f189
./go-geoip -o=asset2.png -ip=me,google.com,slack.com
./go-geoip -onlymarker=true -o=asset3.png -ip=52.93.126.199
./go-geoip -connect=true -onlymarker=false -o=mymap.png -ip=205.206.163.40,2406:da70:4000::1:3bb2,193.218.118.182,2406:db70:4000::1:3f52,amazon.com
./go-geoip -o=asset4.png -ip=2406:da70:4000::1:3f51,199.249.230.186
./go-geoip -o=asset5.png -connect=false -city-mmdb=GeoDB/GeoLite2-City.mmdb -asn-mmdb=GeoDB/GeoLite2-ASN.mmdb -ip=me,github.com

