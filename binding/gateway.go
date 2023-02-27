package binding

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

const (
	HeaderKeyUser = "x-user-data"
	BindTag       = "gateway"
)

var GatewayBinding = gatewayBinding{}

type gatewayBinding struct{}

func (gatewayBinding) Name() string {
	return BindTag
}

func (gatewayBinding) Bind(req *http.Request, obj interface{}) error {
	dataStr := req.Header.Get(HeaderKeyUser)
	if len(dataStr) == 0 {
		return nil
	}
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return err
	}
	dataMap := make(map[string][]string)
	for k, v := range data {
		dataMap[k] = []string{fmt.Sprintf("%v", v)}
	}

	err = binding.MapFormWithTag(obj, dataMap, BindTag)
	if err != nil {
		return err
	}
	return binding.Validator.ValidateStruct(obj)
}
