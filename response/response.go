// @Title
// @Description
// @Author  Wangwengang  2021/9/4 下午4:00
// @Update  Wangwengang  2021/9/4 下午4:00
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"

	"github.com/wwengg/proto/common"
	"github.com/wwengg/proto/rainbow"
)

type Response struct {
	result      rainbow.HttpReply
	c           *gin.Context
	contentType string
	isJson      bool
	isProtobuf  bool
}

func New(isJson bool, isProtobuf bool, c *gin.Context) *Response {
	return &Response{
		c:          c,
		isJson:     isJson,
		isProtobuf: isProtobuf,
	}
}

func (r *Response) Fail(code common.EnumErr, msg string) {
	if r.isJson {
		r.c.JSON(http.StatusOK, rainbow.HttpReply{
			Code:    code,
			Message: msg,
		})
	}

	if r.isProtobuf {
		protoData, _ := proto.Marshal(&rainbow.HttpReply{
			Code:    code,
			Message: msg,
		})
		r.c.Data(http.StatusOK, "application/x-protobuf", protoData)
	}
}

func (r *Response) Success(data []byte) {
	if r.isJson {
		r.c.JSON(http.StatusOK, rainbow.HttpReply{
			Code:    common.EnumErr_SUCCESS,
			Message: "success",
			Data:    data,
		})
	}

	if r.isProtobuf {
		protoData, _ := proto.Marshal(&rainbow.HttpReply{
			Code:    common.EnumErr_SUCCESS,
			Message: "success",
			Data:    data,
		})
		r.c.Data(http.StatusOK, "application/x-protobuf", protoData)

	}
}
