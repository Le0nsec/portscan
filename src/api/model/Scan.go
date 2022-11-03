package model

import "portscan/utils"

// type Records struct {
// 	Host    string `json:"host" binding:"required"`
// 	Port    string `json:"port" binding:"required"`
// 	Threads int    `json:"threads" binding:"required"`
// 	Timeout int    `json:"timeout" binding:"required"`
// 	Created int    `json:"created" binding:"required"`
// 	Status  int    `json:"status" binding:"required"`
// }

type CreateRequest struct {
	Host    string `json:"host" binding:"required"`
	Port    string `json:"port" binding:"required"`
	Threads int    `json:"threads" binding:"required"`
	Timeout int    `json:"timeout" binding:"required"`
}

type Records struct {
	ID      int    `json:"id"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	Threads int    `json:"threads"`
	Timeout int    `json:"timeout"`
	Created int    `json:"created"`
	Status  int    `json:"status"`
}

type RecordsResponse struct {
	List  []Records `json:"list"`
	Total int       `json:"total"`
}

type DetailResponse struct {
	Status int                `json:"status"`
	List   []utils.DetailItem `json:"list"`
}
