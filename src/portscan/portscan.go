package portscan

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"portscan/config"
	. "portscan/config"
	"portscan/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Run(ipList []string, portList []int, mode int, id int) {
	if OutputFile != "" {
		var err error
		F, err = os.OpenFile(OutputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		CheckError(err)
		defer F.Close()
	}
	if OutputDetailFile != "" {
		var err error
		F_detail, err = os.OpenFile(OutputDetailFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		CheckError(err)
		defer F_detail.Close()
	}
	if mode == 0 {
		hostList := []string{}

		if IsPing {
			for _, line := range ipList {
				hostList = append(hostList, strings.SplitN(line, ":", 2)[0])
			}
			hostList = removeDuplicatesAndEmpty(hostList)
			hostList = pingList(hostList)
			if len(hostList) == 0 {
				fmt.Println("no active host")
				return
			}
			openScanList := []string{}
			for _, line := range ipList {
				host := strings.SplitN(line, ":", 2)[0]
				_, exist := find(hostList, host)
				if exist {
					openScanList = append(openScanList, line)
				}
			}
			ipList = ipListSort(openScanList)
		}
		fmt.Println("\nscanning...")
		fmt.Printf("Number of scans: %d\n", len(ipList))
		for _, line := range ipList {
			Ch <- true
			Wg.Add(1)

			pair := strings.SplitN(line, ":", 2)
			host := pair[0]
			port, _ := strconv.Atoi(pair[1])
			go scan(host, port, id)
		}
		Wg.Wait()
	}

	if mode == 1 {
		if IsPing {
			ipList = pingList(ipList)
			ipList = ipListSort(ipList)
			fmt.Println(ipList)
			if len(ipList) == 0 {
				fmt.Println("no active host")
				return
			}
		}
		fmt.Println("\nscanning...")
		if len(portList) > 10 {
			Port = strconv.Itoa(portList[0]) + "," + strconv.Itoa(portList[1]) + "," + strconv.Itoa(portList[2]) + "..." + strconv.Itoa(portList[len(portList)-1])
		}
		if len(ipList) == 1 {
			fmt.Printf("host: %s\n", ipList[0])
			fmt.Printf("port: %s\n", Port)
		} else {
			fmt.Printf("host: %s - %s\n", ipList[0], ipList[len(ipList)-1])
			fmt.Printf("port: %s\n", Port)
		}

		for _, host := range ipList {
			for _, port := range portList {
				Ch <- true
				Wg.Add(1)
				go scan(host, port, id)
			}
		}
		Wg.Wait()
	}
}

func scan(host string, port int, id int) {
	open, str := isOpen(host, port)
	if open {
		fName := filepath.Join(config.DataPath, fmt.Sprintf("%d.txt", id))
		cacheFile, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		CheckError(err)
		defer cacheFile.Close()
		b, _ := json.Marshal(utils.CacheInfo{Host: host, Port: port, Str: str})
		cacheFile.WriteString(fmt.Sprintf("%s\n", string(b)))

		ip := fmt.Sprintf("%s:%d", host, port)
		fmt.Printf("%s open    %s\n", ljust(ip, 21), str)
		F.WriteString(fmt.Sprintf("%s:%d\r\n", host, port))
		F_detail.WriteString(fmt.Sprintf("%s open    %s\n", ljust(ip, 21), str))
		OpenList = append(OpenList, ip)
	}
	<-Ch
	Wg.Done()
}

func isOpen(host string, port int) (bool, string) {
	var msg [128]byte
	var str string
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(Timeout)*time.Millisecond)
	if err != nil {
		//fmt.Println(err)
		return false, str
	}
	conn.SetReadDeadline((time.Now().Add(time.Millisecond * time.Duration(Timeout))))
	_, err = conn.Read(msg[0:])
	if err != nil {
		//fmt.Println(err.Error())
		//fmt.Println(msg)
		isHttpService, str := fetchHttpDetails(host, port)
		if isHttpService {
			return true, str
		}
		return true, str
	}
	conn.Close()
	str = strings.Replace(string(msg[0:]), "\n", "", -1)

	return true, str
}

func fetchHttpDetails(host string, port int) (bool, string) {
	var str string
	// for _, line := range ReqHeaders {
	// 	pair := strings.SplitN(line, ":", 2)
	// 	if len(pair) == 2 {
	// 		k, v := pair[0], strings.Trim(pair[1], " ")
	// 		if strings.ToLower(k) == "host" {
	// 			ReqHost = v
	// 		}
	// 		headers[k] = v
	// 	}
	// }
	url := fmt.Sprintf("http://%s:%d%s", host, port, Path)
	if port == 443 {
		url = fmt.Sprintf("https://%s%s", host, Path)
	}
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Timeout:   time.Duration(Timeout) * time.Millisecond,
		Transport: tr,
	}
	if !Redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//fmt.Println(err)
		return false, str
	}
	req.Host = ReqHost
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println(err.Error())
		return false, str
	}
	defer resp.Body.Close()

	info := &HttpInfo{}
	info.Type = resp.Header.Get("Content-Type")
	info.Server = resp.Header.Get("Server")

	var body [1024]byte
	_, err = resp.Body.Read(body[0:])
	if err != nil {
		str = fmt.Sprintf("%s %s %s %s", resp.Proto, resp.Status, info.Server, info.Type)
		//fmt.Println(err.Error())
		return true, str
	}
	respBody := string(body[0:])
	r := regexp.MustCompile(`(?i)<title>\s?(.*?)\s?</title>`).FindStringSubmatch(respBody)
	if len(r) == 2 {
		info.Title = r[1]
	}
	str = fmt.Sprintf("%s %s %s %s %s", resp.Proto, resp.Status, info.Title, info.Server, info.Type)
	return true, str
}

func ScanByCli(host string, port string, id int) {
	portList, err := ParsePort(port)
	CheckError(err)
	ipList, err := ParseIP(host)
	// 利用error传CIDR判断为CIDR模式
	if err != nil && err.Error() == "CIDR" {
		a := CidrIPList.A
		for _, b := range CidrIPList.B {
			for _, c := range CidrIPList.C {
				s := []string{}
				for _, d := range CidrIPList.D {
					s = append(s, fmt.Sprintf("%d.%d.%d.%d", a, b, c, d))
				}
				Run(s, portList, 1, id)
			}
		}
	} else {
		CheckError(err)
		Run(ipList, portList, 1, id)
	}
	if config.Web {
		err := utils.SaveResult(id)
		if err != nil {
			log.Println("save result err:", err)
		}
	}
}
