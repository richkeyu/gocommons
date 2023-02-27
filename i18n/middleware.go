package i18n

import (
	"context"

	"github.com/richkeyu/gocommons/server"

	"github.com/gin-gonic/gin"
	httpext "github.com/go-playground/pkg/v5/net/http"
	ut "github.com/go-playground/universal-translator"
)

const (
	localQueryKeyName    = "lang"
	localCookieKeyName   = "lang"
	localHeaderKeyName   = "x-accept-language"
	translatorContextKey = "__translator__"
)

var ReplaceLang = map[string]string{
	"zh":    "cn",
	"zh-CN": "cn",
	"zh-HK": "tw",
	"zh-TW": "tw",
}

// GinI18nMiddleware 初始化翻译组件
func GinI18nMiddleware(c *gin.Context) {
	var (
		t       ut.Translator
		ct      Translator
		isFound bool
	)
	// there are many ways to check, this is just checking for query param &
	// Accept-Language header but can be expanded to Cookie's etc....
	if c != nil && c.Request != nil && c.Request.URL != nil {
		locale := c.Request.URL.Query().Get(localQueryKeyName)
		if len(locale) > 0 {
			t, isFound = trans.GetTranslator(locale)
		}
	}
	// cookie
	if !isFound && c.Request != nil {
		localeCookie, err := c.Request.Cookie(localCookieKeyName)
		if err == nil {
			t, isFound = trans.GetTranslator(localeCookie.Value)
		}
	}
	// header
	if !isFound && c.Request != nil && c.Request.Header != nil {
		localeHeader := c.Request.Header.Get(localHeaderKeyName)
		if len(localeHeader) > 0 {
			t, isFound = trans.GetTranslator(localeHeader)
		}
	}
	// Accept-Language header
	if !isFound && c != nil && c.Request != nil {
		lang := httpext.AcceptedLanguages(c.Request)
		for i, l := range lang {
			if nl, ok := ReplaceLang[l]; ok {
				lang[i] = nl
			}
		}
		// get and parse the "Accept-Language" http header and return an array
		t, isFound = trans.FindTranslator(lang...)
	}

	if isFound {
		ct = newTranslator(t)
	} else {
		ct = defaultTrans
	}

	// I would normally wrap ut.Translator with one with my own functions in order
	// to handle errors and be able to use all functions from translator within the templates.
	c.Set(translatorContextKey, ct)
	c.Next()
}

func GetTranslatorFromGin(c *gin.Context) Translator {
	t, exist := c.Get(translatorContextKey)
	if !exist {
		return defaultTrans
	}
	return t.(Translator)
}

func GetTranslatorFromCtx(ctx context.Context) Translator {
	request := server.FromContext(ctx)
	if request == nil {
		return defaultTrans
	}
	t, exist := request.Get(translatorContextKey)
	if !exist {
		return defaultTrans
	}
	return t.(Translator)
}
