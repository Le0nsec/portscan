package api

import (
	"portscan/api/controller"

	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	setCors(r)

	api := r.Group("/api")
	scan := api.Group("/scan")
	scan.POST("/create", controller.CreateScan)
	scan.GET("/records", controller.GetRecords)
	scan.GET("/detail/:id", controller.GetDetailByID)
}
