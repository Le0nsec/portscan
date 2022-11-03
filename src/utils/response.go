package utils

import "github.com/gin-gonic/gin"

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ErrorStrResp(c *gin.Context, code int, msg string) {
	c.JSON(200, Resp{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
	c.Abort()
}

func SuccessResp(c *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, Resp{
			Code: 200,
			Msg:  "success",
			Data: nil,
		})
		return
	}
	c.JSON(200, Resp{
		Code: 200,
		Msg:  "success",
		Data: data[0],
	})
}
