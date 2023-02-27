package ip

import (
	"context"
	"github.com/richkeyu/gocommons/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestUpdateDbFile(t *testing.T) {
	t.Log(os.Getwd())
	_, err := os.Stat(FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(FilePath, 0777)
			assert.Nil(t, err)
			defer os.RemoveAll(FilePath)
		}
	}
	ctx := server.NewContext(context.Background(), &gin.Context{})
	err = UpdateDbFile(ctx)
	assert.Nil(t, err)
}

func TestRun(t *testing.T) {
	Init()
	start := time.Now()
	for {
		if time.Now().Sub(start) > time.Second*60 {
			break
		}
		ip := "201.83.6.2"
		l, err := ToLocation(ip)
		assert.Nil(t, err)
		if err != nil {
			break
		}
		t.Log(db, ip, l.Country_short)
	}
}
