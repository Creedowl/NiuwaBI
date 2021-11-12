package utils

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

func AutoWrap(handler interface{}) gin.HandlerFunc {
	value := reflect.ValueOf(handler)
	if value.Kind() != reflect.Func {
		logrus.Fatalln("handler is not function")
	}
	return func(c *gin.Context) {
		params := []reflect.Value{reflect.ValueOf(c)}
		var req string

		paramsLen := value.Type().NumIn()
		if paramsLen < 1 || paramsLen > 2 {
			logrus.Fatalln("mismatch param length")
		}
		if paramsLen == 2 {
			param := reflect.New(value.Type().In(1)).Interface()
			err := c.ShouldBind(param)
			if err != nil {
				logrus.Fatalf("failed to bind param: %+v\n", err)
			}
			marshaled, _ := jsoniter.Marshal(param)
			req = string(marshaled)
			params = append(params, reflect.ValueOf(param).Elem())
		}

		outsLen := value.Type().NumOut()
		if outsLen < 1 || outsLen > 2 {
			logrus.Fatalln("too many outs")
		}
		if !value.Type().Out(outsLen - 1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			logrus.Fatalln("last output must be error")
		}

		results := value.Call(params)

		var resp gin.H
		if results[1].IsNil() {
			resp = gin.H{
				"code":    200,
				"message": "success",
				"data":    results[0].Interface(),
			}
		} else {
			resp = gin.H{
				"code":    400,
				"message": results[1].Interface().(error).Error(),
				"data":    nil,
			}
		}
		if req != "" {
			logrus.Infof("request param: %s", req)
		}
		marshaled, _ := jsoniter.Marshal(resp)
		logrus.Infof("response: %s", string(marshaled))
		c.JSON(http.StatusOK, resp)
	}
}
