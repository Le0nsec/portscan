package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setCors(r *gin.Engine) {
	conf := cors.DefaultConfig()
	conf.AllowAllOrigins = true
	// conf.AllowHeaders = append(conf.AllowHeaders, "Authorization")
	conf.AllowCredentials = true
	// conf.AllowMethods = append(conf.AllowMethods, "OPTIONS")
	r.Use(cors.New(conf))
}
