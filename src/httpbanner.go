package main

var headers = map[string]string{
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36",
}

type arrayFlags []string

type HttpInfo struct {
	StatusCode int    `json:"status_code"`
	Url        string `json:"url"`
	Title      string `json:"title"`
	Server     string `json:"server"`
	Length     string `json:"length"`
	Type       string `json:"type"`
}

/*
func bannerScan(list []string) {
	// prepare request headers
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

	fmt.Printf("headers:\n")
	for k, v := range headers {
		fmt.Printf("    %s: %s\n", k, v)
	}
	fmt.Printf("\nNumber of scans: %d\n", len(list))

	for _, line := range list {
		ch <- true
		wg.Add(1)

		pair := strings.SplitN(line, ":", 2)
		host := pair[0]
		port, _ := strconv.Atoi(pair[1])
		url := fmt.Sprintf("http://%s:%d%s", host, port, path)
		if port == 443 {
			url = fmt.Sprintf("https://%s%s", host, path)
		}
		go fetch(url)
	}
	wg.Wait()

	if outputJSONFile != "" {
		saveResult(result)
	}
}

func fetch(url string) {

	defer func() {
		<-ch
		wg.Done()
	}()
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
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
		return
	}
	req.Host = reqHost
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("http.Get:", err.Error())
		return
	}
	defer resp.Body.Close()

	info := &HttpInfo{}
	info.Url = url
	info.StatusCode = resp.StatusCode

	// 获取标题
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println("ioutil.ReadAll", err.Error())
		return
	}
	respBody := string(body)
	r := regexp.MustCompile(`(?i)<title>\s?(.*?)\s?</title>`)
	m := r.FindStringSubmatch(respBody)
	if len(m) == 2 {
		info.Title = m[1]
	}

	// 获取响应头Server字段
	info.Server = resp.Header.Get("Server")
	info.Length = resp.Header.Get("Content-Length")

	pair := strings.SplitN(resp.Header.Get("Content-Type"), ";", 2)
	if len(pair) == 2 {
		info.Type = pair[0]
	}

	result = append(result, *info)
	fmt.Printf("%-5d %-6s %-16s %-50s %-60s %s\n", info.StatusCode, info.Length, info.Type, info.Url, info.Server, info.Title)
}

func saveResult([]HttpInfo) {
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.OpenFile(outputJSONFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	f.WriteString(string(output))

}
*/
func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
