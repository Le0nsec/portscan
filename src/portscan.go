package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	wg               sync.WaitGroup
	ch               chan bool
	host             string
	port             string
	loadFile         string
	timeout          int
	verbose          bool
	outputFile       string
	outputDetailFile string
	goroutineNum     int
	reqHost          string
	reqHeaders       arrayFlags
	redirect         bool
	path             string
	isPing           bool
	f                *os.File
	f_detail         *os.File
	openList         []string
	openHostList     []string
	//banner           bool
	//outputJSONFile   string
	//result           []HttpInfo
)

func main() {
	defaultPorts := "7,11,13,15,17,19,21,22,23,25,26,37,38,43,49,51,53,67,70,79,80,81,82,83,84,85,86,88,89,102,104,110,111,113,119,121,135,138,139,143,175,179,199,211,264,311,389,443,444,445,465,500,502,503,505,512,515,548,554,564,587,631,636,646,666,771,777,789,800,801,873,880,902,992,993,995,1000,1022,1023,1024,1025,1026,1027,1080,1099,1177,1194,1200,1201,1234,1241,1248,1260,1290,1311,1344,1400,1433,1471,1494,1505,1515,1521,1588,1720,1723,1741,1777,1863,1883,1911,1935,1962,1967,1991,2000,2001,2002,2020,2022,2030,2049,2080,2082,2083,2086,2087,2096,2121,2181,2222,2223,2252,2323,2332,2375,2376,2379,2401,2404,2424,2455,2480,2501,2601,2628,3000,3128,3260,3288,3299,3306,3307,3310,3333,3388,3389,3390,3460,3541,3542,3689,3690,3749,3780,4000,4022,4040,4063,4064,4369,4443,4444,4505,4506,4567,4664,4712,4730,4782,4786,4840,4848,4880,4911,4949,5000,5001,5002,5006,5007,5009,5050,5084,5222,5269,5357,5400,5432,5555,5560,5577,5601,5631,5672,5678,5800,5801,5900,5901,5902,5903,5938,5984,5985,5986,6000,6001,6068,6379,6488,6560,6565,6581,6588,6590,6664,6665,6666,6667,6668,6669,6998,7000,7001,7005,7014,7071,7077,7080,7288,7401,7443,7474,7493,7537,7547,7548,7634,7657,7777,7779,7911,8000,8001,8008,8009,8010,8020,8025,8030,8040,8060,8069,8080,8081,8082,8086,8087,8088,8089,8090,8098,8099,8112,8123,8125,8126,8139,8161,8200,8291,8333,8334,8377,8378,8443,8500,8545,8554,8649,8686,8800,8834,8880,8883,8888,8889,8983,9000,9001,9002,9003,9009,9010,9042,9051,9080,9090,9100,9151,9191,9200,9295,9333,9418,9443,9527,9530,9595,9653,9700,9711,9869,9944,9981,9999,10000,10001,10162,10243,10333,11001,11211,11300,11310,12300,12345,13579,14000,14147,14265,16010,16030,16992,16993,17000,18001,18081,18245,18246,19999,20000,20547,22105,22222,23023,23424,25000,25105,25565,27015,27017,28017,32400,33338,33890,37215,37777,41795,42873,45554,49151,49152,49153,49154,49155,50000,50050,50070,50100,51106,52869,55442,55553,60001,60010,60030,61613,61616,62078,64738"

	flag.StringVar(&host, "h", "", "scan `host`. format: 127.0.0.1 | 192.168.1.1/24 | 192.168.1.1-5")
	flag.StringVar(&port, "p", defaultPorts, "scan `port`. format: 1-65535 | 21,22,25 | 8080")
	flag.BoolVar(&isPing, "ping", false, "ping before scanning")
	flag.StringVar(&loadFile, "f", "", "load external `file`, ip:port are read by line")
	flag.IntVar(&timeout, "timeout", 4000, "connection `timeout` millisecond")
	flag.StringVar(&outputFile, "o", "", "save open ip:port per line `filepath`")
	flag.StringVar(&outputDetailFile, "O", "", "save details open ports `filepath`")
	flag.StringVar(&path, "path", "/", "request `urlpath`. example: /admin")
	flag.BoolVar(&redirect, "redirect", false, "follow 30x redirect")
	flag.Var(&reqHeaders, "H", "request `headers`. exmaple: -H User-Agent:xx -H Referer:xx")
	flag.IntVar(&goroutineNum, "t", 200, "scan max `threads`")
	flag.BoolVar(&verbose, "v", false, "show verbose")
	//flag.BoolVar(&banner, "banner", false, "whether to get the banner of the active port")
	//flag.StringVar(&outputJSONFile, "oj", "", "save banner json `filepath`, it required the -banner parameter")
	flag.Parse()

	var (
		scanList, ipList []string
		portList         []int
	)
	//限制goroutine数量
	ch = make(chan bool, goroutineNum)

	//host = "172.16.3.250"
	//port = "8080"

	if (host == "" && loadFile == "") || (host != "" && loadFile != "") {
		flag.Usage()
		os.Exit(0)
	}
	/*
		if outputJSONFile != "" && !banner {
			fmt.Println("outputJSONFile required the -banner parameter")
			os.Exit(1)
		}
	*/

	if loadFile != "" {
		lines, err := ReadFileLines(loadFile)
		checkError(err)

		if port != "" {
			portList, _ = ParsePort(port)
			for _, line := range lines {
				line = strings.Trim(line, " ")
				if strings.Contains(line, ":") {
					hostPort := strings.Split(line, ":")
					line = hostPort[0]
				}
				ipList = append(ipList, line)
			}
			run(ipList, portList, 1)
		} else {
			for _, line := range lines {
				line = strings.Trim(line, " ")
				if strings.Contains(line, ":") {
					scanList = append(scanList, line)
				} else {
					fmt.Println("loadfile contents format error")
					os.Exit(1)
				}
			}
			run(scanList, portList, 0)
		}
	} else {
		portList, err := ParsePort(port)
		checkError(err)
		ipList, err = ParseIP(host)
		// 利用error传CIDR判断为CIDR模式
		if err != nil && err.Error() == "CIDR" {
			a := cidrIPList.A
			for _, b := range cidrIPList.B {
				for _, c := range cidrIPList.C {
					s := []string{}
					for _, d := range cidrIPList.D {
						s = append(s, fmt.Sprintf("%d.%d.%d.%d", a, b, c, d))
					}
					run(s, portList, 1)
				}
			}
		} else {
			checkError(err)
			run(ipList, portList, 1)
		}
	}

	fmt.Println("portscan finish...")
	/*
		if banner {
			fmt.Println("\nbanner scanning...")
			bannerScan(openList)
			fmt.Println("bannerscan finish...")
		}
	*/
}

