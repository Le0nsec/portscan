package config

import (
	"os"
	"path/filepath"
	"sync"
)

const (
	DataPath        = "data"
	ServerErrorCode = 500
)

var (
	SqliteDBFile = filepath.Join(DataPath, "portscan.db")
)

var (
	Wg               sync.WaitGroup
	Ch               chan bool
	Host             string
	Port             string
	LoadFile         string
	Timeout          int
	Verbose          bool
	OutputFile       string
	OutputDetailFile string
	GoroutineNum     int
	ReqHost          string
	// ReqHeaders       p.ArrayFlags
	Redirect     bool
	Path         string
	IsPing       bool
	F            *os.File
	F_detail     *os.File
	OpenList     []string
	OpenHostList []string
	//banner           bool
	//outputJSONFile   string
	//result           []HttpInfo

	Web      bool
	HostPort string
)
