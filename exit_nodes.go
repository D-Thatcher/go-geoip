package main

import (
	"bufio"
	"net"
	"os"
	"strings"
)

// readLines is a utility method for parsing the exit_node files
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// OpenNodesAsMap is a method used to read the nodes.txt file into an iterable (default 0 byte alloc)
func OpenNodesAsMap(fpth string) map[string]struct{} {
	lines, err := readLines(fpth)
	if err != nil {
		panic(err)
	}
	m := make(map[string]struct{}, len(lines))
	for i := 0; i < len(lines); i++ {
		var s struct{}
		m[lines[i]] = s
	}
	return m
}

// IsExitNode is a method used to determine if the input IP is in the Nodes map
func IsExitNode(m map[string]struct{}, node string) bool {
	_, contained := m[node]
	return contained
}

// IsAWSExitNode is a method that iterates once over the input IPs and determines whether they're in any AWS hosted subnet
func IsAWSExitNode(lines *[]string, LoIp *[]string) map[string]bool {
	Retmap := make(map[string]bool, len(*LoIp))

	for _, node := range *LoIp {
		if node == "me" {
			node = GetPublicIpV4() // it's cached, dw
		}
		var s []string
		var delim string
		if strings.Contains(node, ".") {
			s = strings.Split(node, ".")
			delim = "."
		} else if strings.Contains(node, ":") {
			s = strings.Split(node, ":")
			delim = ":"
		} else {
			panic("Unexpected input address " + node)
		}
		compStr := s[0] + delim + s[1]

		for _, line := range *lines {
			if line == node {
				Retmap[node] = true
				break
			} else if strings.Contains(line, compStr) {
				_, subnet, _ := net.ParseCIDR(line)
				ip := net.ParseIP(node)
				if subnet.Contains(ip) {
					Retmap[node] = true
					break
				}
			}
		}
	}
	return Retmap
}
