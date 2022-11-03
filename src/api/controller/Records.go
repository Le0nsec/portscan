package controller

import (
	"log"
	"portscan/api/model"
	"portscan/config"
	"portscan/database"
	"portscan/utils"

	"github.com/gin-gonic/gin"
)

func GetRecords(c *gin.Context) {
	var records []model.Records

	if err := getRecords(&records); err != nil {
		log.Println("get records error", err)
		utils.ErrorStrResp(c, config.ServerErrorCode, "Get records failure!")
		return
	}

	count, err := getRecordsCount()
	if err != nil {
		log.Println("get records count error", err)
		utils.ErrorStrResp(c, config.ServerErrorCode, "Get records count failure!")
		return
	}

	resp := model.RecordsResponse{
		List:  records,
		Total: count,
	}

	utils.SuccessResp(c, resp)
}

func getRecords(records *[]model.Records) error {
	command := "SELECT id, host, port, threads, created, status FROM scan;"
	rows, err := database.DB.Query(command)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var b model.Records
		err = rows.Scan(&b.ID, &b.Host, &b.Port, &b.Threads, &b.Created, &b.Status)
		if err != nil {
			return err
		}
		*records = append(*records, b)
	}
	return rows.Err()
}

func getRecordsCount() (int, error) {
	command := "SELECT count(*) FROM scan;"
	rows, err := database.DB.Query(command)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, rows.Err()
}
