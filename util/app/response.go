package app

import (
	"context"
	"net/http"
	"time"

	"github.com/richkeyu/gocommons/i18n"
	"github.com/richkeyu/gocommons/perrors"
	"github.com/richkeyu/gocommons/plog"
	"github.com/richkeyu/gocommons/server"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	defaultSuccessCode = 0
	defaultErrCode     = -1
	defaultErrMsg      = "server internal error"
)

type Resp struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
}

func (r *Resp) IsSuccess() bool {
	return r.Code == defaultSuccessCode
}

func Success(c *gin.Context, data interface{}) {
	Response(c, http.StatusOK, defaultSuccessCode, data)
	return
}

func Error(c *gin.Context, httpCode int, err error) {
	trans := i18n.GetTranslatorFromGin(c)
	// 自定义错误
	if ce, ok := err.(perrors.CustomError); ok {
		var errMsg string
		if trans != nil {
			errMsg = trans.T(server.NewContext(context.Background(), c), ce.Error())
		} else {
			errMsg = ce.Error()
		}
		errCode := ce.Code()
		Response(c, httpCode, errCode, errMsg)
		return
	}

	// 参数校验错误
	if ves, ok := err.(validator.ValidationErrors); ok {
		errMsg := ""
		for _, ve := range ves {
			if trans != nil && i18n.GetValidatorTranslator(trans.Locale()) != nil {
				errMsg += ve.Translate(i18n.GetValidatorTranslator(trans.Locale()))
			} else {
				errMsg += ve.Error()
			}
		}
		Response(c, httpCode, defaultErrCode, errMsg)
		return
	}

	plog.GetDefaultFieldEntryFromGin(c).Errorf("response error: %s", err)
	Response(c, httpCode, defaultErrCode, defaultErrMsg)
	return
}

func Response(c *gin.Context, httpCode, errCode int, data interface{}) {
	var errMsg string
	if errCode != defaultSuccessCode {
		errMsg = data.(string)
		data = nil
	}
	c.JSON(httpCode, Resp{
		Code:      errCode,
		Data:      data,
		Message:   errMsg,
		Timestamp: time.Now().Unix(),
	})
	return
}
