package i18n

import (
	"context"
	"sync"

	"github.com/richkeyu/gocommons/plog"
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
)

var (
	addTransMutex sync.Mutex
)

// Translator https://github.com/go-playground/universal-translator/blob/master/_examples/full-with-files/main.go
// kop 实现的翻译器最低要求
// wraps ut.Translator in order to handle errors transparently
// it is totally optional but recommended as it can now be used directly in
// templates and nobody can add translations where they're not supposed to.
type Translator interface {
	locales.Translator

	// T creates the translation for the locale given the 'key' and params passed in.
	// wraps ut.Translator.T to handle errors
	T(ctx context.Context, key string, params ...string) string

	GetUtTranslator() ut.Translator
}

// implements Translator interface definition above.
// 使用ut实现的版本
type translator struct {
	locales.Translator
	trans ut.Translator
}

func newTranslator(trans ut.Translator) *translator {
	return &translator{
		Translator: trans.(locales.Translator),
		trans:      trans,
	}
}

func (t *translator) GetUtTranslator() ut.Translator {
	return t.trans
}

func (t *translator) T(ctx context.Context, key string, params ...string) string {
	s, err := t.trans.T(key, params...)
	if err != nil {
		t.error(ctx, key, err)
		// 不存在时添加翻译
		//t.AddWithLock(ctx, key, params...)
		// 返回原文
		return key
	}
	return s
}

func (t *translator) AddWithLock(ctx context.Context, key string, params ...string) {
	addTransMutex.Lock()
	defer addTransMutex.Unlock()
	err := t.trans.Add(key, key, false)
	if err != nil {
		t.error(ctx, key, err)
	}
}

func (t *translator) error(ctx context.Context, key string, err error) {
	if conf.Error == ErrorInfo {
		plog.Infof(ctx, "add translating key: '%v' error: '%s'", key, err)
	} else if conf.Error == ErrorReport {
		plog.Errorf(ctx, "add translating key: '%v' error: '%s'", key, err)
	}
}
