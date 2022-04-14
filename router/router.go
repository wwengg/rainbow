// @Title
// @Description
// @Author  Wangwengang  2022/2/16 下午1:04
// @Update  Wangwengang  2022/2/16 下午1:04
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wwengg/arsenal/logger"
	v1 "github.com/wwengg/rainbow/api/v1"
	"github.com/wwengg/rainbow/middleware"
)

func Routers() *gin.Engine {
	var Router = gin.Default()

	//pprof.Register(Router)
	logger.ZapLog.Info("use middleware start")
	// 跨域
	Router.Use(middleware.Cors())
	logger.ZapLog.Info("use middleware success")
	PrivateGroup := Router.Group("")
	{
		initRpcxApi(PrivateGroup)
	}
	return Router
}

func initRpcxApi(Router *gin.RouterGroup) {
	//RpcxRouter := Router.Group("v1").Use(middleware.OperationRecord)
	RpcxRouter := Router.Group("v1")
	{
		RpcxRouter.POST("mall", v1.Http2Rpcx)
		RpcxRouter.POST("json2base64", v1.Json2Base64)
		RpcxRouter.POST("base64ToJson", v1.Base64ToJson)
	}
}
