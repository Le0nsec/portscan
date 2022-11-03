package controller

import (
	"log"
	"portscan/api/model"
	"portscan/config"
	"portscan/database"
	p "portscan/portscan"
	"portscan/utils"

	"github.com/gin-gonic/gin"
)

// CreateScan 创建扫描任务。
func CreateScan(c *gin.Context) {
	var req model.CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("bindjson error:", err)
		utils.ErrorStrResp(c, config.ServerErrorCode, "Bind json error!")
		return
	}

	task := &model.Records{
		Host:    req.Host,
		Port:    req.Port,
		Threads: req.Threads,
		Timeout: req.Timeout,
		Created: utils.Timestamp(),
		Status:  0,
	}

	if err := createScan(task); err != nil {
		log.Println("create scan error:", err)
		utils.ErrorStrResp(c, config.ServerErrorCode, "Create scan failure!")
		return
	}

	log.Println("create scan success:", task.Host, task.Port)
	utils.SuccessResp(c)
}

func createScan(req *model.Records) error {
	command := "INSERT INTO scan (host, port, threads, timeout, created, status) VALUES (?,?,?,?,?,?);"
	re, err := database.DB.Exec(command, req.Host, req.Port, req.Threads, req.Timeout, req.Created, req.Status)
	if err != nil {
		return err
	}
	id, err := re.LastInsertId()
	if err != nil {
		return err
	}

	go startScan(int(id), req.Host, req.Port, req.Threads, req.Timeout)
	return nil
}

func startScan(id int, host string, port string, threads int, timeout int) {
	p.ScanByCli(host, port, id)

	err := utils.ChangeStatus(id)
	if err != nil {
		log.Println("change status error:", err)
		return
	}
}
