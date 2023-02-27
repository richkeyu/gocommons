package client

import "github.com/richkeyu/gocommons/perrors"

const (
	kopSuccessCode = 1
)

type KopResp struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
}

func (r *KopResp) IsSuccess() bool {
	return r.Code == kopSuccessCode
}

func (r *KopResp) GetError() error {
	if r.IsSuccess() {
		return nil
	}
	return perrors.GenError(r.Code, r.Message)
}
