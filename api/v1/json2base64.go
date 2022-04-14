// @Title
// @Description
// @Author  Wangwengang  2022/4/14 下午2:27
// @Update  Wangwengang  2022/4/14 下午2:27
package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Num int `json:"num"`
}

func Json2Base64(c *gin.Context) {
	payload, _ := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	result := make(map[string]interface{})
	_ = json.Unmarshal(payload, &result)
	c.JSON(http.StatusOK, result)
}
