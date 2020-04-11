package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/telmomarques/x360h1080p-web-config-server/config"
	"github.com/telmomarques/x360h1080p-web-config-server/hack/rtspserver"
	"github.com/telmomarques/x360h1080p-web-config-server/hack/websocketstreamserver"
)

var wwwPath = "/mnt/sdcard/hacks/web-config/www"
var port = "80"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(static.Serve("/js", static.LocalFile(wwwPath+"/js", false)))
	r.Use(static.Serve("/css", static.LocalFile(wwwPath+"/css", false)))

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File(wwwPath + "/favicon.ico")
	})

	r.GET("/", func(c *gin.Context) {
		c.Header("no-store", "expires 0")
		c.File(wwwPath + "/index.html")
	})

	apiRoutes := r.Group("/api")

	apiRoutes.GET("/hack/:hackID/config", func(c *gin.Context) {
		hackID := c.Param("hackID")
		c.File(config.GetMetaConfigFilePathForHack(hackID))
	})

	apiRoutes.GET("/hack/:hackID/service", func(c *gin.Context) {
		switch c.Param("hackID") {
		case rtspserver.ID:
			rtspserverServiceStatus := rtspserver.GetServiceStatus()
			c.JSON(http.StatusOK, rtspserverServiceStatus)

		case websocketstreamserver.ID:
			wstServiceStatus := websocketstreamserver.GetServiceStatus()
			c.JSON(http.StatusOK, wstServiceStatus)
		}

		c.Status(http.StatusNotFound)
	})

	apiRoutes.GET("/hack/:hackID/info", func(c *gin.Context) {
		switch c.Param("hackID") {
		case rtspserver.ID:
			c.String(http.StatusOK, rtspserver.Info())

		case websocketstreamserver.ID:
			c.String(http.StatusOK, websocketstreamserver.Info())
		}

		c.Status(http.StatusNotFound)
	})

	apiRoutes.GET("/hack/:hackID/endpoints", func(c *gin.Context) {
		switch c.Param("hackID") {
		case websocketstreamserver.ID:
			c.Data(http.StatusOK, gin.MIMEJSON, []byte(websocketstreamserver.Endpoints()))
		}

		c.Status(http.StatusNotFound)
	})

	apiRoutes.POST("/hack/:hackID/config", func(c *gin.Context) {
		switch c.Param("hackID") {

		case rtspserver.ID:
			var rtspserverConfig rtspserver.RTSPServerConfig
			var httpStatus = http.StatusOK

			c.Bind(&rtspserverConfig)

			success := rtspserver.SaveConfig(rtspserverConfig)
			if !success {
				httpStatus = http.StatusInternalServerError
			}

			c.Status(httpStatus)
			return

		case websocketstreamserver.ID:
			var websocketstreamConfig websocketstreamserver.WebsocketStreamConfig
			var httpStatus = http.StatusOK

			c.Bind(&websocketstreamConfig)

			success := websocketstreamserver.SaveConfig(websocketstreamConfig)
			if !success {
				httpStatus = http.StatusInternalServerError
			}

			c.Status(httpStatus)
			return
		}

		c.Status(http.StatusNotFound)
	})

	return r
}

func main() {
	if len(os.Args) == 2 {
		wwwPath = os.Args[1]
	}

	r := setupRouter()
	r.Run(":" + port)
}
