package app

import (
	"fmt"
	"github.com/richkeyu/gocommons/util/assert"
	"testing"
)

type FormReq struct {
	ID   int    `json:"id" xvalid:"OmitRequired(ga)"`
	Name string `json:"name" xvalid:"OmitRequired(ga)"`
}

type FormReq2 struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestOmitRequired_IsSatisfied(t *testing.T) {
	var fr FormReq

	xvalid := Validation{}
	ok, err := xvalid.Check(&fr, &OmitRequired{})
	if !ok {
		fmt.Println(err)
	}

	assert.Assert(ok == false, "1.not passed omit required check")

	fr.ID = 1

	ok, err = xvalid.Check(&fr, &OmitRequired{})

	assert.Assert(ok == true, "2.not passed omit required check")

	var fr2 FormReq2
	fr2.ID = 2
	ok, err = xvalid.Check(&fr2, &OmitRequired{})

	assert.Assert(ok == true, "ignore omit require check")
}

type embedReq struct {
	FormReq
	Password string `json:"password" xvalid:"OmitRequired(ga)"`
}

func TestOmitRequiredEmbedParse(t *testing.T) {
	var embedReq embedReq

	xvalid := Validation{}
	ok, err := xvalid.Check(&embedReq, &OmitRequired{})
	if err != nil {
		fmt.Println(err)
	}

	assert.Assert(ok == false, "1.omit required embed parse err")

	embedReq.ID = 1
	ok, err = xvalid.Check(&embedReq, &OmitRequired{})
	if err != nil {
		fmt.Println(err)
	}

	assert.Assert(ok == true, "2.omit required embed parse err")
}
