// @Title
// @Description
// @Author  Wangwengang  2021/9/1 下午2:41
// @Update  Wangwengang  2021/9/1 下午2:41
package main

import (
	"fmt"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"

	"github.com/wwengg/arsenal/cache"
	"github.com/wwengg/arsenal/config"
	"github.com/wwengg/arsenal/logger"
	"github.com/wwengg/arsenal/sdk/rpcx"
	"github.com/wwengg/rainbow/router"
)

func main() {
	// Init config
	config.Viper()

	// Init logger
	logger.Setup()

	// Init Redis
	cache.Setup()

	// Init rpcx client
	rpcx.RpcxClientsObj.SetupServiceDiscovery()
	rpcx.RpcxClientsObj.SetFailMode(client.Failover)

	// start gin
	RunServer()

}

type server interface {
	ListenAndServe() error
}

func RunServer() {
	Router := router.Routers()

	//数据中台访问入口
	Router.Static("/admin", "./web")
	Router.Static("/doc", "./yfapi/_book")

	address := fmt.Sprintf(":%d", config.ConfigHub.HttpGateway.Addr)

	s := initServer(address, Router)
	time.Sleep(10 * time.Microsecond)

	logger.ZapLog.Info("server run success on ", zap.String("address", address))
	logger.ZapLog.Error(s.ListenAndServe().Error())

}

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
