package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"portscan/config"
	"portscan/database"
	"time"
)

// Timestamp 用于获取当前10位数时间戳。
func Timestamp() int {
	t := time.Now().Unix()
	return int(t)
}

// ChangeStatus 修改扫描任务状态为 1。
func ChangeStatus(id int) error {
	command := "UPDATE scan SET status = 1 WHERE id = ?;"
	_, err := database.DB.Exec(command, id)
	if err != nil {
		return err
	}
	return nil
}

type CacheInfo struct {
	Host string
	Port int
	Str  string
}

type DetailItem struct {
	Host     string     `json:"host"`
	PortList []PortItem `json:"port_list"`
}

type PortItem struct {
	Port int    `json:"port"`
	Resp string `json:"resp"`
}

// 扫描任务运行完主动触发，读取文件内容后存入数据库。
func SaveResult(id int) error {
	fName := filepath.Join(config.DataPath, fmt.Sprintf("%d.txt", id))
	fileHanle, err := os.OpenFile(fName, os.O_RDONLY, 0755)
	if err != nil {
		ChangeStatus(id)
		return err
	}
	defer fileHanle.Close()
	reader := bufio.NewReader(fileHanle)

	var cacheList []CacheInfo
	for {
		var cache CacheInfo
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		err = json.Unmarshal(line, &cache)
		if err != nil {
			return err
		}
		cacheList = append(cacheList, cache)
	}

	var detailList []DetailItem
	for _, cache := range cacheList {
		portItem := PortItem{
			Port: cache.Port,
			Resp: cache.Str,
		}

		if isInDetailList(cache.Host, &detailList) {
			// Host已存在
			for i, detail := range detailList {
				if cache.Host == detail.Host {
					detailList[i].PortList = append(detail.PortList, portItem)
				}
			}
		} else {
			// Host不存在, 新建
			detailItem := DetailItem{
				Host:     cache.Host,
				PortList: []PortItem{portItem},
			}
			detailList = append(detailList, detailItem)
		}
	}

	jsonByte, err := json.Marshal(detailList)
	if err != nil {
		return err
	}

	command := "INSERT INTO detail (scan_id, json_str) VALUES (?,?);"
	_, err = database.DB.Exec(command, id, jsonByte)
	if err != nil {
		return err
	}

	return nil
}

func isInDetailList(host string, detailList *[]DetailItem) bool {
	for _, detail := range *detailList {
		if host == detail.Host {
			return true
		}
	}
	return false
}
