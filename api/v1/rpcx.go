// @Title
// @Description
// @Author  Wangwengang  2021/9/1 下午3:01
// @Update  Wangwengang  2021/9/1 下午3:01
package v1

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
	"github.com/wwengg/proto/identity"
	"go.uber.org/zap"

	"github.com/wwengg/arsenal/logger"
	"github.com/wwengg/arsenal/sdk/rpcx"
	"github.com/wwengg/proto/common"
	"github.com/wwengg/proto/rainbow"

	"github.com/wwengg/rainbow/response"
)

func Http2Rpcx(c *gin.Context) {
	req := protocol.NewMessage()
	req.SetMessageType(protocol.Request)
	//req.SetOneway(true)

	// json or protobuf
	contentType := c.Request.Header.Get("Content-Type")
	logger.ZapLog.Info(contentType)
	var isJson, isProtobuf bool
	isJson = strings.Contains(contentType, "application/json")
	isProtobuf = strings.Contains(contentType, "application/x-protobuf")

	if isJson {
		req.SetSerializeType(protocol.JSON)
	}
	if isProtobuf {
		req.SetSerializeType(protocol.ProtoBuffer)
	}

	if !isProtobuf && !isJson {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	// response
	r := response.New(isJson, isProtobuf, c)

	// 获取雪花id
	xclient2, err := rpcx.RpcxClientsObj.GetXClient("Identity")
	if err != nil {
		logger.ZapLog.Error("Identity service not found")
		r.Fail(common.EnumErr_ServerError, "error")
		return
	}

	identityClient := identity.NewIdentityClient(xclient2)
	reply, err := identityClient.GetId(context.Background(), nil)
	if err != nil {
		logger.ZapLog.Error(err.Error())
		r.Fail(common.EnumErr_ServerError, "error")
		return
	}
	req.SetSeq(uint64(reply.Id)) // seq

	payload, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		logger.ZapLog.Info("payload err", zap.Error(err))
		r.Fail(common.EnumErr_ServerError, "error")
		return
	}

	// servicePath & serviceMethod
	var args rainbow.HttpArgs
	if isJson {
		_ = json.Unmarshal(payload, &args)
	}
	if isProtobuf {
		_ = proto.Unmarshal(payload, &args)
	}

	req.ServicePath = args.Service
	req.ServiceMethod = args.Method
	xc, err := rpcx.RpcxClientsObj.GetXClient(args.Service)

	req.Payload = args.Data

	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, nil)
	ctx = context.WithValue(ctx, share.ResMetaDataKey, make(map[string]string))
	m, res, err := xc.SendRaw(ctx, req)
	fmt.Println(m)
	if err != nil {
		logger.ZapLog.Error("err", zap.Error(err))
		if sDec, err := base64.StdEncoding.DecodeString(err.Error()); err == nil {
			var reply rainbow.HttpReply
			err = json.Unmarshal(sDec, &reply)
			if err == nil {
				r.Fail(reply.Code, reply.Message)
				return
			}
		}
		r.Fail(common.EnumErr_ServerError, "error")
		return
	}
	if m["X-RPCX-MessageStatusType"] == "Error" {
		logger.ZapLog.Error(m["X-RPCX-ErrorMessage"])

		r.Fail(common.EnumErr_ServerError, m["X-RPCX-ErrorMessage"])
		return
	}
	r.Success(res)
}
