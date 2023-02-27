package i18n

import (
	"testing"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/en_CA"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/nl"
	ut "github.com/go-playground/universal-translator"
	"github.com/stretchr/testify/assert"
)

func TestA(t *testing.T) {
	enTrans := en.New()
	universalTranslator := ut.New(enTrans, enTrans, en_CA.New(), nl.New(), fr.New())

	trans, _ := universalTranslator.GetTranslator("en")
	trans.Add("required", "required", true)
	trans2, _ := universalTranslator.GetTranslator("zh")
	trans2.Add("required", "必填", true)
	r, err := trans2.T("required2")
	assert.Nil(t, err)
	t.Log(r)
}
