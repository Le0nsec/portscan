package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type IPList struct {
	A int
	B []int
	C []int
	D []int
}

var cidrIPList IPList

func ParsePort(portString string) ([]int, error) {

	var portList []int

	pair := strings.Split(portString, ",")
	for _, item := range pair {
		if strings.Contains(item, "-") {
			portRange := strings.Split(item, "-")
			if len(portRange) != 2 {
				return portList, fmt.Errorf("%s is invalid port range", portString)
			}
			start, _ := strconv.Atoi(portRange[0])
			end, _ := strconv.Atoi(portRange[1])
			for i := start; i <= end; i++ {
				portList = append(portList, i)
			}
		} else {
			item, _ := strconv.Atoi(item)
			portList = append(portList, item)
		}
	}

	sort.Ints(portList)
	return portList, nil
}

func ParseIP(ipString string) ([]string, error) {
	ipList := []string{}

	pair := strings.Split(ipString, ",")
	for _, item := range pair {

		if net.ParseIP(item) != nil {
			ipList = append(ipList, item)
		} else if ip, network, err := net.ParseCIDR(item); err == nil {
			n, _ := network.Mask.Size()
			ipSub := strings.Split(ip.Mask(network.Mask).String(), ".")
			cidrIPList.A, _ = strconv.Atoi(ipSub[0])

			if n >= 24 {
				a, _ := strconv.Atoi(ipSub[1])
				cidrIPList.B = append(cidrIPList.B, a)
				a, _ = strconv.Atoi(ipSub[2])
				cidrIPList.C = append(cidrIPList.C, a)
				for i := 1; i < IPRange(n, 32); i++ {
					cidrIPList.D = append(cidrIPList.D, i)
				}
			} else if n >= 16 && n < 24 {
				a, _ := strconv.Atoi(ipSub[1])
				cidrIPList.B = append(cidrIPList.B, a)
				for i := 0; i < IPRange(n, 24); i++ {
					cidrIPList.C = append(cidrIPList.C, i)
				}
				for i := 1; i < 256; i++ {
					cidrIPList.D = append(cidrIPList.D, i)
				}
			} else if n >= 8 && n < 16 {
				for i := 0; i < IPRange(n, 16); i++ {
					cidrIPList.B = append(cidrIPList.B, i)
				}
				for i := 0; i < 256; i++ {
					cidrIPList.C = append(cidrIPList.C, i)
				}
				for i := 1; i < 256; i++ {
					cidrIPList.D = append(cidrIPList.D, i)
				}
			} else {
				return ipList, fmt.Errorf("%s is not supported", item)
			}
			return ipList, fmt.Errorf("CIDR")
		} else if strings.Contains(item, "-") {
			splitIP := strings.SplitN(item, "-", 2)
			ip := net.ParseIP(splitIP[0])
			endIP := net.ParseIP(splitIP[1])
			if endIP != nil {
				if !isStartingIPLower(ip, endIP) {
					return ipList, fmt.Errorf("%s is greater than %s", ip.String(), endIP.String())
				}
				ipList = append(ipList, ip.String())
				for !ip.Equal(endIP) {
					increaseIP(ip)
					ipList = append(ipList, ip.String())
				}
			} else {
				ipOct := strings.SplitN(ip.String(), ".", 4)
				endIP := net.ParseIP(ipOct[0] + "." + ipOct[1] + "." + ipOct[2] + "." + splitIP[1])
				if endIP != nil {
					if !isStartingIPLower(ip, endIP) {
						return ipList, fmt.Errorf("%s is greater than %s", ip.String(), endIP.String())
					}
					ipList = append(ipList, ip.String())
					for !ip.Equal(endIP) {
						increaseIP(ip)
						ipList = append(ipList, ip.String())
					}
				} else {
					return ipList, fmt.Errorf("%s is not an IP Address or CIDR Network", item)
				}
			}
		} else {
			return ipList, fmt.Errorf("%s is not an IP Address or CIDR Network", item)
		}
	}
	return ipList, nil
}

/*
// LinesToIPList processes a list of IP addresses or networks in CIDR format.
// Returning a list of all possible IP addresses.
func LinesToIPList(lines []string) ([]string, error) {
	ipList := []string{}
	for _, line := range lines {
		_ipList, err := ParseIP(line)
		if err != nil {
			return _ipList, fmt.Errorf("%s is not an IP Address", line)
		}
		for _, line := range _ipList {
			ipList = append(ipList, line)
		}
	}
	return ipList, nil
}

func center(s string, width int) string {
	n := width - len(s)
	if n <= 0 {
		return s
	}
	half := n / 2
	if n%2 != 0 && width%2 != 0 {
		half = half + 1
	}
	return strings.Repeat(" ", half) + s + strings.Repeat(" ", (n-half))
}

func rjust(s string, width int) string {
	n := width - len(s)
	if n <= 0 {
		return s
	}
	return strings.Repeat(" ", n) + s
}
*/

func ljust(s string, width int) string {
	n := width - len(s)
	if n <= 0 {
		return s
	}
	return s + strings.Repeat(" ", n)
}

/*
func isValidIPV4(ip string) bool {
	b := net.ParseIP(ip)
	if b.To4() == nil {
		return false
	}
	return true
}
*/

// increases an IP by a single address.
func increaseIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func isStartingIPLower(start, end net.IP) bool {
	if len(start) != len(end) {
		return false
	}
	for i := range start {
		if start[i] > end[i] {
			return false
		}
	}
	return true
}

// ReadFileLines returns all the lines in a file.
func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func IPRange(n int, m int) int {
	return int(math.Pow(2, float64(m-n)))
}

func removeDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func ipListSort(list []string) []string {
	var sortList1, sortList2, sortList3 []string
	for _, i := range list {
		re, _ := regexp.Compile(`(\d+)`)
		rep := re.ReplaceAllString(i, "00$1")
		sortList1 = append(sortList1, rep)
	}
	for _, i := range sortList1 {
		re, _ := regexp.Compile(`0*(\d{3})`)
		rep := re.ReplaceAllString(i, "$1")
		sortList2 = append(sortList2, rep)
	}
	sort.Strings(sortList2)
	for _, i := range sortList2 {
		re, _ := regexp.Compile(`0*(\d+)`)
		rep := re.ReplaceAllString(i, "$1")
		sortList3 = append(sortList3, rep)
	}
	return sortList3
}
