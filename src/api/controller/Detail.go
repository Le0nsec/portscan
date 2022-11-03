package controller

import (
	"encoding/json"
	"portscan/api/model"
	"portscan/config"
	"portscan/database"
	"portscan/utils"

	"github.com/gin-gonic/gin"
)

func GetDetailByID(c *gin.Context) {
	var jsonStr string
	var status int
	id := c.Params.ByName("id")
	if id == "" {
		utils.ErrorStrResp(c, config.ServerErrorCode, "Need id!")
		return
	}
	sql := "SELECT s.status, d.json_str FROM scan as s, detail as d WHERE d.scan_id = ? and d.scan_id = s.id;"
	row := database.DB.QueryRow(sql, id)
	err := row.Scan(&status, &jsonStr)
	if err != nil {
		utils.ErrorStrResp(c, 200, "ID does not exist!")
		return
	}
	var detailList []utils.DetailItem
	err = json.Unmarshal([]byte(jsonStr), &detailList)
	if err != nil {
		utils.ErrorStrResp(c, config.ServerErrorCode, "Unmarshal json error!")
		return
	}

	details := model.DetailResponse{
		Status: status,
		List:   detailList,
	}

	utils.SuccessResp(c, details)
}
