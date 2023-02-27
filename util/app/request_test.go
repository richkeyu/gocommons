package app

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/richkeyu/gocommons/util/assert"

	"github.com/gin-gonic/gin"
)

type testReq struct {
	ID   int    `form:"id" xvalid:"OmitRequired(ga)"`
	Name string `form:"name" xvalid:"OmitRequired(ga)"`
}

func TestBindReqAndValid(t *testing.T) {
	var req testReq
	hreq, err := http.NewRequest("get", "http://localhost/test", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := &gin.Context{Request: hreq}
	err = BindReqAndValid(ctx, &req)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Assert(err != nil, "xvalid validate failed")
}