func run(ipList []string, portList []int, mode int) {
	if outputFile != "" {
		var err error
		f, err = os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		checkError(err)
		defer f.Close()
	}
	if outputDetailFile != "" {
		var err error
		f_detail, err = os.OpenFile(outputDetailFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		checkError(err)
		defer f_detail.Close()
	}
	if mode == 0 {
		hostList := []string{}

		if isPing {
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
			ch <- true
			wg.Add(1)

			pair := strings.SplitN(line, ":", 2)
			host := pair[0]
			port, _ := strconv.Atoi(pair[1])
			go scan(host, port)
		}
		wg.Wait()
	}

	if mode == 1 {
		if isPing {
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
			port = strconv.Itoa(portList[0]) + "," + strconv.Itoa(portList[1]) + "," + strconv.Itoa(portList[2]) + "..." + strconv.Itoa(portList[len(portList)-1])
		}
		if len(ipList) == 1 {
			fmt.Printf("host: %s\n", ipList[0])
			fmt.Printf("port: %s\n", port)
		} else {
			fmt.Printf("host: %s - %s\n", ipList[0], ipList[len(ipList)-1])
			fmt.Printf("port: %s\n", port)
		}

		for _, host := range ipList {
			for _, port := range portList {
				ch <- true
				wg.Add(1)
				go scan(host, port)
			}
		}
		wg.Wait()
	}
}

func scan(host string, port int) {
	open, str := isOpen(host, port)
	if open {
		ip := fmt.Sprintf("%s:%d", host, port)
		fmt.Printf("%s open    %s\n", ljust(ip, 21), str)
		f.WriteString(fmt.Sprintf("%s:%d\r\n", host, port))
		f_detail.WriteString(fmt.Sprintf("%s open    %s\n", ljust(ip, 21), str))
		openList = append(openList, ip)
	}
	<-ch
	wg.Done()
}

func isOpen(host string, port int) (bool, string) {
	var msg [128]byte
	var str string
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout)*time.Millisecond)
	if err != nil {
		//fmt.Println(err)
		return false, str
	}
	conn.SetReadDeadline((time.Now().Add(time.Millisecond * time.Duration(timeout))))
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
	for _, line := range reqHeaders {
		pair := strings.SplitN(line, ":", 2)
		if len(pair) == 2 {
			k, v := pair[0], strings.Trim(pair[1], " ")
			if strings.ToLower(k) == "host" {
				reqHost = v
			}
			headers[k] = v
		}
	}
	url := fmt.Sprintf("http://%s:%d%s", host, port, path)
	if port == 443 {
		url = fmt.Sprintf("https://%s%s", host, path)
	}
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Millisecond,
		Transport: tr,
	}
	if !redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//fmt.Println(err)
		return false, str
	}
	req.Host = reqHost
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
