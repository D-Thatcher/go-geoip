package main

import (
	"bufio"
	"net"
	"os"
	"strings"
)

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

func OpenNodesAsMap(fpth string) map[string]struct{} {
	lines, err := readLines(fpth)
	if err != nil {
		panic(err)
	}

	m := make(map[string]struct{}, len(lines)) // 0 byte allocate

	for i := 0; i < len(lines); i++ {
		var s struct{}
		m[lines[i]] = s
	}
	return m
}

func IsExitNode(m map[string]struct{}, node string) bool {
	_, contained := m[node]
	return contained
}

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
