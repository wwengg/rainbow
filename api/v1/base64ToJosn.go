// @Title
// @Description
// @Author  Wangwengang  2022/4/14 下午6:52
// @Update  Wangwengang  2022/4/14 下午6:52
package v1

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Base64ToJson(c *gin.Context) {
	payload, _ := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if sDec, err := base64.StdEncoding.DecodeString(string(payload)); err == nil {
		result := make(map[string]interface{})
		err = json.Unmarshal(sDec, &result)
		if err == nil {
			c.JSON(http.StatusOK, result)
			return
		}
	}
	c.JSON(http.StatusOK, payload)
}
